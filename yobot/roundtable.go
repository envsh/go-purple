package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	// "github.com/thoj/go-ircevent"
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
)

const MAX_BUS_QUEUE_LEN = 123

type Event struct {
	Proto string
	EType string
	Chan  string
	Args  []interface{}
	// RawEvent interface{}
	Be Backend
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
	this.Proto = e.Proto
	this.EType = e.EType
	this.Chan = e.Chan
	this.Args = make([]interface{}, 0)
	for _, arg := range e.Args {
		this.Args = append(this.Args, arg)
	}
	this.Be = e.Be
	return this
}

type RoundTable struct {
	ctx *Context
}

func NewRoundTable() *RoundTable {
	this := &RoundTable{}
	this.ctx = ctx
	return this
}

func (this *RoundTable) run() {
	go this.handleEvent()
	select {}
}

// 这个方法是可以阻塞运行的，只是把后续的事件延后处理，这样带来程序逻辑简洁。
// 从另一个方面，并不会阻塞发送事件的线程
func (this *RoundTable) handleEvent() {
	for ie := range this.ctx.busch {
		// log.Println(ie)
		switch ie.Proto {
		case PROTO_IRC:
			// this.handleEventIrc(ie.RawEvent.(*irc.Event))
			this.handleEventIrc(ie)
		case PROTO_TOX:
			this.handleEventTox(ie)
		}
	}
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
			rac := this.ctx.acpool.add(ircname, uid)
			rac.conque <- e
		} else {
			rac := this.ctx.acpool.get(ircname, uid)
			rbe := rac.becon.(*IrcBackend2)
			if !rbe.isconnected() {
				log.Println("Oh, maybe unexpected", rbe.getName())
				// TODO 怎么处理呢？
				this.ctx.acpool.remove(ircname, uid)
				rac := this.ctx.acpool.add(ircname, uid)
				rac.conque <- e
				break
			}
			if !rbe.isOn(chname) {
				rbe.join(chname)
			}

			// agent user
			var ac *Account
			if !this.ctx.acpool.has(fromUser, peerPubkey) {
				ac = this.ctx.acpool.add(fromUser, peerPubkey)
				ac.conque <- e
			} else {
				ac := this.ctx.acpool.get(fromUser, peerPubkey)
				be := ac.becon.(*IrcBackend2)
				if !be.isconnected() {
					log.Println("Oh, connection broken, ", chname, be.getName(), fromUser, peerPubkey)
					this.ctx.acpool.remove(fromUser, peerPubkey)
					ac.conque <- e
				} else {
					if !be.isOn(chname) {
						be.join(chname)
					}
					messages := strings.Split(message, "\n") // fix multiple line message
					for _, m := range messages {
						be.sendGroupMessage(m, chname)
					}
				}
			}
		}

	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		ne.Args[0] = PREFIX_ACTION + ne.Args[0].(string)
		this.ctx.busch <- ne

	case EVT_JOIN_GROUP:
		rid := this.ctx.toxagt._tox.SelfGetPublicKey()
		if !this.ctx.acpool.has(ircname, rid) {
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
			be.join(chname)
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
					ac.becon.(*IrcBackend2).ircon.Part(chname)
				}
			} else {
				log.Println("wtf", peerName, this.ctx.acpool.count(), this.ctx.acpool.getNames(5))
			}
		})

	case EVT_FRIEND_MESSAGE:
		friendNumber := e.Args[1].(uint32)
		cmd := e.Args[0].(string)
		segs := strings.Split(cmd, " ")

		switch segs[0] {
		case "info": // show friends count, groups count and group list info
			this.processInfoCmd(friendNumber)
		case "invite":
			if len(segs) > 1 {
				this.processInviteCmd(segs[1:], friendNumber)
			} else {
				this.ctx.toxagt._tox.FriendSendMessage(friendNumber, "invite what?")
			}
		case "id":
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber,
				this.ctx.toxagt._tox.SelfGetAddress())
		case "help":
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, cmdhelp)
		default:
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invalidcmd)
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
			var err error
			if strings.HasPrefix(message, PREFIX_ACTION) {
				_, err = this.ctx.toxagt._tox.GroupActionSend(groupNumber, message[len(PREFIX_ACTION):])
			} else {
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
	case EVT_GROUP_ACTION:
		ne := DupEvent(e)
		ne.EType = EVT_GROUP_MESSAGE
		ne.Args[2] = PREFIX_ACTION + ne.Args[2].(string)
		this.ctx.busch <- ne

	case EVT_JOIN_GROUP:
	case EVT_FRIEND_DISCONNECTED:

	case EVT_DISCONNECTED:
		log.Printf("%+v", e)
		// close reconnect/ by Excess Flood/
		if !this.ctx.acpool.has(be.getName(), be.uid) {
			log.Println("wtf:", be.getName(), be.uid, "//", this.ctx.acpool.getNames(5))
		}
		this.ctx.acpool.remove(be.getName(), be.uid)

	default:
		switch e.EType {
		case "PONG", "PING", "NOTICE": // omit, i known
		default:
			log.Println("unknown evt:", e.EType)
		}
	}

}

func (this *RoundTable) processInviteCmd(channels []string, friendNumber uint32) {
	t := this.ctx.toxagt._tox

	// for groupbot groups

	// for irc groups
	for _, chname := range channels {
		if chname == "" {
			this.ctx.toxagt._tox.FriendSendMessage(friendNumber, invalidcmd)
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
			continue
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
			groupNumber, err := this.ctx.toxagt._tox.AddGroupChat()
			if err != nil {
				log.Println("wtf")
			} else {
				_, err := t.GroupSetTitle(groupNumber, chname)
				_, err = t.InviteFriend(friendNumber, groupNumber)
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
