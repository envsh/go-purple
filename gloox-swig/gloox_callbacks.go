package gloox

/*
#include <stdlib.h>

#include "gloox_callbacks.h"

*/
import "C"

//import "unsafe"

import (
//	"log"
)

var cgobjs = make(map[uint64]interface{})

const objno_base uint64 = 1234567890

var objno_current = objno_base

func next_objno() uint64 {
	objno_current += 1
	return objno_current
}

type baseHandlerX struct {
	cptr  uintptr
	objno uint64
}

func (this *baseHandlerX) Swigcptr() uintptr { return this.cptr }

type MessageHandlerX struct {
	baseHandlerX

	HandlerMessageX func(msg Message, session MessageSession)
}

func NewMessageHandlerX() *MessageHandlerX {
	this := &MessageHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MessageHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}

func (this *MessageHandlerX) SwigIsMessageHandler()          {}
func (this *MessageHandlerX) HandleMessage(a ...interface{}) {}

//export MessageHandlerRCB_handleMessage
func MessageHandlerRCB_handleMessage(gobjno C.uint64_t,
	msgx C.uint64_t, sessionx C.uint64_t) {
	// log.Println(objno, subType)
	msg := SwigcptrMessage(msgx)
	session := SwigcptrMessageSession(sessionx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MessageHandlerX)
	this.HandlerMessageX(msg, session)
	DeleteMessage(msg)
	// freeAll([]*C.char{from, to, subject, body})
}

type ConnectionHandlerX struct {
	baseHandlerX
}

func (this *ConnectionHandlerX) SwigIsConnectionHandler()                                          {}
func (this *ConnectionHandlerX) HandleIncomingConnection(arg2 ConnectionBase, arg3 ConnectionBase) {}

type ConnectionListenerX struct {
	baseHandlerX

	OnConnectX    func()
	OnDisconnectX func(int)
	OnTLSConnectX func(info CertInfo)
}

func (this *ConnectionListenerX) SwigIsConnectionListener()                   {}
func (this *ConnectionListenerX) OnConnect()                                  {}
func (this *ConnectionListenerX) OnDisconnect(arg2 ConnectionError)           {}
func (this *ConnectionListenerX) OnResourceBind(arg2 string)                  {}
func (this *ConnectionListenerX) OnResourceBindError(arg2 Error)              {}
func (this *ConnectionListenerX) OnSessionCreateError(arg2 Error)             {}
func (this *ConnectionListenerX) OnTLSConnect(arg2 CertInfo) (_swig_ret bool) { return true }
func (this *ConnectionListenerX) OnStreamEvent(arg2 StreamEvent)              {}

func NewConnectionListerX() *ConnectionListenerX {
	this := &ConnectionListenerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.ConnectionListenerRCB_new(C.uint64_t(this.objno)))

	return this
}

//export ConnectionListenerRCB_onConnect_go
func ConnectionListenerRCB_onConnect_go(objno C.uint64_t) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*ConnectionListenerX)
	this.OnConnectX()
}

//export ConnectionListenerRCB_onDisconnect_go
func ConnectionListenerRCB_onDisconnect_go(gobjno C.uint64_t, error C.int) {
	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*ConnectionListenerX)
	this.OnDisconnectX(int(error))
}

//export ConnectionListenerRCB_onTLSConnect_go
func ConnectionListenerRCB_onTLSConnect_go(gobjno C.uint64_t, infox C.uint64_t) {
	info := SwigcptrCertInfo(infox)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*ConnectionListenerX)
	this.OnTLSConnectX(info)
	DeleteCertInfo(info)
}

type LogHandlerX struct {
	baseHandlerX
	HandleLogX func(int, int, string)
}

func NewLogHandlerX() *LogHandlerX {
	this := &LogHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.LogHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *LogHandlerX) Delete() {
	C.LogHandlerRCB_delete(C.uint64_t(this.objno))
	this = nil
}
func (this *LogHandlerX) SwigIsLogHandler()                             {}
func (this *LogHandlerX) HandleLog(GlooxLogLevel, GlooxLogArea, string) {}

