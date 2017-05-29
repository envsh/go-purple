package main

import (
	"crypto/tls"
	"log"
	"strings"
	"testing"
	"time"

	irc "github.com/fluffle/goirc/client"
	"github.com/jmz331/gpinyin"
)

/*
E1214 12:21:06.106419   32294 connection.go:410] irc.recv(): read tcp 10.0.0.7:42506->162.213.39.42:6697: read: connection timed out
I1214 12:21:06.107387   32294 connection.go:517] irc.Close(): Disconnected from server.

Thread 1 "yobot.bin" received signal SIGPIPE, Broken pipe.
syscall.Syscall () at /usr/lib/go/src/syscall/asm_linux_amd64.s:27
27              CMPQ    AX, $0xfffffffffffff001

*/

var onEvent123 = func(c *irc.Conn, line *irc.Line) {
	// useing hacked goirc for deadlock
	log.Println(c, line)
}

func TestCrash(t *testing.T) {

	name := "tstgoirc"
	ircfg := irc.NewConfig(name)
	ircfg.SSL = true

	ircfg.SSLConfig = &tls.Config{ServerName: strings.Split(serverssl, ":")[0]}
	ircfg.Server = serverssl
	ircfg.NewNick = func(n string) string { return n + "^" }
	ircon := irc.Client(ircfg)
	ircon.EnableStateTracking()

	for _, cmd := range ircmds {
		if true {
			ircon.HandleFunc(cmd, onEvent123)
		}
	}

	log.Println("ready conn...")
	err := ircon.Connect()
	log.Println("conn done")
	if err != nil {
		// t.Error(err)
	}
	// t.Log(err)
	log.Println(err)

	log.Println("sleeping...")
	time.Sleep(3 * time.Second)
	err = ircon.Close()
	log.Println(err, ircon.Connected())
	select {}
}

func TestTopinyin(t *testing.T) {
	name := "美味的百合仙子"
	namepy := "meiweidebaihexianzi"
	newname := gpinyin.ConvertToPinyinString(name, "", gpinyin.PINYIN_WITHOUT_TONE)
	log.Println(name, "=>", newname)
	if newname != namepy {
		t.Error("topinyin failed")
	}
}
