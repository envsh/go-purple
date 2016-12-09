package main

import (
	"log"

	"github.com/thoj/go-ircevent"
)

type IrcBackend struct {
	BackendBase
	ircon irc.Connection
}

func NewIrcBackend(ctx *Context, name string) *IrcBackend {
	this := &IrcBackend{}
	this.ctx = ctx
	this.conque = make(chan interface{}, 123)
	this.proto = PROTO_IRC
	this.name = name

	return this
}

func (this *IrcBackend) init() {
	var name = "aaaa"
	ircon := irc.IRC(name, name)
	ircon.VerboseCallbackHandler = false
	ircon.UseTLS = true
	ircon.Debug = false

	ircon.AddCallback("*", this.onEvent)
}

func (this *IrcBackend) onEvent(e *irc.Event) {
	log.Printf("%+v\n", e)
}

func (this *IrcBackend) connect() {
	go func() {
		err := this.ircon.Connect(serverssl)
		if err != nil {
			log.Println(err)
		}
	}()
}

func NewEventFromIrcEvent(e *irc.Event) *Event {
	ne := &Event{}
	ne.Proto = PROTO_IRC
	ne.Args = make([]interface{}, 0)
	ne.RawEvent = e

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
	case "JOIN":
		ne.EType = EVT_JOIN_GROUP
	case "ERROR":
		ne.EType = EVT_DISCONNECTED
	case "QUIT":
		ne.EType = EVT_FRIEND_DISCONNECTED
	default:
		ne.EType = EVT_NONE
	}
	return ne
}