//export LogHandlerRCB_handleLog
func LogHandlerRCB_handleLog(objno C.uint64_t, level C.int, area C.int, l *C.char) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*LogHandlerX)
	this.HandleLogX(int(level), int(area), C.GoString(l))
}

type MUCInvitationHandlerX struct {
	baseHandlerX
	HandleMUCInvitationX func(room, from, reason, body, password string, cont bool, thread string)
}

func NewMUCInvitationHandlerX() *MUCInvitationHandlerX {
	this := &MUCInvitationHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MUCInvitationHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *MUCInvitationHandlerX) Delete() {
	C.MUCInvitationHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *MUCInvitationHandlerX) SwigIsMUCInvitationHandler() {}
func (this *MUCInvitationHandlerX) HandleMUCInvitation(JID, JID, string, string, string, bool, string) {
}

func MUCInvitationHandlerRCB_handleMUCInvitation(objno C.uint64_t,
	room, from, reason, body, password *C.char, cont C.int, thread *C.char) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*MUCInvitationHandlerX)
	this.HandleMUCInvitationX(
		C.GoString(room), C.GoString(from), C.GoString(reason),
		C.GoString(body), C.GoString(password),
		c2gobool(cont), C.GoString(thread))
	freeAll([]*C.char{room, from, reason, body, password, thread})
}

type PresenceHandlerX struct {
	baseHandlerX

	HandlePresenceX func(ptype int, from, to, status string)
}

func NewPresenceHandlerX() *PresenceHandlerX {
	this := &PresenceHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.PresenceHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *PresenceHandlerX) Delete() {
	C.PresenceHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *PresenceHandlerX) SwigIsPresenceHandler()       {}
func (this *PresenceHandlerX) HandlePresence(arg2 Presence) {}

//export PresenceHandlerRCB_handlePresence
func PresenceHandlerRCB_handlePresence(objno C.uint64_t,
	ptype C.int, from, to, status *C.char) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*PresenceHandlerX)
	this.HandlePresenceX(
		int(ptype),
		C.GoString(from),
		C.GoString(to),
		C.GoString(status))
	freeAll([]*C.char{from, to, status})
}

type SubscriptionHandlerX struct {
	baseHandlerX

	HandleSubscriptionX func(subtype int, from, to, status string)
}

func NewSubscriptionHandlerX() *SubscriptionHandlerX {
	this := &SubscriptionHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.SubscriptionHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *SubscriptionHandlerX) Delete() {
	C.SubscriptionHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *SubscriptionHandlerX) SwigIsSubscriptionHandler()           {}
func (this *SubscriptionHandlerX) HandleSubscription(arg2 Subscription) {}

//export SubscriptionHandlerRCB_handleSubscription
func SubscriptionHandlerRCB_handleSubscription(objno C.uint64_t,
	ptype C.int, from, to, status *C.char) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*SubscriptionHandlerX)
	this.HandleSubscriptionX(
		int(ptype),
		C.GoString(from),
		C.GoString(to),
		C.GoString(status))
	freeAll([]*C.char{from, to, status})
}

type Statistics struct {
	TotalBytesSent            uint64
	TotalBytesReceived        uint64
	CompressedBytesSent       uint64
	CompressedBytesReceived   uint64
	UncompressedBytesSent     uint64
	UncompressedBytesReceived uint64
	TotalStanzasSent          uint64
	TotalStanzasReceived      uint64
	IqStanzasSent             uint64
	IqStanzasReceived         uint64
	MessageStanzasSent        uint64
	MessageStanzasReceived    uint64
	S10nStanzasSent           uint64
	S10nStanzasReceived       uint64
	PresenceStanzasSent       uint64
	PresenceStanzasReceived   uint64
	Encryption                int
	Compression               int
}

type StatisticsHandlerX struct {
	baseHandlerX

	HandleStatisticsX func(stats Statistics)
}

