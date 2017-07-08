package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go-purple/fetchtitle"
	"mkuse/hlpbot"

	"github.com/fluffle/goirc/state"
	"github.com/mvdan/xurls"
)

const (
	EVT_NONE                = "none"
	EVT_CONNECTED           = "connected"
	EVT_DISCONNECTED        = "disconnected"
	EVT_FRIEND_CONNECTED    = "friend_connected"
	EVT_FRIEND_DISCONNECTED = "friend_disconnected"
	EVT_FRIEND_MESSAGE      = "friend_message"
	EVT_JOIN_GROUP          = "join_group"
	EVT_LEAVE_GROUP         = "leave_group"
	EVT_GROUP_MESSAGE       = "group_message"
	EVT_GROUP_ACTION        = "group_action"
	EVT_GROUP_NAMES         = "group_names"
	EVT_FETCH_URL_META      = "fetch_url_meta"
	EVT_GOT_URL_META        = "got_url_meta"
)

const MAX_BUS_QUEUE_LEN = 123

type Event struct {
	Proto string
	EType string
	Chan  string
	Args  []interface{}
	// RawEvent interface{}
	Be    Backend
	Ident string
	Host  string
	// for some debug case
	No  uint64
	ttl int // 有时需要重新重理该event，重新进队。但某些逻辑需要单一处理，比如get url
}

func NewEvent(proto string, etype string, ch string, args ...interface{}) *Event {
	this := &Event{}
	this.Proto = proto
	this.EType = etype
	this.Chan = ch
	this.Args = args
	return this
}

func (e *Event) Dup() *Event {
	return DupEvent(e)
}
func DupEvent(e *Event) *Event {
	this := &Event{}
	*this = *e
	this.Args = make([]interface{}, 0)
	for _, arg := range e.Args {
		this.Args = append(this.Args, arg)
	}
	return this
}

type RoundTable struct {
	ctx     *Context
	mflt    *MessageFilter
	tracker state.Tracker
	doneC   chan bool
	assit   *helper.Assistant

	// for debug case, event handle timeout check bus
	echktoC chan eventCheckTimeout
}
type eventCheckTimeout struct {
	no     uint64
	action int // 1 begin, 2 end, 3 timeout
}

func NewRoundTable() *RoundTable {
	this := &RoundTable{}
	this.ctx = ctx
	this.mflt = NewMessageFilter()
	this.tracker = state.NewTracker(ircname)
	this.doneC = make(chan bool)
	this.assit = helper.NewAssistant()
	this.echktoC = make(chan eventCheckTimeout)

	if pxyurl != "" {
		fetchtitle.SetProxy(pxyurl)
	}
	return this
}

func (this *RoundTable) run() {
	go this.handleHelperResponse()
	go this.handleEvent()
	select {}
}
func (this *RoundTable) stop() {
	this.ctx.sendBusEvent(NewEvent(PROTO_SYS, "shutdown", ""))
}

// 这个方法是可以阻塞运行的，只是把后续的事件延后处理，这样带来程序逻辑简洁。
// 从另一个方面，并不会阻塞发送事件的线程
var handledEventCount uint64 = 0

func (this *RoundTable) handleEvent() {

	dispEvent := func(ie *Event) bool {
		nie := DupEvent(ie)
		// 这个要是可以的话，可以算作一种新和程序模型
		tmer := time.AfterFunc(3*handleEventTimeout, func() {
			log.Println(nie.No, *nie, 3*handleEventTimeout)
			panic("timeout event")
		})

		switch ie.Proto {
		case PROTO_IRC:
			// this.handleEventIrc(ie.RawEvent.(*irc.Event))
			this.handleEventIrc(ie)
		case PROTO_TOX:
			this.handleEventTox(ie)
		case PROTO_TABLE:
			this.handleEventTable(ie)
		case PROTO_SYS:
			tmer.Stop()
			return true
		}

		if !tmer.Stop() {
			log.Println(lerrorp, "wtf")
		}
		return false
	}

	for ie := range this.ctx.busch {
		// log.Println(ie)
		btime := time.Now()
		handledEventCount += 1
		log.Println("begin handle event:", handledEventCount, len(this.ctx.busch))

		ie.No = handledEventCount
		ie.ttl += 1
		breakit := dispEvent(ie)

		log.Println("end handle event:", handledEventCount, len(this.ctx.busch))
		this.ctx.rstats.handleEventTime(btime)
		this.ctx.msgbus.Publish(ie)

		if breakit {
			break
		}
	}

	log.Println("endup event handler")
	this.doneC <- true
}
func (this *RoundTable) done() <-chan bool {
	return this.doneC
}

