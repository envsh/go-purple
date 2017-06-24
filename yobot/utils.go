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
