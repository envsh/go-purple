package main

import (
	"crypto/tls"
	"log"
	"strings"

	irc "github.com/fluffle/goirc/client"
)

type IrcBackend2 struct {
	RelaxCallObject
	BackendBase
	ircon *irc.Conn
	// ircfg *irc.Config
	rmers map[string]irc.Remover
}

func NewIrcBackend2(ctx *Context, name string) *IrcBackend2 {
	this := &IrcBackend2{}
	this.ctx = ctx
	this.conque = make(chan interface{}, MAX_BUS_QUEUE_LEN)
	this.proto = PROTO_IRC
	this.name = name
	this.rmers = make(map[string]irc.Remover, 0)

	this.init()
	return this
}

func (this *IrcBackend2) init() {
	var name = this.name
	ircfg := irc.NewConfig(name)
	ircfg.SSL = true

	ircfg.SSLConfig = &tls.Config{ServerName: strings.Split(serverssl, ":")[0]}
	ircfg.Server = serverssl
	ircfg.NewNick = func(n string) string { return n + "^" }
	ircon := irc.Client(ircfg)
	ircon.EnableStateTracking()

	for _, name := range ircmds {
		rmer := ircon.HandleFunc(name, this.onEvent)
		this.rmers[name] = rmer
	}
	for _, no := range []string{"353"} {
		rmer := ircon.HandleFunc(no, this.onEvent)
		this.rmers[name] = rmer
	}

	this.ircon = ircon

}

func (this *IrcBackend2) setName(name string) {
	this.name = name
	this.ircon.Nick(name)
	// this.ircon.Config().Me.IsOn("#thehehe")
}

func (this *IrcBackend2) getName() string {
	// this.ircon.Me().Nick == "Powered by GoIRC"
	// nick vs name的区别
	// zuck07 // Powered by GoIRC // Nick: zuck07 // Hostmask: goirc@
	// Real Name: Powered by GoIRC // Channels: #a #b #c
	// log.Println(ircon.Me().Nick, ircon.Me().Name, ircon.Me().String())
	if this.ircon != nil && this.ircon.Me().Nick != this.name {
		if this.ircon.Me().Nick != this.ircon.Config().NewNick(this.name) {
			log.Println("wtf", this.ircon.Me().Nick, this.name)
		}
	}
	return this.name
}

func (this *IrcBackend2) isOn(channel string) bool {
	cp, on := this.ircon.Me().IsOn(channel)
	log.Printf("%v, %v, %v, %v\n", cp, on, channel, this.getName())
	return on
}

func (this *IrcBackend2) clearEvents() {
	rmers := this.rmers
	this.rmers = make(map[string]irc.Remover, 0)

	for _, rmer := range rmers {
		rmer.Remove()
	}
}

func (this *IrcBackend2) onEvent(ircon *irc.Conn, line *irc.Line) {
	// log.Printf("%+v\n", e)
	// filter logout
	switch line.Cmd {
	case "332": // channel title
	case "353": // channel users
	case "372":
		// case "376":
		// log.Printf("%s<- %+v", e.Connection.GetNick(), e)
	case "PONG", "PING", "NOTICE": // omit, i known
	default:
		log.Printf("%s<- %+v", ircon.Me().Nick, line.Raw)

		ce := NewEventFromIrcEvent2(ircon, line)
		ce.Be = this
		this.nonblockSendBusch(ce)
	}

	switch line.Cmd {
	case "353":

	}
}

func (this *IrcBackend2) nonblockSendBusch(ce *Event) {
	select {
	case this.ctx.busch <- ce:
	default:
		log.Println("send busch blocked")
	}
}

func (this *IrcBackend2) connect() {
	//	go func() {
	log.Println(this.name)
	err := this.ircon.Connect()
	if err != nil {
		log.Println(err)
		// 并不会触发disconnect事件，需要手动触发
		ce := NewEvent(PROTO_IRC, EVT_DISCONNECTED, "unknown", err.Error())
		ce.Be = this
		this.nonblockSendBusch(ce)
	}
	log.Println(this.name, "done")
	//	}()
}

func (this *IrcBackend2) reconnect() error {
	if this.ircon.Connected() {
		return this.ircon.Close() // just close for reconnect
	}
	return nil
	// return this.ircon.Connect()
}

func (this *IrcBackend2) disconnect() {
	if this.ircon.Connected() {
		err := this.ircon.Close()
		if err != nil {
			log.Println(err)
		}
	}
	this.ircon = nil
}

func (this *IrcBackend2) isconnected() bool {
	return this.ircon.Connected()
}

func (this *IrcBackend2) join(channel string) {
	this.ircon.StateTracker().NewChannel(channel)
	this.ircon.Join(channel)
}

func (this *IrcBackend2) sendMessage(msg string, user string) bool {
	this.ircon.Privmsg(user, msg)
	return true
}

func (this *IrcBackend2) sendGroupMessage(msg string, channel string) bool {
	this.ircon.Privmsg(channel, msg)
	return true
}

func NewEventFromIrcEvent2(ircon *irc.Conn, line *irc.Line) *Event {
	ne := &Event{}
	ne.Proto = PROTO_IRC
	ne.Args = make([]interface{}, 0)
	// ne.RawEvent = e

	ne.Args = append(ne.Args, line.Nick)
	for _, arg := range line.Args {
		ne.Args = append(ne.Args, arg)
	}
	ne.Ident = line.Ident
	ne.Host = line.Host

	log.Printf("%+v\n", line)
	ne.EType = line.Cmd
	switch line.Cmd {
	case irc.CONNECTED:
		ne.EType = EVT_CONNECTED
	case irc.PRIVMSG:
		ne.Chan = line.Args[0]
		if ne.Chan == ircon.Me().Nick {
			ne.EType = EVT_FRIEND_MESSAGE
		} else {
			ne.EType = EVT_GROUP_MESSAGE
		}
	case irc.ACTION:
		ne.EType = EVT_GROUP_ACTION
		ne.Chan = line.Args[0]
	case irc.JOIN:
		ne.EType = EVT_JOIN_GROUP
	case irc.DISCONNECTED:
		ne.EType = EVT_DISCONNECTED
	case irc.QUIT:
		if line.Nick == ircon.Me().Nick {
			ne.EType = EVT_DISCONNECTED
		} else {
			ne.EType = EVT_FRIEND_DISCONNECTED
		}
	default:
		ne.EType = line.Cmd
	}
	return ne
}

// from goirc/commands.go
var ircmds = []string{
	irc.REGISTER,
	irc.CONNECTED,
	irc.DISCONNECTED,
	irc.ACTION,
	irc.AWAY,
	irc.CAP,
	irc.CTCP,
	irc.CTCPREPLY,
	irc.ERROR,
	irc.INVITE,
	irc.JOIN,
	irc.KICK,
	irc.MODE,
	irc.NICK,
	irc.NOTICE,
	irc.OPER,
	irc.PART,
	irc.PASS,
	irc.PING,
	irc.PONG,
	irc.PRIVMSG,
	irc.QUIT,
	irc.TOPIC,
	irc.USER,
	irc.VERSION,
	irc.VHOST,
	irc.WHO,
	irc.WHOIS,
}
