package wechat

import (
	"log"

	"github.com/kitech/colog"
)

type Wechat struct {
	OnEvent func(evt *Event, userData interface{})

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
// 客户端设置的回调函数会在Iterate线程内被调用
func (this *Wechat) Iterate(userData interface{}) {
	// 非阻塞读取chan，一次处理所有事件
	hasEvent := true
	for hasEvent {
		select {
		case evt := <-this.eqch:
			this.OnEvent(evt, userData)
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

func (this *Wechat) IsLogined() bool {
	return this.poller.state.logined
}

// 发送拿出请求，并退出事件循环。
func (this *Wechat) Logout() bool {
	return true
}

func (this *Wechat) GetInitData() string {
	return this.poller.state.wxInitRawData
}

func (this *Wechat) GetContactData() string {
	return this.poller.state.wxContactRawData
}

// 这些发送都是需要等待响应的？
// 应该在没有登陆或者网络不通的情况会发送失败。
func (this *Wechat) SendMessage(friendId string, message string) bool {
	return true
}

func (this *Wechat) SendFile(friendId string, filename string) bool {
	return true
}

func (this *Wechat) SendVoice(friendId string, filename string) bool {
	return true
}
