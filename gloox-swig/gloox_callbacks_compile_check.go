package gloox

import (
	"log"
)

func gloox_callbacks_interface_compile_check() {
	var conlsn ConnectionListener = NewConnectionListerX()
	log.Println(conlsn)

	// var conh ConnectionHandler = NewConnectionH

	var mucroomh MUCRoomHandler = NewMUCRoomHandlerX()
	log.Println(mucroomh)

	var mucroomcfgh MUCRoomConfigHandler = NewMUCRoomConfigHandlerX()
	log.Println(mucroomcfgh)

	var presh PresenceHandler = NewPresenceHandlerX()
	log.Println(presh)

	var logh LogHandler = NewLogHandlerX()
	log.Println(logh)

	var msgh MessageHandler = NewMessageHandlerX()
	log.Println(msgh)

	var mucinvh MUCInvitationHandler = NewMUCInvitationHandlerX()
	log.Println(mucinvh)

	var subh SubscriptionHandler = NewSubscriptionHandlerX()
	log.Println(subh)

	var statsh StatisticsHandler = NewStatisticsHandlerX()
	log.Println(statsh)

	var tagh TagHandler = NewTagHandlerX()
	log.Println(tagh)
}
