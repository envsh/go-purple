package main

import (
	"fmt"
	"log"
	"strings"
	// "time"

	"github.com/thoj/go-ircevent"
)

const (
	EVT_NONE = iota
	EVT_CONNECTED
	EVT_DISCONNECTED
	EVT_FRIEND_CONNECTED
	EVT_FRIEND_DISCONNECTED
	EVT_FRIEND_MESSAGE
	EVT_JOIN_GROUP
	EVT_GROUP_MESSAGE
)

type Event struct {
	Proto    int
	EType    int
	Chan     string
	Args     []interface{}
	RawEvent interface{}
	Be       Backend
}

func NewEvent(proto int, etype int, ch string, args ...interface{}) *Event {
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
			this.handleEventIrc(ie.RawEvent.(*irc.Event))
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
			ircon := rac.ircon
			if !ircon.Connected() {
				log.Println("Oh, maybe unexpected")
			}
			ircon.Join(chname)

			// agent user
			var ac *Account
			if !this.ctx.acpool.has(fromUser) {
				ac = this.ctx.acpool.add(fromUser)
				ac.conque <- e
			} else {
				ircon := this.ctx.acpool.get(fromUser).ircon
				if !ircon.Connected() {
					log.Println("Oh, connection broken, ", chname)
					ircon.Reconnect()
				}
				ircon.Join(chname)
				messages := strings.Split(message, "\n") // fix multiple line message
				for _, m := range messages {
					ircon.Privmsg(chname, m)
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
			ircon := this.ctx.acpool.get(ircname).ircon
			ircon.Join(chname)
		}

	}
}

func (this *RoundTable) handleEventIrc(e *irc.Event) {
	// filter logout
	switch e.Code {
	case "332": // channel title
	case "353": // channel users
	case "372":
	// case "376":
	// log.Printf("%s<- %+v", e.Connection.GetNick(), e)
	default:
		log.Printf("%s<- %+v", e.Connection.GetNick(), e)
	}

	ircon := e.Connection
	switch e.Code {
	case "376": // MOTD end
		// ircon.Join("#tox-cn123")
		ac := this.ctx.acpool.get(ircon.GetNick())
		for len(ac.conque) > 0 {
			e := <-ac.conque
			this.ctx.busch <- e
		}

	case "PING":
	case "PRIVMSG":
		// 检查是否是root用户连接
		if ircon.GetNick() != ircname {
			break // forward message only by root user
		}
		// 检查来源是否是我们自己的连接
		isourcon := false
		for name, _ := range this.ctx.acpool.acs {
			if e.Nick == name {
				isourcon = true
				break
			}
		}
		if isourcon {
			break
		}

		chname := e.Arguments[0]
		message := e.Arguments[1]
		message = fmt.Sprintf("[%s] %s", e.Nick, message)

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

	case "JOIN":
	case "ERROR":
		// close reconnect/ by Excess Flood/
		this.ctx.acpool.remove(ircname)
	case "QUIT":

	}

}
