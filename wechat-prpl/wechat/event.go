package wechat

import ()

type EVT_TYPE int

// 事后事件
const (
	EVT_NONE EVT_TYPE = iota
	EVT_GOT_UUID
	EVT_GOT_QRCODE
	EVT_SCAN_DATA
	EVT_REDIR_URL
	EVT_LOGIN_STATUS
	EVT_GOT_BASEINFO
	EVT_GOT_CONTACT
	EVT_GOT_MESSAGE
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
	case EVT_SCAN_DATA:
		s = "EVT_SCAN_DATA"
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
	}
	return
}

type Event struct {
	Type  EVT_TYPE
	SType string
	Args  []string
}

func newEvent(evt EVT_TYPE, args []string) *Event {
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
