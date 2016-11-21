package wechat

import (
	"encoding/json"
)

/*
在库内部模拟/备份一个完整会话状态数据
包括好友列表。收到/发送的消息。
兼备一点统计功能。
*/
type innerSession struct {
}

func (this *innerSession) OnEvent(evt *Event, ud interface{}) {
	wxh := ud.(*Wechat)

	switch evt.Type {
	case EVT_GOT_UUID:
	case EVT_GOT_QRCODE:
	case EVT_GOT_QRLINK:
	case EVT_SCANED_DATA:
	case EVT_GOT_BASEINFO:
	case EVT_GOT_CONTACT:
	case EVT_RAW_MESSAGE:
		msgs := parseMessages(evt.Args[0])
		for _, msg := range msgs {
			emsg, _ := json.Marshal(msg)
			wxh.eqch <- newEvent(EVT_GOT_MESSAGE, string(emsg))
		}
	case EVT_GOT_MESSAGE: // don't here

	}
}
