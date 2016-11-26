package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/godbus/dbus"
)

type Controller struct {
	conn  *dbus.Conn
	ppobj dbus.BusObject
	sigch chan *dbus.Signal
}

func NewController() *Controller {
	this := &Controller{}
	return this
}

var service = "im.pidgin.gopurple.PurpleService"
var path = "/im/pidgin/gopurple/PurpleObject"
var iface = "im.pidgin.gopurple.PurpleInterface"

func (this *Controller) init() {
	var err error
	this.conn, err = dbus.SessionBus()
	if err != nil {
		log.Println(err)
	}

	this.ppobj = this.conn.Object(service, dbus.ObjectPath(path))

	sigstr := fmt.Sprintf(
		"type='signal',path='%s',interface='%s',sender='%s'", path, iface, service)
	this.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, sigstr)
	this.sigch = make(chan *dbus.Signal, 123)
	this.conn.Signal(this.sigch)

}

func (this *Controller) serve() {
	log.Println("waiting signals...")
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
	log.Printf("name=%s, path=%v, iface=%s, %v\n", sig.Name, sig.Path, sig.Sender, sig.Body)

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
}

func (this *Controller) BusCall(method string, args ...interface{}) *dbus.Call {
	fullMethod := fmt.Sprintf("%s.%s", iface, method)
	call := this.ppobj.Call(fullMethod, 0, args...)
	return call
}
