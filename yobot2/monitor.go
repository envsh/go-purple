package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"go-purple/purple"

	"github.com/godbus/dbus"
)

type Context struct {
	ctrls []ProtoController
}

var ctx = &Context{make([]ProtoController, 0)}

type Controller struct {
	conn  *dbus.Conn
	ppobj dbus.BusObject
	sigch chan *dbus.Signal
}

func NewController() *Controller {
	this := &Controller{}
	ctx.ctrls = append(ctx.ctrls, this)
	ctx.ctrls = append(ctx.ctrls, &IrcController{})
	ctx.ctrls = append(ctx.ctrls, &ToxController{})

	return this
}

var service = purple.GetDBusService() //  "im.pidgin.gopurple.PurpleService"
var path = purple.GetDBusPath()       // "/im/pidgin/gopurple/PurpleObject"
var iface = purple.GetDBusInterface() // "im.pidgin.gopurple.PurpleInterface"

func (this *Controller) init() {
	var err error
	this.conn, err = dbus.SessionBus()
	if err != nil {
		log.Println(err)
	}

	this.ppobj = this.conn.Object(service, dbus.ObjectPath(path))

	sigstr := fmt.Sprintf(
		"type='signal',path='%s',interface='%s',sender='%s'", path, iface, service)
	log.Println(sigstr)
	this.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, sigstr)
	this.sigch = make(chan *dbus.Signal, 123)
	this.conn.Signal(this.sigch)

}

func (this *Controller) serve() {
	log.Println("waiting signals...", purple.GoID())
	for sig := range this.sigch {
		// log.Println(sig)
		this.dispatch(sig)
	}
	log.Println("wait signals done")
}

func (this *Controller) dispatch(sig *dbus.Signal) {
	// log.Printf("name=%s, path=%v, iface=%s\n", sig.Name, sig.Path, sig.Sender)
	name := sig.Name[len(iface)+1:]
	// log.Println(name)
	if strings.HasPrefix(name, "Irc") && strings.HasSuffix(name, "Text") {
		return
	}
	log.Printf("name=%s, path=%v, iface=%s, %+v\n", sig.Name, sig.Path, sig.Sender, sig.Body)

	// TODO 根据事件名字，动态查找对应的处理函数
	switch name {
	case "AccountStatusChanged":
		call := this.BusCall("PurpleStatusGetName", sig.Body[1])
		log.Println(sig, call.Err, call.Body)
		switch call.Body[0].(string) {
		case "Away":
		case "Avaliable":
			log.Println("join...")
		}
	}

	// 还是少用点动态特性吧，不然和python一样了
	hname := "On" + name
	for _, ctrl := range ctx.ctrls {
		thisv := reflect.ValueOf(ctrl)
		mthv := thisv.MethodByName(hname)
		if mthv.IsValid() {
			argsv := []reflect.Value{reflect.ValueOf(sig), reflect.ValueOf(name)}
			mthv.Call(argsv)
		}
	}
}

func (this *Controller) BusCall(method string, args ...interface{}) *dbus.Call {
	fullMethod := fmt.Sprintf("%s.%s", iface, method)
	call := this.ppobj.Call(fullMethod, 0, args...)
	return call
}

func (this *Controller) OnAccountStatusChanged(sig *dbus.Signal, name string) {
	call := this.BusCall("PurpleStatusGetName", sig.Body[1])
	log.Println(sig, call.Err, call.Body)
	switch call.Body[0].(string) {
	case "Away":
	case "Avaliable":
		log.Println("join...")
	}
}

func (this *Controller) OnSignedOn(sig *dbus.Signal, name string) {
	log.Println(name)
	call := this.BusCall("PurpleConnectionGetAccount", sig.Body[0])
	log.Println(call.Err, call.Body)
}

func (this *Controller) getProto() string { return "roundtable" }

type ProtoController interface {
	getProto() string
}

type IrcController struct {
}

func (this *IrcController) getProto() string { return "irc" }

type ToxController struct {
}

func (this *ToxController) getProto() string { return "gotox" }