func (this *RoundTable) handleEventTox(e *Event) {
	log.Printf("%+v, %d", e, len(e.Args))
	switch e.EType {
	case EVT_GROUP_MESSAGE: // TODO 这个逻辑有点复杂了
		if this.processGroupCmd(e.Args[0].(string), e.Args[1].(int), e.Args[2].(int)) {
			break
		}

		peerName, err := this.ctx.toxagt._tox.GroupPeerName(e.Args[1].(int), e.Args[2].(int))
		if err != nil {
			log.Println(err)
		}
		peerPubkey, err := this.ctx.toxagt._tox.GroupPeerPubkey(e.Args[1].(int), e.Args[2].(int))
		if err != nil {
			log.Println(err)
		}
		var fromUser = fmt.Sprintf("%s[t]", peerName)
		groupTitle := e.Chan
		message := e.Args[0].(string)
		var chname string = groupTitle
		if key, found := chmap.GetKey(groupTitle); found {
			// forward message to...
			chname = key.(string)
		}
		if err := this.assit.MaybeFilter(helper.PROTO_TOX, peerName, message); err != nil {
			log.Println(lwarningp, "filtered message:", err)
		}
		_, message, _ = this.assit.MaybeTransform(helper.PROTO_TOX, peerName, message)

		// plugins process
		if e.ttl <= 1 {
			defer this.ctx.sendBusEvent(NewEvent(PROTO_TABLE, EVT_FETCH_URL_META, chname, message, fromUser))
		}

		// root user, 确保root connection存在
		uid := this.ctx.toxagt._tox.SelfGetPublicKey()
		if !this.ctx.acpool.has(ircname, uid) {
			log.Println("wtf, try fix")
			/*
				rets := this.ctx.acpool.Call1(func() Any { return this.ctx.acpool.add(ircname, uid) })
				rac := rets[0].(*Account)
			*/
			// 原来阻塞的connect调用，有助于简化后续的批量处理逻辑
			rac := this.ctx.acpool.add(ircname, uid)
			rac.conque <- e
			break // nonblock-connect后添加的逻辑
		} else {
			rac := this.ctx.acpool.get(ircname, uid)
			rbe := rac.becon.(*IrcBackend2)
			if !rbe.isconnected() {
				log.Println("Oh, maybe unexpected", rbe.getName())
				// TODO 怎么处理呢？
				/*
					rets := this.ctx.acpool.Call1(func() Any {
						this.ctx.acpool.remove(ircname, uid)
						return this.ctx.acpool.add(ircname, uid)
					})
					rac := rets[0].(*Account)
				*/
				// rac := this.ctx.acpool.add(ircname, uid)
				rac.becon.connect() // fix readd problem
				rac.conque <- e
				break
			}
			if !rbe.isOn(chname) {
				// rbe.Call0(func() { rbe.join(chname) })
				rbe.join(chname)
				this.tracker.NewChannel(chname)
			}

			// agent user
			var ac *Account
			if !this.ctx.acpool.has(fromUser, peerPubkey) {
				/*
					rets := this.ctx.acpool.Call1(func() Any {
						return this.ctx.acpool.add(fromUser, peerPubkey)
					})
					ac = rets[0].(*Account)
				*/
				ac = this.ctx.acpool.add(fromUser, peerPubkey)
				ac.conque <- e // 重复的消息处理，导致插件多次处理该消息
			} else {
				ac := this.ctx.acpool.get(fromUser, peerPubkey)
				be := ac.becon.(*IrcBackend2)
				if !be.isconnected() {
					log.Println("Oh, connection broken, ", chname, be.getName(), fromUser, peerPubkey)
					// this.ctx.acpool.Call0(func() { this.ctx.acpool.remove(fromUser, peerPubkey) })
					// this.ctx.acpool.remove(fromUser, peerPubkey)
					ac.becon.connect()
					ac.conque <- e
				} else {
					if !be.isOn(chname) {
						// be.Call0(func() { be.join(chname) })
						be.join(chname)
						this.tracker.NewChannel(chname)
					}
					messages := strings.Split(message, "\n") // fix multiple line message
					for _, m := range messages {
						// be.Call0(func() { be.sendGroupMessage(m, chname) })
						be.sendGroupMessage(m, chname)
					}
				}
			}
		}

	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		ne.Args[0] = PREFIX_ACTION + ne.Args[0].(string)
		this.ctx.sendBusEvent(ne)

	case EVT_JOIN_GROUP:
		if e.ttl <= 1 && len(e.Args) > 2 {
			peerName := e.Args[2].(string)
			if true && peerName != "" && peerName != ircname &&
				peerName != "Tox User" {
				ne := DupEvent(e)
				ne.Args[1] = ""
				defer this.assit.MaybeCmdAsync(e.Chan,
					fmt.Sprintf("welcome %s", peerName), ne)
			}
		}
		rid := this.ctx.toxagt._tox.SelfGetPublicKey()
		if !this.ctx.acpool.has(ircname, rid) {
			/*
				rets := this.ctx.acpool.Call1(func() Any {
					return this.ctx.acpool.add(ircname, rid)
				})
				ac := rets[0].(*Account)
			*/
			ac := this.ctx.acpool.add(ircname, rid)
			ac.conque <- e
		} else {
			groupTitle := e.Chan
			chname := groupTitle
			if key, found := chmap.GetKey(groupTitle); found {
				// forward message to...
				chname = key.(string)
			}
			ac := this.ctx.acpool.get(ircname, rid)
			be := ac.becon.(*IrcBackend2)
			if !be.isconnected() {
				ac.becon.connect()
				ac.conque <- e
				break
			}
			// be.Call0(func() { be.join(chname) })
			be.join(chname)
			this.tracker.NewChannel(chname)
		}

	case EVT_LEAVE_GROUP:
		// 在一段时间内如果用户未再次进入该群组，则视其为离开，并关闭其irc连接
		// 或者让其也离开对应的irc群组
		chname := e.Chan
		peerName := e.Args[0].(string)
		peerPubkey := e.Args[1].(string)
		groupNumber := e.Args[2].(int)
		time.AfterFunc(leaveChannelTimeout*time.Second, func() {
			// names := this.ctx.toxagt._tox.GroupGetNames(groupNumber)
			peerPubkeys := this.ctx.toxagt._tox.GroupGetPeerPubkeys(groupNumber)
			found := false
			for _, pk := range peerPubkeys {
				if pk == peerPubkey {
					found = true
				}
			}
			if !found {
				ac := this.ctx.acpool.get(peerName+"[t]", peerPubkey)
				if ac == nil {
					log.Println("wtf", peerName, this.ctx.acpool.count(), this.ctx.acpool.getNames(5))
				} else {
					ac.becon.(*IrcBackend2).ircon.Part(chname, "tox rejoin need.")
				}
			} else {
				log.Println("wtf", peerName, this.ctx.acpool.count(), this.ctx.acpool.getNames(5))
			}
		})

		if e.ttl <= 1 {
			ne := DupEvent(e)
			ne.Args[1] = ""
			defer this.assit.MaybeCmdAsync(ne.Chan,
				fmt.Sprintf("welleave %s", peerName), ne)
		}

	case EVT_FRIEND_MESSAGE:
		friendNumber := e.Args[1].(uint32)
		cmd := e.Args[0].(string)
		segs := strings.Split(cmd, " ")

		if this.ctx.toxagt.isGroupbot(friendNumber) {
			log.Println("skip groupbot response message:", cmd)
			break
		}

		switch segs[0] {
		case "info": // show friends count, groups count and group list info
			this.processInfoCmd(friendNumber)
		case "invite":
			// TODO args check put in process function
			if len(segs) > 1 {
				this.processInviteCmd(segs[1:], friendNumber)
			} else {
				// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.FriendSendMessage(friendNumber, "invite what?") })
				this.ctx.toxagt._tox.FriendSendMessage(friendNumber, "invite what?")
			}
		case "leave": // TODO
			log.Println("should leave")
			// this.processLeaveCmd(segs, friendNumber)
		case "id":
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.FriendSendMessage(friendNumber, this.ctx.toxagt._tox.SelfGetAddress()) })
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, this.ctx.toxagt._tox.SelfGetAddress())
		case "runstats":
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, this.ctx.rstats.collect())
		case "help":
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, cmdhelp)
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.FriendSendMessage(friendNumber, cmdhelp) })
			// this.ctx.toxagt.Call(func() { this.ctx.toxagt._tox.FriendSendMessage(friendNumber, cmdhelp) })
		default:
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invalidcmd+": "+segs[0])
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invalidcmd+": "+segs[0]) })
		}
	}
}

