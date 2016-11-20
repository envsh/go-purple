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

func parseMessage(msgo *simplejson.Json) *Message {
	msg := &Message{}
	msg.MsgId = msgo.Get("MsgId").MustString()
	msg.FromUserName = msgo.Get("FromUserName").MustString()
	msg.ToUserName = msgo.Get("ToUserName").MustString()
	msg.MsgType = msgo.Get("MsgType").MustInt()
	msg.Content = msgo.Get("Content").MustString()
	msg.CreateTime = msgo.Get("CreateTime").MustUint64()
	return msg
}

func parseMessages(data string) (msgs []*Message) {
	jso, err := simplejson.NewJson([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
	retv := jso.GetPath("BaseResponse", "Ret").MustInt()
	if retv != 0 {
		log.Println("Invalid resonse:", retv)
		return
	}

	cntv := jso.GetPath("AddMsgCount").MustInt()
	log.Println("parsering msgs:", cntv)
	msgs = make([]*Message, cntv)
	for idx := 0; idx < cntv; idx++ {
		msgo := jso.Get("AddMsgList").GetIndex(idx)
		m := parseMessage(msgo)
		msgs[idx] = m
	}

	return
}

func parseMessages2(data string) (msgs []*Message) {
	p := NewParser(data)
	if !p.RetOK() {
		return
	}

	msgs = make([]*Message, 0)
	p.Each("AddMsg", func(itemo *simplejson.Json) {
		m := parseMessage(itemo)
		msgs = append(msgs, m)
	})
	log.Println("parsed msgs:", len(msgs))

	return
}
