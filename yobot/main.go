package main

import (
	"flag"
	"log"
	"syscall"
	// "strings"
	"os"
	"os/signal"
	"runtime"
	debug1 "runtime/debug"
	"time"

	"go-pprofui"

	"github.com/emirpasic/gods/maps/hashbidimap"
	"github.com/fluffle/goirc/logging/glog"
	"github.com/kitech/colog"
)

var debug bool
var pxyurl string
var pprof bool

const (
	ltracep   = "trace: "
	ldebugp   = "debug: "
	linfop    = "info: "
	lwarningp = "warning: "
	lerrorp   = "error: "
	lalertp   = "alert: "
)

func init() {
	flag.BoolVar(&debug, "debug", debug, "purple debug switch")
	flag.StringVar(&pxyurl, "proxy", pxyurl, "proxy, http://")
	flag.BoolVar(&pprof, "pprof", pprof, "enable net/http/pprof: *:6060")

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
	msgbus *MsgBusClient
}

func (this *Context) sendBusEvent(e *Event) bool {
	sendok := true
	defer func() {
		if x := recover(); x != nil {
			sendok = false
			log.Printf("wow: %v", x)
		}
	}()

	select {
	case this.busch <- e:
	default:
		log.Println("send busch blocked:", len(this.busch))
		// TODO 这种情况是为什么呢，应该怎么办呢？
		debug1.PrintStack()
	}
	return sendok
}

var ctx *Context

// ./bot -debug -v 2 -logtostderr
func main() {
	flag.Parse()
	glog.Init()

	log.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	if true {
		//	go func() { log.Println(http.ListenAndServe(":6060", nil)) }()
		go func() { pprofui.Main(":6060") }()
	}

	ctx = &Context{}
	ctx.busch = make(chan *Event, MAX_BUS_QUEUE_LEN)
	ctx.acpool = NewAccountPool()
	ctx.toxagt = NewToxAgent()
	ctx.toxagt.start()
	ctx.rtab = NewRoundTable()
	ctx.msgbus = newMsgBusClient()

	go ctx.rtab.run()

	shutdownHandler := func() {
		ctx.rtab.stop()
		<-ctx.rtab.done()
		ctx.acpool.disconnectAll()
		ctx.toxagt.stop()
		log.Println("shutdown done.")
		os.Exit(0)
	}

	// TODO system signal, elegant shutdown
	elegantShutdown := func(hfunc func()) {
		var niceCloseC = make(chan os.Signal, 0)
		signal.Notify(niceCloseC, os.Interrupt, syscall.SIGPIPE)
		intrTimes := 0
		for {
			select {
			case sig := <-niceCloseC:
				log.Println("received sig:", sig.String())
				switch sig {
				case os.Interrupt:
					intrTimes += 1
					if intrTimes > 1 {
						log.Println("force shutdown...")
						os.Exit(0)
						return
					}
					hfunc()
					goto endfor
				case syscall.SIGPIPE:
					// 为啥并没有捕捉到这个信号，程序依旧崩溃了，是因为在gdb中吗
					// 果然是因为gdb的问题吗，测试SIGINT捕捉不到
					log.Println("wow, SIGPIPE occurs. omit.")
				}
			} // end select
		} //end for

	endfor:
		return
	}

	elegantShutdown(shutdownHandler)
}

// TODO multiple servers,
//const serverssl = "weber.freenode.net:6697"
const serverssl = "irc.freenode.net:6697"
const toxname = "zuck05l" // hlpbot
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