func NewStatisticsHandlerX() *StatisticsHandlerX {
	this := &StatisticsHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.StatisticsHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *StatisticsHandlerX) Delete() {
	C.StatisticsHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *StatisticsHandlerX) SwigIsStatisticsHandler()                {}
func (this *StatisticsHandlerX) HandleStatistics(stats StatisticsStruct) {}

//export StatisticsHandlerRCB_handleStatistics
func StatisticsHandlerRCB_handleStatistics(objno C.uint64_t,
	totalBytesSent,
	totalBytesReceived,
	compressedBytesSent,
	compressedBytesReceived,
	uncompressedBytesSent,
	uncompressedBytesReceived,
	totalStanzasSent,
	totalStanzasReceived,
	iqStanzasSent,
	iqStanzasReceived,
	messageStanzasSent,
	messageStanzasReceived,
	s10nStanzasSent,
	s10nStanzasReceived,
	presenceStanzasSent,
	presenceStanzasReceived C.long,
	encryption,
	compression C.int) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*StatisticsHandlerX)

	stats := Statistics{
		uint64(totalBytesSent),
		uint64(totalBytesReceived),
		uint64(compressedBytesSent),
		uint64(compressedBytesReceived),
		uint64(uncompressedBytesSent),
		uint64(uncompressedBytesReceived),
		uint64(totalStanzasSent),
		uint64(totalStanzasReceived),
		uint64(iqStanzasSent),
		uint64(iqStanzasReceived),
		uint64(messageStanzasSent),
		uint64(messageStanzasReceived),
		uint64(s10nStanzasSent),
		uint64(s10nStanzasReceived),
		uint64(presenceStanzasSent),
		uint64(presenceStanzasReceived),
		int(encryption),
		int(compression)}

	this.HandleStatisticsX(stats)
}

////
type TagHandlerX struct {
	baseHandlerX

	HandleTagX func()
}

func NewTagHandlerX() *TagHandlerX {
	this := &TagHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.TagHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *TagHandlerX) Delete() {
	C.TagHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *TagHandlerX) SwigIsTagHandler() {}
func (this *TagHandlerX) HandleTag(tag Tag) {}

//export TagHandlerRCB_handleTag
func TagHandlerRCB_handleTag(objno C.uint64_t) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*TagHandlerX)

	this.HandleTagX()
}

////////
type MUCRoomHandlerX struct {
	baseHandlerX

	HandleMUCParticipantPresenceX func(room MUCRoom, part MUCRoomParticipant, presence Presence)
	HandleMUCMessageX             func(room MUCRoom, msg Message, priv bool)
	HandleMUCRoomCreationX        func(room MUCRoom)
	HandleMUCSubjectX             func(room MUCRoom, nick, subject string)
	HandleMUCInviteDeclineX       func(room MUCRoom, jid, reason string)
	HandleMUCErrorX               func(room MUCRoom, error int)
	HandleMUCInfoX                func(room MUCRoom, features int, name string)
	HandleMUCItemsX               func(room MUCRoom)
}

func NewMUCRoomHandlerX() *MUCRoomHandlerX {
	this := &MUCRoomHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MUCRoomHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *MUCRoomHandlerX) Delete() {
	C.MUCRoomHandlerRCB_delete(C.uint64_t(this.objno))
}
func (this *MUCRoomHandlerX) SwigIsMUCRoomHandler() {}
func (this *MUCRoomHandlerX) HandleMUCParticipantPresence(arg2 MUCRoom, arg3 MUCRoomParticipant, arg4 Presence) {
}
func (this *MUCRoomHandlerX) HandleMUCMessage(arg2 MUCRoom, arg3 Message, arg4 bool)                 {}
func (this *MUCRoomHandlerX) HandleMUCRoomCreation(arg2 MUCRoom) (_swig_ret bool)                    { return false }
func (this *MUCRoomHandlerX) HandleMUCSubject(arg2 MUCRoom, arg3 string, arg4 string)                {}
func (this *MUCRoomHandlerX) HandleMUCInviteDecline(arg2 MUCRoom, arg3 JID, arg4 string)             {}
func (this *MUCRoomHandlerX) HandleMUCError(arg2 MUCRoom, arg3 GlooxStanzaError)                     {}
func (this *MUCRoomHandlerX) HandleMUCInfo(arg2 MUCRoom, arg3 int, arg4 string, arg5 DataForm)       {}
func (this *MUCRoomHandlerX) HandleMUCItems(arg2 MUCRoom, arg3 Std_list_Sl_gloox_Disco_Item_Sm__Sg_) {}

