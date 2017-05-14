package gloox

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	log.Println(1)
	msgh := NewMessageHandlerX()
	msgh.HandlerMessageX = func(subType int, from, to, subject, body string) {
		log.Println(subType, from, to, subject, body)
	}

	// from gloox_simple_config_test.go
	jid := NewJID(jidname) // eg: "simbot0@xmpp.jp/gogloox")
	cli := NewClient(jid, jidpass)
	cli.RegisterMessageHandler(msgh)

	conlsr := NewConnectionListerX()
	conlsr.OnConnectX = func() {
		log.Println()
	}
	conlsr.OnDisconnectX = func(error int) {
		log.Println(error)
	}
	conlsr.OnTLSConnectX = func() {
		log.Println()
	}
	if true {
		cli.RegisterConnectionListener(conlsr)
	}

	logh := NewLogHandlerX()
	logh.HandleLogX = func(level int, area int, l string) {
		log.Println("log", level, area, l)
	}
	cli.LogInstance().RegisterLogHandler(LogLevelDebug, int(LogAreaAll), logh)

	retok := cli.Connect(true)
	log.Println(retok)
	select {}
}
