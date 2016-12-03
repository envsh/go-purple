package main

import (
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
	Proto int
	EType int
	Chan  string
	Args  []interface{}
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
	ctx        *Context
	irconQueue chan interface{}
}

func NewRoundTable() *RoundTable {
	this := &RoundTable{}
	this.ctx = ctx
	this.irconQueue = make(chan interface{}, 123)
	return this
}

func (this *RoundTable) run() {
	go this.handleEvent()
	// this.ctx.acpool.add(ircname)
	select {}
}

func (this *RoundTable) handleEvent() {
	for ie := range this.ctx.busch {
		// log.Println(ie)
		switch e := ie.(type) {
		case *irc.Event:
			this.handleEventIrc(e)
		case *Event:
			this.handleEventTox(e)
		}
	}
}

func (this *RoundTable) handleEventTox(e *Event) {
	log.Printf("%+v", e)
	switch e.EType {
	case EVT_GROUP_MESSAGE:
		groupTitle := e.Chan
		message := e.Args[0].(string)
		var chname string = groupTitle
		if _, ok := chanMap[groupTitle]; ok {
			// forward message to...
			chname = chanMap[groupTitle]
		}
		if !this.ctx.acpool.has(ircname) {
			this.ctx.acpool.add(ircname)
			this.irconQueue <- e
		} else {
			ircon := this.ctx.acpool.get(ircname).ircon
			ircon.Join(chname)
			messages := strings.Split(message, "\n") // fix multiple line message
			for _, m := range messages {
				ircon.Privmsg(chname, m)
			}
		}
	case EVT_JOIN_GROUP:
		if !this.ctx.acpool.has(ircname) {
			this.ctx.acpool.add(ircname)
			this.irconQueue <- e
		} else {
			groupTitle := e.Chan
			chname := groupTitle
			if _, ok := chanMap[groupTitle]; ok {
				// forward message to...
				chname = chanMap[groupTitle]
			}
			ircon := this.ctx.acpool.get(ircname).ircon
			ircon.Join(chname)
		}
	}
}

func (this *RoundTable) handleEventIrc(e *irc.Event) {
	// filter
	switch e.Code {
	case "372":
	default:
		log.Printf("%+v", e)
	}

	// ircon := e.Connection
	switch e.Code {
	case "376": // MOTD end
		// ircon.Join("#tox-cn123")
		for len(this.irconQueue) > 0 {
			e := <-this.irconQueue
			this.ctx.busch <- e
		}
	case "PING":
	case "PRIVMSG":
		chname := e.Arguments[0]
		message := e.Arguments[1]

		if toname_, ok := chanMap2[chname]; ok {
			chname = toname_
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