func (this *RoundTable) handleEventIrc(e *Event) {
	be := e.Be.(*IrcBackend2)

	switch e.EType {
	case EVT_CONNECTED: // MOTD end
		log.Printf("%+v", e)
		ac := this.ctx.acpool.get(be.getName(), be.uid)
		for len(ac.conque) > 0 {
			ne := <-ac.conque
			this.ctx.sendBusEvent(ne)
		}
		// TODO rejoin. 现在的方式是需要用户触发一条消息才能rejoin
		// 怎么能判断是二次重连接呢？需要判断吗？

	case EVT_GROUP_MESSAGE:
		nick := e.Args[0].(string)
		// 检查是否是root用户连接
		if be.uid != this.ctx.toxagt._tox.SelfGetPublicKey() {
			break // forward message only by root user
		}
		// 检查来源是否是我们自己的连接发的消息
		if this.ctx.acpool.isOurs(nick) {
			break
		}
		// TODO 两机器人消息转发循环问题，zuck07 and zuck05...
		if strings.HasPrefix(nick, ircname[0:5]) {
			log.Println("maybe partizan...", nick)
		}
		// Ident == ircIdent && nick != ircname, 这应该是同类
		if e.Ident == ircIdent && nick != ircname {
			log.Println("must be partizan, omit.", nick, e.Ident)
			break
		}

		chname := e.Args[1].(string)
		message := e.Args[2].(string)

		// filter and transform message
		if err := this.assit.MaybeFilter(helper.PROTO_IRC, nick, message); err != nil {
			log.Println(lwarningp, "filtered message:", err)
		}
		nick, message, _ = this.assit.MaybeTransform(helper.PROTO_IRC, nick, message)
		rmessage := message //save for MaybeCmd call

		if strings.HasPrefix(message, PREFIX_ACTION) {
			message = fmt.Sprintf("%s[%s] %s", PREFIX_ACTION, nick, message[len(PREFIX_ACTION):])
		} else {
			message = fmt.Sprintf("[%s] %s", nick, message)
		}

		if val, found := chmap.Get(chname); found {
			chname = val.(string)
		}

		// TODO maybe multiple result
		groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
		if groupNumber == -1 {
			groupNumber = this.ctx.toxagt.getToxGroupByName(strings.ToLower(chname))
		}
		if groupNumber == -1 {
			log.Println("group not exists:", chname, strings.ToLower(chname))
		} else {
			// TODO convert to helper.Filter
			message = this.mflt.Filter(message)
			var err error // Any
			if strings.HasPrefix(message, PREFIX_ACTION) {
				/*
					rets := this.ctx.toxagt.Call2(func() (Any, Any) {
						return this.ctx.toxagt._tox.GroupActionSend(groupNumber, message[len(PREFIX_ACTION):])
					})
					err = rets[1]
				*/
				_, err = this.ctx.toxagt._tox.GroupActionSend(groupNumber, message[len(PREFIX_ACTION):])
			} else {
				/*
						log.Println("send tox group message begin...")
						rets := this.ctx.toxagt.Call2(func() (Any, Any) {
							return this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
						})
						log.Println("send tox group message end", rets)

					if len(rets) != 2 {
						log.Println("wtf", groupNumber, message)
					}
					err = rets[1]
				*/
				_, err = this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
			}
			if err != nil {
				pno := this.ctx.toxagt._tox.GroupNumberPeers(groupNumber)
				if pno == 1 { // less log
					// should be 1, should be me
				} else {
					log.Println(err, chname, groupNumber, message, pno)
				}
			}
		}
		// plugins process
		this.ctx.sendBusEvent(NewEvent(PROTO_TABLE, EVT_FETCH_URL_META, chname, rmessage, nick))

	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		// only do transform here, filter and cmd do later
		nick, message, _ := this.assit.MaybeTransform(helper.PROTO_IRC, ne.Args[0].(string), ne.Args[2].(string))
		ne.Args[0] = nick
		ne.Args[2] = message
		ne.Args[2] = PREFIX_ACTION + ne.Args[2].(string)
		this.ctx.sendBusEvent(ne)

	case EVT_FRIEND_MESSAGE:
		log.Println("not impled", e.EType)
	case EVT_JOIN_GROUP:
		if be.getName() != ircname {
			break
		}
		nick, chname := e.Args[0].(string), e.Args[1].(string)
		this.tracker.NewNick(nick)
		this.tracker.Associate(chname, nick)
		groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
		if groupNumber == -1 {
			groupNumber = this.ctx.toxagt.getToxGroupByName(strings.ToLower(chname))
		}
		if groupNumber == -1 {
			log.Println("group not exists:", chname, strings.ToLower(chname))
		} else {
			message := fmt.Sprintf("        **%s [%s@%s] 进入了聊天室。**", nick, e.Ident, e.Host)
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message) })
			this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
		}
	case "PART":
		nick, chname := e.Args[0].(string), e.Args[1].(string)
		groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
		if groupNumber == -1 {
			groupNumber = this.ctx.toxagt.getToxGroupByName(strings.ToLower(chname))
		}
		if groupNumber == -1 {
			log.Println("group not exists:", chname, strings.ToLower(chname))
		} else {
			message := fmt.Sprintf("        **%s 离开了聊天室 (quit: %s)**", nick, e.EType)
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message) })
			this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
		}
		this.tracker.Dissociate(chname, nick)
	case EVT_FRIEND_DISCONNECTED:
		nick, quitmsg := e.Args[0].(string), e.Args[1].(string)

		groupNumbers := this.ctx.toxagt._tox.GetChatList()
		for _, groupNumber := range groupNumbers {
			groupTitle, err := this.ctx.toxagt._tox.GroupGetTitle(int(groupNumber))
			if err != nil {
				log.Println(err)
			}
			// TODO 确定是否在该群里
			// _, ok := ac.becon.(*IrcBackend2).ircon.StateTracker().IsOn(groupTitle, nick)
			_, ok := this.tracker.IsOn(groupTitle, nick)
			if ok {
				message := fmt.Sprintf("        **%s 离开了聊天室 (quit: %s)**", nick, quitmsg)
				// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(int(groupNumber), message) })
				_, err := this.ctx.toxagt._tox.GroupMessageSend(int(groupNumber), message)
				if err != nil {
					log.Println(err)
				}
				this.tracker.Dissociate(groupTitle, nick)
			}
		}

	case EVT_DISCONNECTED:
		log.Printf("%+v", e)
		// close reconnect/ by Excess Flood/
		if !this.ctx.acpool.has(be.getName(), be.uid) {
			log.Println("wtf:", be.getName(), be.uid, "//", this.ctx.acpool.getNames(5))
		}
		// this.ctx.acpool.Call0(func() { this.ctx.acpool.remove(be.getName(), be.uid) })
		this.ctx.acpool.remove(be.getName(), be.uid)
		if be.getName() == ircname {
			log.Println("woo, might be reconnect the master connection")
			time.AfterFunc(2000*time.Millisecond, func() {
				if this.ctx.acpool.has(be.getName(), be.uid) {
					return
				}
				groupNames := this.ctx.toxagt.getToxGroupNames()
				log.Println("woo, reconnect the master connection to:", groupNames)
				ac := this.ctx.acpool.add(be.getName(), be.uid)
				for groupNumber, groupTitle := range groupNames {
					ac.conque <- NewEvent(PROTO_TOX, EVT_JOIN_GROUP, groupTitle, groupNumber, 0)
				}
				this.ctx.rstats.masterIrcReconnect()
			})
		}

	case "TOPIC":
		if be.getName() != ircname {
			break
		}
		nick, chname := e.Args[0].(string), e.Args[1].(string)
		topic := e.Args[2].(string)
		groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
		if groupNumber == -1 {
			groupNumber = this.ctx.toxagt.getToxGroupByName(strings.ToLower(chname))
		}
		if groupNumber == -1 {
			if nch, ok := chmap.Get(chname); ok {
				groupNumber = this.ctx.toxagt.getToxGroupByName(nch.(string))
			}
		}
		if groupNumber == -1 {
			log.Println("group not exists:", chname, strings.ToLower(chname))
		} else {
			message := fmt.Sprintf("        **%s 将话题改为：%s**", nick, topic)
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message) })
			this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
		}
	case EVT_GROUP_NAMES:
		chname := e.Args[3].(string)
		for _, nickx := range strings.Split(e.Args[4].(string), " ") {
			nick := strings.TrimLeft(nickx, "@ ")
			if this.tracker.GetNick(nick) == nil {
				this.tracker.NewNick(nick)
			}
			this.tracker.Associate(chname, nick)
			if _, ok := this.tracker.IsOn(chname, nick); !ok {
				panic("wtf")
			}
		}
	default:
		switch e.EType {
		case "PONG", "PING", "NOTICE": // omit, i known
		default:
			log.Println("unknown evt:", e.EType, e.Args)
		}
	}

}

