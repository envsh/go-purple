package main

import (
	"flag"
	"log"

	"github.com/kitech/colog"

	"yobot/purple"
)

var username string = "yournicknameu@irc.freenode.net"

func init() {
	colog.Register()
	colog.SetFlags(log.Lshortfile | log.LstdFlags | colog.Flags())

	flag.StringVar(&username, "u", username, "your username of irc")
}

func main() {
	log.Println("good")

	pc := purple.NewPurpleCore()
	pc.InitPurple()

	acc := pc.AccountsFind(username, "prpl-irc")
	if acc == nil {
		acc = purple.NewAccountCreate(username, "prpl-irc", "")
		log.Println(acc)
	}
	acc.SetEnabled(true)
	acc.Connect()
	// pc.ToRoom(acc)

	go func() {
		pc.Loop()
	}()

	log.Println(pc)
	select {}
}
