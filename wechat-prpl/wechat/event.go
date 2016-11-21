package wechat

import (
	"fmt"
)

type EVT_TYPE int

// 事后事件
const (
	EVT_NONE EVT_TYPE = iota
	EVT_GOT_UUID
	EVT_GOT_QRCODE // TODO 存储成链接，再发送一次事件。这样客户端的选择更多样了
	EVT_GOT_QRLINK // 存储QRCODE数据的http资源
	EVT_SCANED_DATA
	EVT_REDIR_URL
	EVT_LOGIN_STATUS
	EVT_GOT_BASEINFO
	EVT_GOT_CONTACT
	EVT_RAW_MESSAGE // 原始消息，没有解析
	EVT_GOT_MESSAGE // 格式化的消息，一条原始消息可能触发多少
	EVT_NEW_MPARTICLE
	EVT_NEW_SUBSCRIBE
	EVT_FRIEND_REQUEST
	EVT_LOGOUT
)

func (this EVT_TYPE) String() (s string) {
	switch this {
	case EVT_NONE:
		s = "EVT_NONE"
	case EVT_GOT_UUID:
		s = "EVT_GOT_UUID"
	case EVT_GOT_QRCODE:
		s = "EVT_GOT_QRCODE"
	case EVT_SCANED_DATA:
		s = "EVT_SCANED_DATA"
	case EVT_REDIR_URL:
		s = "EVT_REDIR_URL"
	case EVT_LOGIN_STATUS:
		s = "EVT_LOGIN_STATUS"
	case EVT_GOT_BASEINFO:
		s = "EVT_GOT_BASEINFO"
	case EVT_GOT_CONTACT:
		s = "EVT_GOT_CONTACT"
	case EVT_GOT_MESSAGE:
		s = "EVT_GOT_MESSAGE"
	case EVT_LOGOUT:
		s = "EVT_LOGOUT"
	default:
		s = fmt.Sprintf("EVT_%d?", this)
	}
	return
}

type Event struct {
	Type  EVT_TYPE
	SType string
	Args  []string
}

func newEvent(evt EVT_TYPE, args ...string) *Event {
	this := &Event{}
	this.Type = evt
	this.Args = args
	return this
}

func newEvent2(stype string) *Event {
	this := &Event{}
	this.SType = stype
	return this
}