func (this *RoundTable) processInviteCmd(channels []string, friendNumber uint32) {
	t := this.ctx.toxagt._tox

	// for groupbot groups

	// for irc groups
	for _, chname := range channels {
		if chname == "" {
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invalidcmd+": "+chname)
			continue
		}

		groupNumbers := this.ctx.toxagt._tox.GetChatList()
		found := false
		var groupNumber int
		for _, gn := range groupNumbers {
			groupTitle, err := this.ctx.toxagt._tox.GroupGetTitle(int(gn))
			if err != nil {
				log.Println("wtf")
			} else {
				if groupTitle == chname {
					found = true
					groupNumber = int(gn)
				}
			}
		}
		if found {
			log.Println("already exists:", chname)
			_, err := t.InviteFriend(friendNumber, groupNumber)
			if err != nil {
				log.Println("wtf")
			}
			// 冗余（或者并不冗余）触发一下进群事件，帮助确认成功进入 irc群
			var peerNumber int = 0
			var title = chname
			this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber))
			continue // 这个continue会导致不正确的join问题，所以要手工触发一下事件
		}

		_, err := strconv.Atoi(chname)
		if err == nil {
			// for groupbot groups
			friendNumber, err := this.ctx.toxagt._tox.FriendByPublicKey(groupbot)
			if err != nil {
				log.Println(err)
			}
			invcmd := fmt.Sprintf("invite %s", chname)
			log.Println("send groupbot invite:", chname, friendNumber, err)
			ret, err := this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invcmd)
			if err != nil {
				log.Println(err, ret)
			}
			go func() {
			}()
		} else {
			// for irc groups
			/*
				rets := this.ctx.toxagt.Call2(func() (Any, Any) {
					return this.ctx.toxagt._tox.AddGroupChat()
				})
				groupNumber, err := rets[0].(int), rets[1]
			*/
			groupNumber, err := this.ctx.toxagt._tox.AddGroupChat()
			if err != nil {
				log.Println("wtf")
			} else {
				/*
					rets := this.ctx.toxagt.Call2(func() (Any, Any) {
						t.GroupSetTitle(groupNumber, chname)
						return t.InviteFriend(friendNumber, groupNumber)
					})
					err := rets[1]
					if err != nil {
						log.Println("wtf")
					}
				*/
				t.GroupSetTitle(groupNumber, chname)
				_, err := t.InviteFriend(friendNumber, groupNumber)
				if err != nil {
					log.Println("wtf")
				}

				// 这个逻辑不会自动触发groupTitle事件，手工触发一下
				var peerNumber int = 0
				var title = chname
				this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber))
			}
		}
		// ac := this.ctx.acpool.get(chname)
	}
}

