package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go-purple/fetchtitle"

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
}

func NewEvent(proto string, etype string, ch string, args ...interface{}) *Event {
	this := &Event{}
	this.Proto = proto
	this.EType = etype
	this.Chan = ch
	this.Args = args
	return this
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
	ctx                  *Context
	mflt                 *MessageFilter
	tracker              state.Tracker
	masterReconnectTimes int
	doneC                chan bool
}

func NewRoundTable() *RoundTable {
	this := &RoundTable{}
	this.ctx = ctx
	this.mflt = NewMessageFilter()
	this.tracker = state.NewTracker(ircname)
	this.doneC = make(chan bool)

	if pxyurl != "" {
		fetchtitle.SetProxy(pxyurl)
	}
	return this
}

func (this *RoundTable) run() {
	go this.handleEvent()
	select {}
}
func (this *RoundTable) stop() {
	this.ctx.busch <- NewEvent(PROTO_SYS, "shutdown", "")
}

// 这个方法是可以阻塞运行的，只是把后续的事件延后处理，这样带来程序逻辑简洁。
// 从另一个方面，并不会阻塞发送事件的线程
var handledEventCount uint64

func (this *RoundTable) handleEvent() {
	for ie := range this.ctx.busch {
		handledEventCount += 1
		log.Println("begin handle event:", handledEventCount, len(this.ctx.busch))
		// log.Println(ie)
		this.ctx.msgbus.Publish(ie)
		switch ie.Proto {
		case PROTO_IRC:
			// this.handleEventIrc(ie.RawEvent.(*irc.Event))
			this.handleEventIrc(ie)
		case PROTO_TOX:
			this.handleEventTox(ie)
		case PROTO_TABLE:
			this.handleEventTable(ie)
		case PROTO_SYS:
			goto endfunc
		}
		log.Println("end handle event:", handledEventCount, len(this.ctx.busch))
	}
endfunc:
	log.Println("endup event handler")
	this.doneC <- true
}
func (this *RoundTable) done() <-chan bool {
	return this.doneC
}

func (this *RoundTable) handleEventTox(e *Event) {
	log.Printf("%+v", e)
	switch e.EType {
	case EVT_GROUP_MESSAGE:
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

		// root user
		uid := this.ctx.toxagt._tox.SelfGetPublicKey()
		if !this.ctx.acpool.has(ircname, uid) {
			log.Println("wtf, try fix")
			/*
				rets := this.ctx.acpool.Call1(func() Any { return this.ctx.acpool.add(ircname, uid) })
				rac := rets[0].(*Account)
			*/
			rac := this.ctx.acpool.add(ircname, uid)
			rac.conque <- e
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
				rac := this.ctx.acpool.add(ircname, uid)
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
				ac.conque <- e
			} else {
				ac := this.ctx.acpool.get(fromUser, peerPubkey)
				be := ac.becon.(*IrcBackend2)
				if !be.isconnected() {
					log.Println("Oh, connection broken, ", chname, be.getName(), fromUser, peerPubkey)
					// this.ctx.acpool.Call0(func() { this.ctx.acpool.remove(fromUser, peerPubkey) })
					this.ctx.acpool.remove(fromUser, peerPubkey)
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
		// plugins process
		this.ctx.busch <- NewEvent(PROTO_TABLE, EVT_FETCH_URL_META, chname, message, fromUser)

	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		ne.Args[0] = PREFIX_ACTION + ne.Args[0].(string)
		this.ctx.busch <- ne

	case EVT_JOIN_GROUP:
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

		// ircon.Join("#tox-cn123")
		ac := this.ctx.acpool.get(be.getName(), be.uid)
		for len(ac.conque) > 0 {
			e := <-ac.conque
			this.ctx.busch <- e
		}

	case EVT_GROUP_MESSAGE:
		nick := e.Args[0].(string)
		// 检查是否是root用户连接
		if be.uid != this.ctx.toxagt._tox.SelfGetPublicKey() {
			break // forward message only by root user
		}

		// 检查来源是否是我们自己的连接发的消息
		for _, ac := range this.ctx.acpool.acs {
			be := ac.becon.(*IrcBackend2)
			beme := be.ircon.Me()
			if nick == beme.Nick {
				// if suffix with ^^^, the beme.Nick contains it, and nick contains it too.
				// log.Printf("drop by my ourcon:%s, %s, %v\n", nick, be.getName(), beme)
				return
			}
		}

		// TODO 两机器人消息转发循环问题，zuck07 and zuck05...

		chname := e.Args[1].(string)
		message := e.Args[2].(string)
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
		this.ctx.busch <- NewEvent(PROTO_TABLE, EVT_FETCH_URL_META, chname, message, nick)
	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		ne.Args[2] = PREFIX_ACTION + ne.Args[2].(string)
		this.ctx.busch <- ne

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
			}
			// TODO 确定是否在该群里
			// _, ok := ac.becon.(*IrcBackend2).ircon.StateTracker().IsOn(groupTitle, nick)
			_, ok := this.tracker.IsOn(groupTitle, nick)
			if ok {
				message := fmt.Sprintf("        **%s 离开了聊天室 (quit: %s)**", nick, quitmsg)
				// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(int(groupNumber), message) })
				this.ctx.toxagt._tox.GroupMessageSend(int(groupNumber), message)
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
				this.masterReconnectTimes += 1
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
			log.Println("group not exists:", chname, strings.ToLower(chname))
		} else {
			message := fmt.Sprintf("        **%s 将话题改为：%s**", nick, topic)
			// this.ctx.toxagt.Call0(func() { this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message) })
			this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
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
			this.ctx.busch <- NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber)
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
				this.ctx.busch <- NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber)
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
				this.ctx.sendBusEvent(NewEvent(PROTO_TABLE,
					EVT_GOT_URL_META, e.Chan, titleLine, e.Args[1]))
			}
		}()
	case EVT_GOT_URL_META:
		chname := e.Chan
		message := e.Args[0].(string)
		nick := e.Args[1].(string)
		message = fmt.Sprintf("%s: %s", nick, message)

		// check has bot，
		hasTitleBot := false
		ac := this.ctx.acpool.get(ircname, this.ctx.toxagt._tox.SelfGetPublicKey())
		for _, bot := range []string{"smbot", "varia", "xmppbot", "anotitlebot", "TideBot", "ttlbot"} {
			_, ok := ac.becon.(*IrcBackend2).ircon.StateTracker().IsOn(e.Chan, bot)
			if ok {
				hasTitleBot = true
				// break
			}
		}
		found := isInBotResponseWhiteChannel(chname)

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
