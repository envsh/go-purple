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
		log.Println("Disconnect:", error,
			cli.StreamErrorText(),
			cli.StreamErrorCData(),
			cli.StreamErrorAppCondition(),
			"/")

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
	mucih := NewMUCInvitationHandlerX()
	mucih.HandleMUCInvitationX = func(room, from, reason, body, password string, cont bool, thread string) {

	}
	cli.RegisterMUCInvitationHandler(mucih)

	psh := NewPresenceHandlerX()
	psh.HandlePresenceX = func(ptype int, from, to, status string) {
		log.Println(ptype, from, to, status)
	}
	cli.RegisterPresenceHandler(psh)

	subh := NewSubscriptionHandlerX()
	subh.HandleSubscriptionX = func(ptype int, from, to, status string) {
		log.Println(ptype, from, to, status)
	}
	cli.RegisterSubscriptionHandler(subh)

	statsh := NewStatisticsHandlerX()
	statsh.HandleStatisticsX = func(stats Statistics) {
		log.Println("stats:", stats.Encryption)
	}
	cli.RegisterStatisticsHandler(statsh)

	tagh := NewTagHandlerX()
	tagh.HandleTagX = func() {
		log.Println()
	}
	cli.RegisterTagHandler(tagh, "", "")

	dryrun := false
	if !dryrun {
		retok := cli.Connect(true)
		log.Println(retok, cli.AuthError(), cli.StreamErrorCData())
		select {}
	}
}
