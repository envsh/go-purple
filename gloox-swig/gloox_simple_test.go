package gloox

import (
	"log"
	"testing"
)

func Test1(t *testing.T) {
	log.Println(1)
	msgh := NewMessageHandlerX()
	msgh.HandlerMessageX = func(msg Message, session MessageSession) {
		stan := SwigcptrStanza(msg.Swigcptr())
		log.Println(msg.Subtype(), stan.From(), stan.To(),
			msg.Subject(), msg.Body())
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
	conlsr.OnTLSConnectX = func(info CertInfo) {
		log.Println("TLSConnect:")
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
		log.Println("stats:", stats)
	}
	cli.RegisterStatisticsHandler(statsh)

	tagh := NewTagHandlerX()
	tagh.HandleTagX = func() {
		log.Println("tag:")
	}
	cli.RegisterTagHandler(tagh, "", "")

	// how about extensions: otr, file, attension
	// gext := NewGPGEncrypted()
	// cli.RegisterStanzaExtension(gext)

	mucroomh := NewMUCRoomHandlerX()
	mucroomh.HandleMUCParticipantPresenceX = func(room MUCRoom, part MUCRoomParticipant, presence Presence) {
		log.Println(room)
		log.Println(part)
		log.Println(presence)
	}

	room := NewMUCRoom(cli, jid, mucroomh)
	log.Println(room)

	var roomcfgh MUCRoomConfigHandler = NewMUCRoomConfigHandlerX()
	log.Println(roomcfgh)

	//
	dryrun := true
	if !dryrun {
		retok := cli.Connect(true)
		log.Println(retok, cli.AuthError(), cli.StreamErrorCData())
		select {}
	}
}
