package main

import (
	"fmt"
	"log"
	"strings"
	// "time"
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
	EVT_GROUP_MESSAGE       = "group_message"
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
		peerName, err := this.ctx.toxagt._tox.GroupPeerName(e.Args[1].(int), e.Args[2].(int))
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
		if !this.ctx.acpool.has(ircname) {
			log.Println("wtf, try fix")
			rac := this.ctx.acpool.add(ircname)
			rac.conque <- e
		} else {
			rac := this.ctx.acpool.get(ircname)
			be := rac.becon.(*IrcBackend)
			if !be.isconnected() {
				log.Println("Oh, maybe unexpected")
			}
			be.join(chname)

			// agent user
			var ac *Account
			if !this.ctx.acpool.has(fromUser) {
				ac = this.ctx.acpool.add(fromUser)
				ac.conque <- e
			} else {
				ac := this.ctx.acpool.get(fromUser)
				be := ac.becon.(*IrcBackend)
				if !be.isconnected() {
					log.Println("Oh, connection broken, ", chname)
					err := be.reconnect()
					if err != nil {
						log.Println(err)
					}
				}
				be.join(chname)
				messages := strings.Split(message, "\n") // fix multiple line message
				for _, m := range messages {
					be.sendGroupMessage(m, chname)
				}
			}
		}

	case EVT_JOIN_GROUP:

		if !this.ctx.acpool.has(ircname) {
			ac := this.ctx.acpool.add(ircname)
			ac.conque <- e
		} else {
			groupTitle := e.Chan
			chname := groupTitle
			if key, found := chmap.GetKey(groupTitle); found {
				// forward message to...
				chname = key.(string)
			}
			ac := this.ctx.acpool.get(ircname)
			be := ac.becon.(*IrcBackend)
			be.join(chname)
		}

	}
}

func (this *RoundTable) handleEventIrc(e *Event) {
	be := e.Be.(*IrcBackend)

	switch e.EType {
	case EVT_CONNECTED: // MOTD end
		// ircon.Join("#tox-cn123")
		ac := this.ctx.acpool.get(be.getName())
		for len(ac.conque) > 0 {
			e := <-ac.conque
			this.ctx.busch <- e
		}

	case EVT_GROUP_MESSAGE:
		nick := e.Args[0].(string)
		// 检查是否是root用户连接
		if be.getName() != ircname {
			break // forward message only by root user
		}
		// 检查来源是否是我们自己的连接发的消息
		if _, ok := this.ctx.acpool.acs[e.Args[0].(string)]; ok {
			break
		}

		chname := e.Args[1].(string)
		message := e.Args[2].(string)
		message = fmt.Sprintf("[%s] %s", nick, message)

		if val, found := chmap.Get(chname); found {
			chname = val.(string)
		}

		// TODO maybe multiple result
		groupNumber := this.ctx.toxagt.getToxGroupByName(chname)
		if groupNumber == -1 {
			log.Println("group not exists:", chname)
		} else {
			_, err := this.ctx.toxagt._tox.GroupMessageSend(groupNumber, message)
			if err != nil {
				// should be 1
				pno := this.ctx.toxagt._tox.GroupNumberPeers(groupNumber)
				log.Println(err, chname, groupNumber, message, pno)
			}
		}

	case EVT_JOIN_GROUP:
	case EVT_DISCONNECTED:
		// close reconnect/ by Excess Flood/
		this.ctx.acpool.remove(ircname)
	default:
		log.Println("unknown evt:", e.EType)

	}

}