//export MUCRoomHandlerRCB_handleMUCParticipantPresence
func MUCRoomHandlerRCB_handleMUCParticipantPresence(gobjno C.uint64_t,
	roomx C.uint64_t, partx C.uint64_t, presencex C.uint64_t /* MUCRoom* room, const MUCRoomParticipant participant,	const Presence& presence */) {
	room := SwigcptrMUCRoom(roomx)
	part := SwigcptrMUCRoomParticipant(partx)
	presence := SwigcptrPresence(presencex)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)

	if this.HandleMUCParticipantPresenceX != nil {
		this.HandleMUCParticipantPresenceX(room, part, presence)
	}
	if true {
		// DeleteMUCRoomParticipant(part)
		// DeletePresence(presence)
		// DeleteMUCRoom(room)
	}
}

//export MUCRoomHandlerRCB_handleMUCMessage
func MUCRoomHandlerRCB_handleMUCMessage(gobjno C.uint64_t,
	roomx C.uint64_t, msgx C.uint64_t, privx C.int /* MUCRoom* room, const Message& msg, bool priv */) {
	room := SwigcptrMUCRoom(roomx)
	msg := SwigcptrMessage(msgx)
	priv := c2gobool(privx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCMessageX != nil {
		this.HandleMUCMessageX(room, msg, priv)
	}
	DeleteMessage(msg)
}

//export MUCRoomHandlerRCB_handleMUCRoomCreation
func MUCRoomHandlerRCB_handleMUCRoomCreation(gobjno C.uint64_t,
	roomx C.uint64_t /* MUCRoom* room*/) C.int {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCRoomCreationX != nil {
		this.HandleMUCRoomCreationX(room)
	}

	return 0
}

//export MUCRoomHandlerRCB_handleMUCSubject
func MUCRoomHandlerRCB_handleMUCSubject(gobjno C.uint64_t,
	roomx C.uint64_t, nick *C.char, subject *C.char /* MUCRoom* room, const std::string& nick,	const std::string& subject */) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCSubjectX != nil {
		this.HandleMUCSubjectX(room, C.GoString(nick), C.GoString(subject))
	}
	freeAll([]*C.char{nick, subject})
}

//export MUCRoomHandlerRCB_handleMUCInviteDecline
func MUCRoomHandlerRCB_handleMUCInviteDecline(gobjno C.uint64_t,
	roomx C.uint64_t, invitee, reason *C.char /* MUCRoom* room, const JID& invitee,	const std::string& reason */) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCInviteDeclineX != nil {
		this.HandleMUCInviteDeclineX(room, C.GoString(invitee), C.GoString(reason))
	}
	freeAll([]*C.char{invitee, reason})
}

//export MUCRoomHandlerRCB_handleMUCError
func MUCRoomHandlerRCB_handleMUCError(gobjno C.uint64_t,
	roomx C.uint64_t, error C.int /* MUCRoom* room, StanzaError error*/) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCErrorX != nil {
		this.HandleMUCErrorX(room, int(error))
	}
}

//export MUCRoomHandlerRCB_handleMUCInfo
func MUCRoomHandlerRCB_handleMUCInfo(gobjno C.uint64_t,
	roomx C.uint64_t, features C.int, name *C.char, infoForm C.uint64_t /* MUCRoom* room, int features, const std::string& name,	const DataForm* infoForm */) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCInfoX != nil {
		this.HandleMUCInfoX(room, int(features), C.GoString(name))
	}
	freeAll([]*C.char{name})
}

