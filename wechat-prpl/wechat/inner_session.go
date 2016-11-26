package wechat

import (
	"encoding/json"
	"sync"
)

/*
在库内部模拟/备份一个完整会话状态数据
包括好友列表。收到/发送的消息。
兼备一点统计功能。
*/
type innerSession struct {
	me     *User
	users  []*User
	usersm map[string]*User // userName => *User, userNick => *User, pinyin => *User

	mpas   []*MPArticle
	mpsubs []*MPArticle
}

var insess *innerSession = nil
var insess_once sync.Once

func newInnerSession() *innerSession {
	insess_once.Do(func() {
		this := &innerSession{}
		this.users = make([]*User, 0)
		this.usersm = make(map[string]*User, 0)
		insess = this
	})
	return insess
}

func (this *innerSession) OnEvent(evt *Event, ud interface{}) {
	wxh := ud.(*Wechat)

	switch evt.Type {
	case EVT_GOT_UUID:
	case EVT_GOT_QRCODE:
	case EVT_GOT_QRLINK:
	case EVT_SCANED_DATA:
	case EVT_GOT_BASEINFO:
		this.onBaseInfo(evt)
	case EVT_GOT_CONTACT:
		this.onBaseContacts(evt)
	case EVT_RAW_MESSAGE:
		msgs := parseMessages(evt.Args[0])
		for _, msg := range msgs {
			emsg, _ := json.Marshal(msg)
			wxh.eqch <- newEvent(EVT_GOT_MESSAGE, string(emsg))
		}
	case EVT_GOT_MESSAGE: // don't here

	}
}

func (this *innerSession) onBaseInfo(evt *Event) {
	users := ParseWXInitData(evt.Args[0])
	this.users = users
	this.users = append(this.users, this.me)

	for _, u := range users {
		this.usersm[u.UserName] = u
		this.usersm[u.NickName] = u
	}
}

func (this *innerSession) onBaseContacts(evt *Event) {
	users := ParseContactData(evt.Args[0])
	for _, u := range users {
		if _, ok := this.usersm[u.UserName]; !ok {
			this.users = append(this.users, u)
			this.usersm[u.UserName] = u
			this.usersm[u.NickName] = u
		}
	}
}