var myonlineTime = time.Now()

func (this *RoundTable) processInfoCmd(friendNumber uint32) {
	info := ""

	onlineFriendCount := this.ctx.toxagt.getOnlineFriendCount()

	info += fmt.Sprintf("Uptime: %s\n\n", time.Now().Sub(myonlineTime).String())
	info += fmt.Sprintf("Friends: %d (%d online)\n\n",
		this.ctx.toxagt._tox.SelfGetFriendListSize(), onlineFriendCount)

	// irc connections
	ircConnectionCount := len(this.ctx.acpool.acs)
	ircActiveConnectionCount := 0
	for _, ac := range this.ctx.acpool.acs {
		if ac.becon.(*IrcBackend2).isconnected() {
			ircActiveConnectionCount += 1
		}
	}
	info += fmt.Sprintf("Connections: %d (%d active)\n\n",
		ircConnectionCount, ircActiveConnectionCount)

	groupNumbers := this.ctx.toxagt._tox.GetChatList()
	for _, groupNumber := range groupNumbers {
		groupTitle, err := this.ctx.toxagt._tox.GroupGetTitle(int(groupNumber))
		if err != nil {
			log.Println(err)
		}
		peerCount := this.ctx.toxagt._tox.GroupNumberPeers(int(groupNumber))
		info += fmt.Sprintf("Group %d | Text | peers: %d | Title: %s\n\n",
			groupNumber, peerCount, groupTitle)
	}

	info = strings.TrimRight(info, "\n")
	this.ctx.toxagt._tox.FriendSendMessage(friendNumber, info)
}

