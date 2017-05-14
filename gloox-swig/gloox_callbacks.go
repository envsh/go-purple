package gloox

/*
#include <stdlib.h>

#include "gloox_callbacks.h"

*/
import "C"

//import "unsafe"

import (
// "log"
)

var cgobjs = make(map[uint64]interface{})

const objno_base uint64 = 1234567890

var objno_current = objno_base

func next_objno() uint64 {
	objno_current += 1
	return objno_current
}

type BaseHandlerX struct {
	cptr  uintptr
	objno uint64
}

func (this *BaseHandlerX) Swigcptr() uintptr { return this.cptr }

type MessageHandlerX struct {
	BaseHandlerX

	HandlerMessageX func(subType int, from, to, subject, body string)
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
func MessageHandlerRCB_handleMessage(objno C.uint64_t, subType C.int,
	from, to, subject, body *C.char) {
	// log.Println(objno, subType)
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*MessageHandlerX)
	this.HandlerMessageX(int(subType), C.GoString(from), C.GoString(to),
		C.GoString(subject), C.GoString(body))
	freeAll([]*C.char{from, to, subject, body})
}

type ConnectionHandlerX struct {
	BaseHandlerX
}

func (this *ConnectionHandlerX) SwigIsConnectionHandler() {}
func (this *ConnectionHandlerX) HandleIncomingConnection(arg2 Gloox_ConnectionBase, arg3 Gloox_ConnectionBase) {

}

type ConnectionListenerX struct {
	BaseHandlerX

	OnConnectX    func()
	OnDisconnectX func(int)
	OnTLSConnectX func()
}

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
func ConnectionListenerRCB_onDisconnect_go(objno C.uint64_t, error C.int) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*ConnectionListenerX)
	this.OnDisconnectX(int(error))
}

//export ConnectionListenerRCB_onTLSConnect_go
func ConnectionListenerRCB_onTLSConnect_go(objno C.uint64_t) {
	thisx, _ := cgobjs[uint64(objno)]
	this := thisx.(*ConnectionListenerX)
	this.OnTLSConnectX()
}

type LogHandlerX struct {
	BaseHandlerX
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
	BaseHandlerX
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
	BaseHandlerX

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
	BaseHandlerX

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
	BaseHandlerX

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
	BaseHandlerX

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
