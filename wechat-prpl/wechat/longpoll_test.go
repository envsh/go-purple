package wechat

import (
	"log"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	t.Run("uuid", func(t *testing.T) {
		log.Println(123)
		str := `window.QRLogin.code = 200; window.QRLogin.uuid = "gYmgd1grLg==";`
		log.Println(str)
		code, uuid := parseuuid(str)
		if code != 200 || len(uuid) == 0 {
			t.Fail()
		}
		str2 := `window.QRLogin.code = 200; window.QRLogin.uuid = "wdXyNY-p5g==";`
		code, uuid = parseuuid(str2)
		if code != 200 || len(uuid) == 0 {
			t.Fail()
		}
	})
	t.Run("ticket", func(t *testing.T) {
		str := `<error><ret>0</ret><message>OK</message><skey>@crypt_3ea2fe08_723d1e1bd7b4171657b58c6d2849b367</skey><wxsid>9qxNHGgi9VP4/Tx6</wxsid><wxuin>979270107</wxuin><pass_ticket>%2BEdqKi12tfvM8ZZTdNeh4GLO9LFfwKLQRpqWk8LRYVWFkDE6%2FZJJXurz79ARX%2FIT</pass_ticket><isgrayscale>1</isgrayscale></error>`
		parseTicket(str)
		str2 := `<error><ret>0</ret><message></message><skey>@crypt_3ea2fe08_fd75b74965cfcdc0202733f6a11485b1</skey><wxsid>MAzn6b0ZO2e7J+qa</wxsid><wxuin>979270107</wxuin><pass_ticket>r%2F5XCAHISKCRB87xpK52Pqf3n3umRTOSfKxMy%2BkXX4GQrRNbP0SFuP4TOwiGDVSh</pass_ticket><isgrayscale>1</isgrayscale></error>`
		parseTicket(str2)
		t.Fatal("heheh")
	})
	t.Run("netpoll", func(t *testing.T) {
		// t.SkipNow()
		wx := NewWechat()
		wx.OnEvent = func(evt *Event) {
			if len(evt.Args) == 0 {
				log.Println(evt.Type, len(evt.Args))
			} else {
				if len(evt.Args[0]) >= 68 {
					log.Println(evt.Type, len(evt.Args), len(evt.Args[0]), evt.Args[0][0:68])
				} else {
					log.Println(evt.Type, len(evt.Args), len(evt.Args[0]), evt.Args[0])
				}
			}
		}

		for idx := 0; idx < 3000; idx++ {
			wx.Iterate()
			time.Sleep(300 * time.Millisecond)
		}
	})
	t.Run("load cookies", func(t *testing.T) {
		t.SkipNow()
		wx := NewWechat()
		wx.poller.loadCookies()
	})
}
