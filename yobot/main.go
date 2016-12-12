package main

import (
	"flag"
	"log"
	// "strings"
	"time"

	"github.com/emirpasic/gods/maps/hashbidimap"
	"github.com/kitech/colog"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", debug, "purple debug switch")
	colog.Register()
	colog.SetFlags(log.Flags() | log.Lshortfile | log.LstdFlags)
	time.Sleep(0)
}

type Context struct {
	// busch  chan interface{}
	busch  chan *Event
	toxagt *ToxAgent // it's root tox
	acpool *AccountPool
	rtab   *RoundTable
}

var ctx *Context

func main() {
	flag.Parse()

	ctx = &Context{}
	// ctx.busch = make(chan interface{}, 123)
	ctx.busch = make(chan *Event, MAX_BUS_QUEUE_LEN)
	ctx.acpool = NewAccountPool()
	ctx.toxagt = NewToxAgent()
	ctx.toxagt.start()
	ctx.rtab = NewRoundTable()

	ctx.rtab.run()

	// TODO system signal, elegant shutdown
}

// TODO multiple servers,
const serverssl = "weber.freenode.net:6697"
const toxname = "zuck07"
const ircname = toxname
const leaveChannelTimeout = 270 // seconds

var chmap = hashbidimap.New()

func init() {
	// irc <=> tox
	chmap.Put("#tox-cn123", "testks")
	chmap.Put("#tox-cn", "Chinese 中文")
	chmap.Put("#tox-en", "#tox")
	chmap.Put("#tox-ru", "Russian Tox Chat (Use Kalina: kalina@toxme.io or 12EDB939AA529641CE53830B518D6EB30241868EE0E5023C46A372363CAEC91C2C948AEFE4EB)")
}

var PREFIX_ACTION = "/me "

var statusMessage = "Send me the message 'invite', 'info', 'help' for a full list of commands"

var cmdhelp = "info : Print my current status and list active group chats\n\n" +
	"id : Print my Tox ID\n\n" +
	"invite : Request invite to default group chat\n\n" +
	"invite <n> <p> : Request invite to group chat n (with password p if protected)\n\n" +
	"group <type> <pass> : Creates a new groupchat with type: text | audio (optional password)"

var invalidcmd = "Invalid command. Type help for a list of commands"
