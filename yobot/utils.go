package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func isLocalUrl(u string) bool {
	uo, err := url.Parse(u)
	if err != nil {
		log.Println(err)
	} else {
		hostAndPort := strings.Split(uo.Host, ":")
		if len(hostAndPort) > 0 {
			if hostAndPort[0] == "127.0.0.1" || hostAndPort[0] == "localhost" {

				return true
			}
		}
	}
	return false
}

// d just an optional parameter here
func sendChanTimeouted(c interface{}, v interface{}, d ...time.Duration) bool {
	cty := reflect.TypeOf(c)
	if cty.Kind() != reflect.Chan {
		return false
	}
	vty := reflect.TypeOf(v)
	if vty.AssignableTo(cty.Elem()) || vty.ConvertibleTo(cty.Elem()) {
	} else {
		return false
	}

	rd := 10 * time.Second
	if len(d) > 0 {
		rd = d[0]
	}

	sendok := true
	defer func() {
		if x := recover(); x != nil {
			log.Printf("wow should be closed channel: %v", x)
			sendok = false
		}
	}()

	cv := reflect.ValueOf(c)
	vv := reflect.ValueOf(v)
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectSend, Chan: cv, Send: vv},
		{Dir: reflect.SelectDefault},
	}
	chosen, _, _ := reflect.Select(cases)
	if chosen == 1 {
		log.Println("send busch blocked:", cv.Len())
		// 这种情况是为什么呢，应该怎么办呢？
		// 這種情況下，是chan寫滿了，所以block了。
		// debug1.PrintStack()
		ctx, ccfn := context.WithTimeout(context.Background(), rd)
		defer ccfn()
		func(ctx context.Context) {
			cases := []reflect.SelectCase{
				{Dir: reflect.SelectSend, Chan: cv, Send: vv},
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
			}
			chosen, _, _ := reflect.Select(cases)
			if chosen == 1 {
				log.Println(ctx.Err())
				msg := fmt.Sprintf("send busch timeout: %d, dropped", cv.Len())
				log.Println(msg)
			}
		}(ctx) // 使用一个函数，明确ctx的使用
	}

	return sendok
}

type emojiRange struct {
	Begin rune
	End   rune
}

func isEmojiChar(c rune) bool {
	for _, r := range emojiTable {
		if c >= r.Begin && c <= r.End {
			return true
		}
	}
	return false
}

var emojiTable = []emojiRange{{0x0080, 0x02AF},
	{0x0300, 0x03FF},
	{0x0600, 0x06FF},
	{0x0C00, 0x0C7F},
	{0x1DC0, 0x1DFF},
	{0x1E00, 0x1EFF},
	{0x2000, 0x209F},
	{0x20D0, 0x214F},
	{0x2190, 0x23FF},
	{0x2460, 0x25FF},
	{0x2600, 0x27EF},
	{0x2900, 0x29FF},
	{0x2B00, 0x2BFF},
	{0x2C60, 0x2C7F},
	{0x2E00, 0x2E7F},
	{0x3000, 0x303F},
	{0xA490, 0xA4CF},
	{0xE000, 0xF8FF},
	{0xFE00, 0xFE0F},
	{0xFE30, 0xFE4F},
	{0x1F000, 0x1F02F},
	{0x1F0A0, 0x1F0FF},
	{0x1F100, 0x1F64F},
	{0x1F680, 0x1F6FF},
	{0x1F910, 0x1F96B},
	{0x1F980, 0x1F9E0},
}