//export MUCRoomHandlerRCB_handleMUCItems
func MUCRoomHandlerRCB_handleMUCItems(gobjno C.uint64_t,
	roomx C.uint64_t /* MUCRoom* room, const Disco::ItemList& items */) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomHandlerX)
	if this.HandleMUCItemsX != nil {
		this.HandleMUCItemsX(room)
	}
}

////////////////
type MUCRoomConfigHandlerX struct {
	baseHandlerX

	HandleMUCConfigListX   func(arg2 MUCRoom, operation int /*, arg3 Std_list_Sl_gloox_MUCListItem_Sg_, arg4 GlooxMUCOperation*/)
	HandleMUCConfigFormX   func(arg2 MUCRoom, form DataForm /*, arg3 DataForm*/)
	HandleMUCConfigResultX func(arg2 MUCRoom, arg3 bool, operation int /*, arg4 GlooxMUCOperation*/)
	HandleMUCRequestX      func(arg2 MUCRoom, form DataForm /*, arg3 DataForm*/)
}

func (this *MUCRoomConfigHandlerX) HandleMUCConfigList(arg2 MUCRoom, arg3 Std_list_Sl_gloox_MUCListItem_Sg_, arg4 GlooxMUCOperation) {
}
func (this *MUCRoomConfigHandlerX) HandleMUCConfigForm(arg2 MUCRoom, arg3 DataForm) {}
func (this *MUCRoomConfigHandlerX) HandleMUCConfigResult(arg2 MUCRoom, arg3 bool, arg4 GlooxMUCOperation) {
}
func (this *MUCRoomConfigHandlerX) HandleMUCRequest(arg2 MUCRoom, arg3 DataForm) {}

func (this *MUCRoomConfigHandlerX) SwigIsMUCRoomConfigHandler() {}

func NewMUCRoomConfigHandlerX() *MUCRoomConfigHandlerX {
	this := &MUCRoomConfigHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MUCRoomConfigHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}

func (this *MUCRoomConfigHandlerX) delete() {
	C.MUCRoomConfigHandlerRCB_delete(C.uint64_t(this.cptr))
}

//export MUCRoomConfigHandlerRCB_handleMUCConfigList
func MUCRoomConfigHandlerRCB_handleMUCConfigList(gobjno C.uint64_t,
	roomx C.uint64_t, itemsx C.uint64_t, operation C.int) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomConfigHandlerX)
	if this.HandleMUCConfigListX != nil {
		this.HandleMUCConfigListX(room, int(operation))
	}
}

//export MUCRoomConfigHandlerRCB_handleMUCConfigForm
func MUCRoomConfigHandlerRCB_handleMUCConfigForm(gobjno C.uint64_t,
	roomx C.uint64_t, formx C.uint64_t) {
	room := SwigcptrMUCRoom(roomx)
	form := SwigcptrDataForm(formx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomConfigHandlerX)
	if this.HandleMUCConfigFormX != nil {
		this.HandleMUCConfigFormX(room, form)
	}
	DeleteDataForm(form)
}

//export MUCRoomConfigHandlerRCB_handleMUCConfigResult
func MUCRoomConfigHandlerRCB_handleMUCConfigResult(gobjno C.uint64_t,
	roomx C.uint64_t, successx C.int, operation C.int) {
	room := SwigcptrMUCRoom(roomx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomConfigHandlerX)
	if this.HandleMUCConfigResultX != nil {
		this.HandleMUCConfigResultX(room, c2gobool(successx), int(operation))
	}
}

//export MUCRoomConfigHandlerRCB_handleMUCRequest
func MUCRoomConfigHandlerRCB_handleMUCRequest(gobjno C.uint64_t,
	roomx C.uint64_t, formx C.uint64_t) {
	room := SwigcptrMUCRoom(roomx)
	form := SwigcptrDataForm(formx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MUCRoomConfigHandlerX)
	if this.HandleMUCRequestX != nil {
		this.HandleMUCRequestX(room, form)
	}
	DeleteDataForm(form)
}