// 如果是cmd则返回true
func (this *RoundTable) processGroupCmd(msg string, groupNumber, peerNumber int) bool {
	groupTitle, err := this.ctx.toxagt._tox.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println(err)
	}

	uid := this.ctx.toxagt._tox.SelfGetPublicKey()
	segs := strings.Split(msg, " ")
	if len(segs) == 1 {
		switch segs[0] {
		case "names":
		case "nc": // name count of peer irc

			ac := this.ctx.acpool.get(ircname, uid)
			if ac == nil {
				log.Println("not connected to ", groupTitle)
				this.ctx.toxagt._tox.GroupMessageSend(groupNumber, "not connected to irc:"+groupTitle)
			} else {
				ircon := ac.becon.(*IrcBackend2)
				// ircon.ircon.SendRaw("/users")
				ircon.ircon.Raw("/users")
			}
			return true
		case "ping":
			ac := this.ctx.acpool.get(ircname, uid)
			if ac == nil {
				log.Println("not connected to ", groupTitle)
				this.ctx.toxagt._tox.GroupMessageSend(groupNumber, "not connected to irc:"+groupTitle)
			} else {
				ircon := ac.becon.(*IrcBackend2)
				// ircon.ircon.SendRaw(fmt.Sprintf("/whois %s", ircname))
				ircon.ircon.Whois(ircname)
			}
			return true
		case "raw":
			this.ctx.toxagt._tox.GroupMessageSend(groupNumber, "raw what?")
			return true
		}
	} else if len(segs) > 1 {
		switch segs[0] {
		case "raw":
			ac := this.ctx.acpool.get(ircname, uid)
			if ac == nil {
				log.Println("not connected to ", groupTitle)
				this.ctx.toxagt._tox.GroupMessageSend(groupNumber, "not connected to irc:"+groupTitle)
			} else {
				ircon := ac.becon.(*IrcBackend2)
				// ircon.ircon.SendRaw(fmt.Sprintf("%s", strings.Join(segs[1:], " ")))
				ircon.ircon.Raw(fmt.Sprintf("%s", strings.Join(segs[1:], " ")))
			}
			return true
		}
	}
	return false
}

