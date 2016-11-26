/*
 解析消息数据
 消息数据结构
*/

package wechat

import (
	"log"

	"github.com/bitly/go-simplejson"
)

type Message struct {
	MsgId        string
	FromUserName string
	ToUserName   string
	MsgType      int
	Content      string
	CreateTime   uint64
	NewMsgId     uint64 // == uint64(MsgId)
}

func ParseMessage(msg string) *Message {
	jso, err := simplejson.NewJson([]byte(msg))
	if err != nil {
		log.Println(err)
		return nil
	}
	return parseMessageJson(jso)
}

func parseMessageJson(msgo *simplejson.Json) *Message {
	msg := &Message{}
	msg.MsgId = msgo.Get("MsgId").MustString()
	msg.FromUserName = msgo.Get("FromUserName").MustString()
	msg.ToUserName = msgo.Get("ToUserName").MustString()
	msg.MsgType = msgo.Get("MsgType").MustInt()
	msg.Content = msgo.Get("Content").MustString()
	msg.CreateTime = msgo.Get("CreateTime").MustUint64()

	msg.SwapFromTo()
	return msg
}

func parseMessages(data string) (msgs []*Message) {
	p := NewParser(data)
	if !p.RetOK() {
		return
	}

	msgs = make([]*Message, 0)
	p.Each("AddMsg", func(itemo *simplejson.Json) {
		m := parseMessageJson(itemo)
		msgs = append(msgs, m)
	})
	log.Println("parsed msgs:", len(msgs))

	return
}

// 确保ToUserName一定是对话的对方，不是me
func (this *Message) SwapFromTo() {
	me := newInnerSession().me
	log.Println(me)
	if this.ToUserName == me.UserName {
		this.ToUserName, this.FromUserName = this.FromUserName, this.ToUserName
	}
}