////
type MessageSessionHandlerX struct {
	baseHandlerX

	HandleMessageSessionX func(MessageSession)
}

func NewMessageSessionHandlerX() *MessageSessionHandlerX {
	this := &MessageSessionHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MessageSessionHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *MessageSessionHandlerX) delete() {
	C.MessageSessionHandlerRCB_delete(C.uint64_t(this.cptr))
}
func (this *MessageSessionHandlerX) HandleMessageSession(MessageSession) {}
func (this *MessageSessionHandlerX) SwigIsMessageSessionHandler()        {}

//export MessageSessionHandlerRCB_handleMessageSession
func MessageSessionHandlerRCB_handleMessageSession(gobjno C.uint64_t, sessionx C.uint64_t) {
	session := SwigcptrMessageSession(sessionx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MessageSessionHandlerX)
	if this.HandleMessageSessionX != nil {
		this.HandleMessageSessionX(session)
	}
}

////
type MessageEventHandlerX struct {
	baseHandlerX

	HandleMessageEventX func(JID, int)
}

func NewMessageEventHandlerX() *MessageEventHandlerX {
	this := &MessageEventHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.MessageEventHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *MessageEventHandlerX) delete() {
	C.MessageEventHandlerRCB_delete(C.uint64_t(this.cptr))
}
func (this *MessageEventHandlerX) HandleMessageEvent(JID, GlooxMessageEventType) {}
func (this *MessageEventHandlerX) SwigIsMessageEventHandler()                    {}

//export MessageEventHandlerRCB_handleMessageEvent
func MessageEventHandlerRCB_handleMessageEvent(gobjno C.uint64_t, jidx C.uint64_t, event C.int) {
	jid := SwigcptrJID(jidx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*MessageEventHandlerX)
	if this.HandleMessageEventX != nil {
		this.HandleMessageEventX(jid, int(event))
	}
	DeleteJID(jid)
}

////
type ChatStateHandlerX struct {
	baseHandlerX

	HandleChatStateX func(JID, int)
}

func NewChatStateHandlerX() *ChatStateHandlerX {
	this := &ChatStateHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.ChatStateHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *ChatStateHandlerX) delete() {
	C.ChatStateHandlerRCB_delete(C.uint64_t(this.cptr))
}
func (this *ChatStateHandlerX) HandleChatState(JID, ChatStateType) {}
func (this *ChatStateHandlerX) SwigIsChatStateHandler()            {}

//export ChatStateHandlerRCB_handleChatState
func ChatStateHandlerRCB_handleChatState(gobjno C.uint64_t, jidx C.uint64_t, state C.int) {
	jid := SwigcptrJID(jidx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*ChatStateHandlerX)
	if this.HandleChatStateX != nil {
		this.HandleChatStateX(jid, int(state))
	}
	DeleteJID(jid)
}

////
type EventHandlerX struct {
	baseHandlerX

	HandleEventX func(Event)
}

func NewEventHandlerX() *EventHandlerX {
	this := &EventHandlerX{}
	this.objno = next_objno()
	cgobjs[this.objno] = this
	this.cptr = (uintptr)(C.EventHandlerRCB_new(C.uint64_t(this.objno)))

	return this
}
func (this *EventHandlerX) delete() {
	C.EventHandlerRCB_delete(C.uint64_t(this.cptr))
}
func (this *EventHandlerX) HandleEvent(Event)   {}
func (this *EventHandlerX) SwigIsEventHandler() {}

//export EventHandlerRCB_handleEvent
func EventHandlerRCB_handleEvent(gobjno C.uint64_t, eventx C.uint64_t) {
	event := SwigcptrEvent(eventx)

	thisx, _ := cgobjs[uint64(gobjno)]
	this := thisx.(*EventHandlerX)
	if this.HandleEventX != nil {
		this.HandleEventX(event)
	}
	DeleteEvent(event)
}
