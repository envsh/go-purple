package main

import (
	"flag"
	"log"
	"time"

	"go-purple/purple"

	"github.com/kitech/colog"
)

var username string = "yournicknameu@weber.freenode.net"
var debug bool

func init() {
	flag.StringVar(&username, "u", username, "your username of irc")
	flag.BoolVar(&debug, "debug", debug, "purple debug switch")
	colog.Register()
	log.SetFlags(log.Flags() | log.Lshortfile | log.LstdFlags)
	time.Sleep(0)
}

var gbot *Yobot

func main() {
	flag.Parse()

	bot := NewYobot()
	gbot = bot
	bot.init()
	bot.run()
}

type Yobot struct {
	pc   *purple.PurpleCore
	ctrl *Controller
	am   *AccountManager
}

func NewYobot() *Yobot {
	this := &Yobot{}
	return this
}

func (this *Yobot) init() {
	if debug {
		purple.DebugSetEnabled(true)
	}
	purple.UtilSetUserDir(purple.UserDir() + "-yobot")

	/////
	this.pc = purple.NewPurpleCore()

	this.ctrl = NewController()
	this.ctrl.init()

	this.am = NewAccountManager()
	this.am.init()
}

func (this *Yobot) run() {
	go this.ctrl.serve()
	this.pc.MainLoop()
}
