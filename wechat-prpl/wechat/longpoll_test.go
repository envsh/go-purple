package wechat

import (
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	t.Run("uuid", func(t *testing.T) {
		str := `window.QRLogin.code = 200; window.QRLogin.uuid = "gYmgd1grLg==";`
		code, uuid := parseuuid(str)
		if code != 200 || len(uuid) == 0 {
			t.Error(code, uuid)
		}
		str2 := `window.QRLogin.code = 200; window.QRLogin.uuid = "wdXyNY-p5g==";`
		code, uuid = parseuuid(str2)
		if code != 200 || len(uuid) == 0 {
			t.Error(code, uuid)
		}
	})
	t.Run("ticket", func(t *testing.T) {
		str := `<error><ret>0</ret><message>OK</message><skey>@crypt_3ea2fe08_723d1e1bd7b4171657b58c6d2849b367</skey><wxsid>9qxNHGgi9VP4/Tx6</wxsid><wxuin>979270107</wxuin><pass_ticket>%2BEdqKi12tfvM8ZZTdNeh4GLO9LFfwKLQRpqWk8LRYVWFkDE6%2FZJJXurz79ARX%2FIT</pass_ticket><isgrayscale>1</isgrayscale></error>`
		ret, _, _, _, _ := parseTicket(str)
		if ret != 0 {
			t.Error(ret)
		}
		str2 := `<error><ret>0</ret><message></message><skey>@crypt_3ea2fe08_fd75b74965cfcdc0202733f6a11485b1</skey><wxsid>MAzn6b0ZO2e7J+qa</wxsid><wxuin>979270107</wxuin><pass_ticket>r%2F5XCAHISKCRB87xpK52Pqf3n3umRTOSfKxMy%2BkXX4GQrRNbP0SFuP4TOwiGDVSh</pass_ticket><isgrayscale>1</isgrayscale></error>`
		ret, _, _, _, _ = parseTicket(str2)
		if ret != 0 {
			t.Error(ret)
		}
	})
	t.Run("load cookies", func(t *testing.T) {
		t.SkipNow()
		wx := NewWechat()
		wx.poller.loadCookies()
	})
	t.Run("evt string", func(t *testing.T) {
		EVT_GOT_QRCODE.String()
	})
	t.Run("parse init data", func(t *testing.T) {
		bcc, err := ioutil.ReadFile("wxinit_fmt.json")
		if err != nil {
			log.Println(err)
		}
		users := parseWXInitData(string(bcc))
		if users == nil {
			t.Error("invalid init data")
		}
	})
	t.Run("parse contact data", func(t *testing.T) {
		bcc, err := ioutil.ReadFile("wxcontact_fmt.json")
		if err != nil {
			log.Println(err)
		}
		users := parseContactData(string(bcc))
		if users == nil {
			t.Error("invalid init data")
		}
	})
}

func TestStateMachine(t *testing.T) {
	t.Run("netpoll", func(t *testing.T) {
		t.SkipNow()
		wx := NewWechat()
		wx.OnEvent = func(evt *Event, ud interface{}) {
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
			wx.Iterate(idx)
			time.Sleep(300 * time.Millisecond)
		}
	})
}
