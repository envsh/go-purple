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
	busch  chan interface{}
	toxagt *ToxAgent
	acpool *AccountPool
	rtab   *RoundTable
}

var ctx *Context

func main() {
	flag.Parse()

	ctx = &Context{}
	ctx.busch = make(chan interface{}, 123)
	ctx.acpool = NewAccountPool()
	ctx.toxagt = NewToxAgent()
	ctx.toxagt.start()
	ctx.rtab = NewRoundTable()

	ctx.rtab.run()

	// TODO system signal, elegant shutdown
}

const serverssl = "weber.freenode.net:6697"
const toxname = "zuck07"
const ircname = toxname

var chmap = hashbidimap.New()

func init() {
	// irc <=> tox
	chmap.Put("#tox-cn123", "testks")
	chmap.Put("#tox-cn", "Chinese 中文")
	chmap.Put("#tox-en", "#tox")
}
