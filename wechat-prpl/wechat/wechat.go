package wechat

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/kitech/colog"
)

type Wechat struct {
	OnEvent func(evt *Event, userData interface{})

	// private
	eqch      chan *Event
	poller    *longPoll
	inses     *innerSession
	msgIDBase int64
}

func NewWechat() *Wechat {
	this := &Wechat{}
	this.inses = newInnerSession()
	this.msgIDBase = time.Now().UnixNano()

	this.eqch = make(chan *Event, 8000)
	this.poller = newLongPoll(this.eqch)
	this.poller.start()

	return this
}

// ~200ms
// 客户端设置的回调函数会在Iterate线程内被调用
func (this *Wechat) Iterate(userData interface{}) {
	if this.OnEvent == nil {
		log.Panicln("can not nil event handler: Wechat.OnEvent")
	}
	// 非阻塞读取chan，一次处理所有事件
	hasEvent := true
	for hasEvent {
		select {
		case evt := <-this.eqch:
			this.OnEvent(evt, userData)
			// TODO 不应该占用客户端时间有调用时间，也许和innerSession之间异步更好。
			// 比如EventLoop来做个多向的事件分发。或者进程内的pub/sub模型。
			this.inses.OnEvent(evt, this)
		default:
			hasEvent = false
			break
		}
	}
}

func (this *Wechat) Kill() {
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Flags())
	colog.Register()
}

func (this *Wechat) IsLogined() bool {
	return this.poller.state.Logined
}

// 发送拿出请求，并退出事件循环。
func (this *Wechat) Logout() bool {
	return true
}

func (this *Wechat) GetInitData() string {
	return this.poller.state.WxInitRawData
}

func (this *Wechat) GetContactData() string {
	return this.poller.state.WxContactRawData
}

func (this *Wechat) Me() *User              { return newInnerSession().me }
func (this *Wechat) Articles() []*MPArticle { return newInnerSession().mpas }
func (this *Wechat) nextMsgId() int64       { this.msgIDBase += 1; return this.msgIDBase }

// 这些发送都是需要等待响应的？
// 应该在没有登陆或者网络不通的情况会发送失败。
func (this *Wechat) SendMessage(fromUserName, toUserName string, message string) bool {
	passTicket := strings.Replace(this.poller.state.WxPassTicket, "%", "%25", -1)

	nsurl := fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxsendmsg?lang=en_US&pass_ticket=%s",
		this.poller.state.UrlBase, passTicket)
	uin, _ := strconv.Atoi(this.poller.state.Wxuin)

	BaseRequest := map[string]interface{}{
		"Uin":      uin,
		"Sid":      this.poller.state.Wxsid,
		"Skey":     this.poller.state.WxSKey,
		"DeviceID": this.poller.state.Wxdevid,
	}
	msgID := this.nextMsgId()
	msgType := 1
	Msg := map[string]interface{}{
		"Type":         msgType,
		"Content":      message,
		"FromUserName": fromUserName,
		"ToUserName":   toUserName,
		"LocalID":      fmt.Sprintf("%d", msgID),
		"ClientMsgId":  fmt.Sprintf("%d", msgID),
	}

	jso := simplejson.New()
	jso.Set("BaseRequest", BaseRequest)
	jso.Set("Msg", Msg)
	jso.Set("Scene", 0)

	postData, _ := jso.Encode()
	log.Println(string(postData))

	// TODO options 在有并发请求时可能会冲突
	this.poller.rops.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	this.poller.rops.Headers["Content-Type"] = "application/json;charset=UTF-8"
	this.poller.rops.JSON = postData
	resp, err := this.poller.rses.Post(nsurl, this.poller.rops)
	delete(this.poller.rops.Headers, "Content-Type")
	this.poller.rops.JSON = nil
	defer resp.Close()
	this.poller.saveContent("sendmsg.json", resp.Bytes(), resp, nsurl)

	if err != nil {
		log.Println(err)
	}

	p := NewParser(resp.String())
	if !p.RetOK() {
		return false
	}
	return true
}

func (this *Wechat) SendFile(friendId string, filename string) bool {
	return true
}

func (this *Wechat) SendVoice(friendId string, filename string) bool {
	return true
}

func (this *Wechat) GetBatchContact() {}
func (this *Wechat) GetIcon()         {}
func (this *Wechat) GegMsgImg()       {}
func (this *Wechat) GegMsgImgUrl()    {}
func (this *Wechat) GegMsgFileUrl()   {}
func (this *Wechat) GegMsgVoice()     {}

// promise not nil
func (this *Wechat) GetSavedata() []byte {
	return []byte(this.poller.saveData)
}

func (this *Wechat) GetUser(userName string) *User {
	if u, ok := insess.usersm[userName]; ok {
		return u
	}
	return nil
}
