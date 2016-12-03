package main

import (
	"flag"
	"log"
	"time"

	"go-purple/purple"

	"github.com/kitech/colog"
)

func init() {
	colog.Register()
	colog.SetFlags(log.Flags() | log.Lshortfile | log.LstdFlags)
	time.Sleep(0)
}

func main() {
	flag.Parse()

	log.Println(purple.UserDir())
	purple.UtilSetUserDir("/home/gzleo/.purple-yobot")
	log.Println(purple.UserDir(), purple.GoID())

	pc := purple.NewPurpleCore()

	ctrl := NewController()
	ctrl.init()
	go ctrl.serve()

	as := NewAccountServer(pc)
	if false {
		go as.run()
	}

	// go func() { pc.MainLoop() }()
	pc.MainLoop()
	// select {}
}
