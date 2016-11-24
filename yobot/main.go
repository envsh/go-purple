package main

import (
	"flag"
	"log"

	"go-purple/purple"
)

var username string = "yournicknameu@irc.freenode.net"

func init() {
	flag.StringVar(&username, "u", username, "your username of irc")
}

func main() {
	flag.Parse()

	pc := purple.NewPurpleCore()
	acc := pc.AccountsFind(username, "prpl-irc")
	if acc == nil {
		acc = purple.NewAccountCreate(username, "prpl-irc", "")
		log.Println(acc)
	}
	acc.SetEnabled(true)
	acc.Connect()

	pc.MainLoop()
}
