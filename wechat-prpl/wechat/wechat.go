package wechat

import (
	"log"

	"github.com/kitech/colog"
)

type Wechat struct {
	OnEvent func(evt *Event)

	// private
	eqch   chan *Event
	poller *longPoll
}

func NewWechat() *Wechat {
	this := &Wechat{}

	this.eqch = make(chan *Event, 8000)
	this.poller = newLongPoll(this.eqch)
	this.poller.start()

	return this
}

// 200ms
func (this *Wechat) Iterate() {
	// 非阻塞读取chan，一次处理所有事件
	hasEvent := true
	for hasEvent {
		select {
		case evt := <-this.eqch:
			this.OnEvent(evt)
		default:
			hasEvent = false
			break
		}
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Flags())
	colog.Register()
}
