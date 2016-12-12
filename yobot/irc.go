package main

import (
	"log"

	"github.com/thoj/go-ircevent"
)

type IrcBackend struct {
	BackendBase
	ircon *irc.Connection
}

func NewIrcBackend(ctx *Context, name string) *IrcBackend {
	this := &IrcBackend{}
	this.ctx = ctx
	this.conque = make(chan interface{}, MAX_BUS_QUEUE_LEN)
	this.proto = PROTO_IRC
	this.name = name

	this.init()
	return this
}

func (this *IrcBackend) init() {
	var name = this.name
	ircon := irc.IRC(name, name)
	ircon.VerboseCallbackHandler = false
	ircon.UseTLS = true
	ircon.Debug = false

	ircon.AddCallback("*", this.onEvent)

	this.ircon = ircon
}

func (this *IrcBackend) setName(name string) {
	this.name = name
	this.ircon.Nick(name)
}

func (this *IrcBackend) getName() string {
	if this.ircon.GetNick() != this.name {
		log.Println("wtf", this.ircon.GetNick(), this.name)
	}
	return this.name
}

func (this *IrcBackend) onEvent(e *irc.Event) {
	// log.Printf("%+v\n", e)
	// filter logout
	switch e.Code {
	case "332": // channel title
	case "353": // channel users
	case "372":
		// case "376":
		// log.Printf("%s<- %+v", e.Connection.GetNick(), e)
	default:
		log.Printf("%s<- %+v", e.Connection.GetNick(), e)

		ce := NewEventFromIrcEvent(e)
		ce.Be = this
		this.ctx.busch <- ce
	}
}

func (this *IrcBackend) connect() {
	go func() {
		err := this.ircon.Connect(serverssl)
		if err != nil {
			log.Println(err)
		}
	}()
}

func (this *IrcBackend) reconnect() error {
	return this.ircon.Reconnect()
}

func (this *IrcBackend) disconnect() {
	this.ircon.Disconnect()
	this.ircon = nil
}

func (this *IrcBackend) isconnected() bool {
	return this.ircon.Connected()
}

func (this *IrcBackend) join(channel string) {
	this.ircon.Join(channel)
}

func (this *IrcBackend) sendMessage(msg string, user string) bool {
	this.ircon.Privmsg(user, msg)
	return true
}

func (this *IrcBackend) sendGroupMessage(msg string, channel string) bool {
	this.ircon.Privmsg(channel, msg)
	return true
}

func NewEventFromIrcEvent(e *irc.Event) *Event {
	ne := &Event{}
	ne.Proto = PROTO_IRC
	ne.Args = make([]interface{}, 0)
	// ne.RawEvent = e

	ne.Args = append(ne.Args, e.Nick)
	for _, arg := range e.Arguments {
		ne.Args = append(ne.Args, arg)
	}

	switch e.Code {
	case "376":
		ne.EType = EVT_CONNECTED
	case "PRIVMSG":
		// TODO 如何区分好友消息和群组消息
		ne.EType = EVT_GROUP_MESSAGE
		ne.Chan = e.Arguments[0]
	case "CTCP_ACTION":
		ne.EType = EVT_GROUP_ACTION
		ne.Chan = e.Arguments[0]
	case "JOIN":
		ne.EType = EVT_JOIN_GROUP
	case "ERROR":
		ne.EType = EVT_DISCONNECTED
	case "QUIT":
		ne.EType = EVT_FRIEND_DISCONNECTED
	default:
		ne.EType = e.Code
	}
	return ne
}
