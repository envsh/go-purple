package gloox

/*
#include <stdlib.h>

#include "gloox_callbacks.h"

*/
import "C"
import "unsafe"

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
	C.free((unsafe.Pointer)(from))
	C.free((unsafe.Pointer)(to))
	C.free((unsafe.Pointer)(subject))
	C.free((unsafe.Pointer)(body))
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