func (this *RoundTable) handleEventTable(e *Event) {
	switch e.EType {
	case EVT_FETCH_URL_META:
		msg := e.Args[0].(string)
		msgsegs := strings.Split(msg, " ")
		if len(msgsegs) > 1 && strings.HasPrefix(msgsegs[1], "Title:") {
			break // 过滤掉返回的meta结果的meta
		}

		urlst := xurls.Strict.FindAllString(msg, -1)
		go func() {
			for _, u := range urlst {
				if isLocalUrl(u) {
					log.Println("oh bad guy:", e.Args[1], u)
				}
				title, mime, err := fetchtitle.FetchMeta(u, 7)
				titleLine := fmtUrlMeta(title, mime, err, u)
				if false {
					this.ctx.sendBusEvent(NewEvent(PROTO_TABLE,
						EVT_GOT_URL_META, e.Chan, titleLine, e.Args[1]))
				}
			}
		}()

		this.assit.MaybeCmdAsync(e.Chan, msg, e)

	case EVT_GOT_URL_META:
		log.Printf("%+v\n", e)
		chname := e.Chan
		message := e.Args[0].(string)
		nick := e.Args[1].(string)
		if strings.HasSuffix(nick, "[t]") {
			//	nick = nick[0 : len(nick)-3]
		}
		if nick == "" { // 通知消息，無需暱稱
			message = fmt.Sprintf("      %s", message)
		} else {
			message = fmt.Sprintf("%s: %s", nick, message)
		}

		// check has bot，
		hasTitleBot := false
		ac := this.ctx.acpool.get(ircname, this.ctx.toxagt._tox.SelfGetPublicKey())
		// check ac, test code
		if ac == nil {
			log.Println("account not found:", ircname, this.ctx.toxagt._tox.SelfGetPublicKey())
		}
		if ac != nil && ac.becon == nil {
			log.Println("ac.becon is nil")
		}
		if ac != nil && ac.becon != nil && ac.becon.(*IrcBackend2).ircon == nil {
			log.Println("ac ircon is nil")
		}
		if ac == nil || ac.becon == nil || ac.becon.(*IrcBackend2).ircon == nil {
			time.AfterFunc(2*time.Second, func() { this.ctx.sendBusEvent(e) })
			break
		}

		for _, bot := range []string{"smbot", "varia", "xmppbot", "anotitlebot", "TideBot", "ttlbot"} {
			_, ok := ac.becon.(*IrcBackend2).ircon.StateTracker().IsOn(e.Chan, bot)
			if ok {
				hasTitleBot = true
				// break
			}
		}
		found := helper.IsInBotResponseWhiteChannel(chname)

		// 总是发到tox端上
		if true {
			// TODO move to ToxAgent.relaxSendMessage(chname)
			groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
			if groupNumber == -1 {
				groupNumber = this.ctx.toxagt.getToxGroupByName(strings.ToLower(chname))
			}
			if groupNumber == -1 {
				if nch, ok := chmap.Get(chname); ok {
					groupNumber = this.ctx.toxagt.getToxGroupByName(nch.(string))
				}
			}
			if groupNumber == -1 {
				log.Println("group not exists:", chname, strings.ToLower(chname))
			} else {
				_, err := this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
				if err != nil {
					log.Println(err, groupNumber, chname)
				}
			}
		}

		// 本群中有titlebot，或者不在白名单，则不响应此消息
		if hasTitleBot || !found {
			log.Println("got meta: ", message, hasTitleBot, found)
			break
		}
		if true {
			// TODO move to Account.relaxSendMessage(ircname)
			// find channel connection
			ac := this.ctx.acpool.get(ircname, this.ctx.toxagt._tox.SelfGetPublicKey())
			if ac == nil {
				log.Println("account not found:", ircname)
			} else {
				if nch, ok := chmap.GetKey(chname); ok {
					ac.becon.(*IrcBackend2).sendGroupMessage(message, nch.(string))
				} else {
					ac.becon.(*IrcBackend2).sendGroupMessage(message, chname)
				}
			}
		}
	}

}

func (this *RoundTable) handleHelperResponse() {
	for res := range this.assit.ResultC() {
		log.Println(res.Result)
		e := res.Extra.(*Event)
		titleLine := res.Result

		this.ctx.sendBusEvent(NewEvent(PROTO_TABLE,
			EVT_GOT_URL_META, e.Chan, titleLine, e.Args[1]))
	}
	log.Println("done")
}
