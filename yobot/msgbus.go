package main

import (
	"log"

	"github.com/bitly/go-simplejson"
	"github.com/nats-io/go-nats"
)

// publish all messages to nats message bus
type MsgBusClient struct {
	nc  *nats.Conn
	bus chan *Event
}

func newMsgBusClient() *MsgBusClient {
	this := &MsgBusClient{}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println(err, nil)
	}
	this.nc = nc
	this.bus = make(chan *Event, MAX_BUS_QUEUE_LEN*2)
	go this.polling()

	return this
}

func (this *MsgBusClient) Publish(e *Event) {
	this.bus <- e
}

func (this *MsgBusClient) polling() {
	for {
		select {
		case e, ok := <-this.bus:
			if e == nil || !ok {
				return
			}

			jso := simplejson.New()
			jso.Set("Chan", e.Chan)
			jso.Set("EType", e.EType)
			jso.Set("Proto", e.Proto)
			jso.Set("Args", e.Args)
			jso.Set("Ident", e.Ident)

			jsb, err := jso.Encode()
			if err != nil {
				log.Println(err)
				break
			}
			if this.nc == nil {
				this.nc, err = nats.Connect(nats.DefaultURL)
				if err != nil {
					break
				}
				log.Println("msgbus onlined.", nats.DefaultURL)
			}
			if this.nc.IsClosed() {
				log.Println("try reconnect msgbus")
				this.nc, err = this.nc.Opts.Connect()
			}
			err = this.nc.Publish("yobotmsg", jsb)
			if err != nil {
				log.Println(err)
				if err == nats.ErrConnectionClosed {
					this.nc, err = this.nc.Opts.Connect()

				}
			}
		}
	}
}
