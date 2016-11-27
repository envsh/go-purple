package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-purple/purple"

	"github.com/kitech/colog"
	"github.com/kitech/go-toxcore"
	"github.com/thoj/go-ircevent"
)

var username string = "yournicknameu@weber.freenode.net"
var debug bool

func init() {
	flag.StringVar(&username, "u", username, "your username of irc")
	flag.BoolVar(&debug, "debug", debug, "purple debug switch")
	colog.Register()
	colog.SetFlags(log.Flags() | log.Lshortfile | log.LstdFlags)
	time.Sleep(0)
}

var gbot *Yobot
var gtox *tox.Tox

func main() {
	flag.Parse()

	// TODO system signal, elegant shutdown
	bot := NewYobot()
	gbot = bot
	bot.init()
	bot.run()
}

type Yobot struct {
	pc    *purple.PurpleCore
	ctrl  *Controller
	am    *AccountManager
	ircon *irc.Connection
	acp   *AccountPool
}

func NewYobot() *Yobot {
	this := &Yobot{}
	return this
}

const serverssl = "weber.freenode.net:6697"
const toxname = "zuck07"
const ircname = toxname

var chanMap = map[string]string{
	"testks": "#tox-cn123", "Chinese 中文": "#tox-cn",
	"#tox": "tox-en",
}

func (this *Yobot) init() {
	this.acp = NewAccountPool()
	this.acp.add(ircname)

	this.setupTox()
	this.save_account()
	go this.Iterate()

	gtox.CallbackSelfConnectionStatus(func(t *tox.Tox, status uint32, d interface{}) {
		log.Println(status)
		fn, err := t.FriendByPublicKey(groupbot)
		log.Println(fn, err)
		if err != nil {
			t.FriendAdd(groupbot, "me here")
		}
		this.save_account()
	}, nil)
	gtox.CallbackFriendRequest(func(t *tox.Tox, pubkey, msg string, d interface{}) {
		log.Println("hehhe", pubkey, msg)
		friendNumber, err := t.FriendAddNorequest(pubkey)
		if err != nil {
			log.Println(err, friendNumber)
		}
		this.save_account()
	}, nil)

	gtox.CallbackFriendConnectionStatus(func(t *tox.Tox, friendNumber uint32, status uint32, d interface{}) {
		log.Println(friendNumber, status)
		this.save_account()
		pubkey, err := gtox.FriendGetPublicKey(friendNumber)
		if err != nil {
			log.Println(err, pubkey)
		}

		defer func() {
			if strings.HasPrefix(groupbot, pubkey) {
				// t.FriendSendMessage(friendNumber, "invite 1")
				// t.FriendSendMessage(friendNumber, "invite 2")
				_, err := t.FriendSendMessage(friendNumber, "invite 5")
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}, nil)
	gtox.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
		log.Println(friendNumber, msg, purple.MyTid2())
	}, nil)
	gtox.CallbackGroupMessage(func(t *tox.Tox, groupNumber int,
		peerNumber int, message string, d interface{}) {
		log.Println(groupNumber, peerNumber, message)
		groupTitle, err := t.GroupGetTitle(groupNumber)
		if err != nil {
			log.Println(err, groupTitle)
		}
		pubkeys := t.GroupGetPeerPubkeys(groupNumber)
		log.Println(pubkeys)
		groupbotIn := false
		for _, pubkey := range pubkeys {
			if strings.HasPrefix(groupbot, pubkey) {
				groupbotIn = true
			}
		}
		selfMessage := false
		peerPubkey, err := t.GroupPeerPubkey(groupNumber, peerNumber)
		if strings.HasPrefix(t.SelfGetAddress(), peerPubkey) {
			selfMessage = true
		}
		if selfMessage {
			return
		}
		// should be
		if groupbotIn {
			if toname, ok := chanMap[groupTitle]; ok {
				// forward message to...
				this.acp.get(ircname).ircon.Join(toname)
				this.acp.get(ircname).ircon.Privmsg(toname, message)
			} else {
				log.Println("unsupported group:", groupTitle)
			}
		} else {
			// forward message to...
			this.acp.get(ircname).ircon.Join(groupTitle)
			this.acp.get(ircname).ircon.Privmsg(groupTitle, message)
		}
	}, nil)

	gtox.CallbackGroupInvite(func(t *tox.Tox,
		friendNumber uint32, itype uint8, data []byte, d interface{}) {
		log.Println(friendNumber, len(data), itype)
		pubkey, err := gtox.FriendGetPublicKey(friendNumber)
		if err != nil {
			log.Println(err, pubkey)
		}

		acceptInvite := func(interface{}) {
			var groupNumber int
			var err error
			switch itype {
			case tox.GROUPCHAT_TYPE_AV:
				groupNumber, err = t.JoinAVGroupChat(friendNumber, data)
				if err != nil {
					log.Println(err, groupNumber)
				}
			case tox.GROUPCHAT_TYPE_TEXT:
				groupNumber, err = t.JoinGroupChat(friendNumber, data)
				if err != nil {
					log.Println(err, groupNumber)
				}
			default:
				log.Panicln("wtf")
			}
			if err == nil {
				// 立即取Title一般会失败的
				groupTitle, err := t.GroupGetTitle(groupNumber)
				if err != nil {
					log.Println(err, groupTitle)
				}
			}
		}
		if strings.HasPrefix(groupbot, pubkey) {
			acceptInvite(nil)
		} else if strings.HasPrefix(pubkey, "398C8") {
			acceptInvite(nil)
		}
	}, nil)

	gtox.CallbackGroupTitle(func(t *tox.Tox,
		groupNumber int, peerNumber int, title string, d interface{}) {
		log.Println(groupNumber, peerNumber, title)
	}, nil)

	if false {
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
}

func (this *Yobot) run() {
	// go this.ctrl.serve()
	// this.pc.MainLoop()
	go this.handleIrcEvent()
	select {}
}

var chanMap2 = map[string]string{
	"#tox-cn123": "testks", "#tox-cn": "Chinese 中文", "#tox-en": "#Tox",
}

func (this *Yobot) getToxGroupByName(name string) int {
	chats := gtox.GetChatList()
	log.Println(len(chats), chats)
	for _, groupNumber := range chats {
		groupTitle, err := gtox.GroupGetTitle(int(groupNumber))
		if err != nil {
			log.Println(err, groupNumber, groupTitle)
		} else {
			if groupTitle == name {
				return int(groupNumber)
			}
		}
	}
	return -1
}
func (this *Yobot) handleIrcEvent() {
	for ie := range busch {
		e := ie.(*irc.Event)
		log.Println(e)
		ircon := e.Connection
		switch e.Code {
		case "376": // MOTD end
			ircon.Join("#tox-cn123")
		case "PING":
			log.Println(e)
		case "PRIVMSG":
			log.Println(e)
			chname := e.Arguments[0]
			message := e.Arguments[1]

			var toname string = chname
			if toname_, ok := chanMap2[chname]; ok {
				toname = toname_
			}

			groupNumber := this.getToxGroupByName(toname)
			if groupNumber == -1 {
				log.Println("group not exists:", toname)
			} else {
				_, err := gtox.GroupMessageSend(groupNumber, message)
				if err != nil {
					log.Println(err, toname, groupNumber, message)
				}
			}
		case "JOIN":
			log.Println(e)

		}
	}
}

var bsnodes = []string{
	"biribiri.org", "33445", "F404ABAA1C99A9D37D61AB54898F56793E1DEF8BD46B1038B9D822E8460FAB67",
	"178.62.250.138", "33445", "788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B",
	"205.185.116.116", "33445", "A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702",
}

var groupbot = "56A1ADE4B65B86BCD51CC73E2CD4E542179F47959FE3E0E21B4B0ACDADE51855D34D34D37CB5"

func (this *Yobot) setupTox() {
	toxops := tox.NewToxOptions()
	gtox = tox.NewTox(toxops)

	toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
	this.load_account(toxops)

	// retry 50 times
	for port := 0; port < 50; port++ {
		toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
		gtox = tox.NewTox(toxops)
		if gtox != nil {
			log.Println("TOXID:", gtox.SelfGetAddress())
			break
		}
	}
	if gtox == nil {
		log.Panicln("null")
	}

	gtox.SelfSetName(toxname)

	for i := 0; i < len(bsnodes); i += 3 {
		port, _ := strconv.Atoi(bsnodes[i+1])
		ok1, err1 := gtox.Bootstrap(bsnodes[i], uint16(port), bsnodes[i+2])
		ok2, err2 := gtox.AddTcpRelay(bsnodes[i], uint16(port), bsnodes[i+2])
		if !ok1 || !ok2 || err1 != nil || err2 != nil {
			log.Println(ok1, ok2, err1, err2)
		}
	}

}

func (this *Yobot) Iterate() {
	stopped := false
	tick := time.Tick(100 * time.Millisecond)
	id := gtox.SelfGetAddress()
	for !stopped {
		select {
		case <-tick:
			gtox.Iterate()
		}
	}
	log.Println("stopped", id)
}

var tox_save_file = "./tox.save"

func (this *Yobot) load_account(toxops *tox.ToxOptions) {
	data, err := ioutil.ReadFile(tox_save_file)
	if err != nil || len(data) == 0 {
		log.Println("load data error:", err)
	} else {
		toxops.Savedata_data = data
		toxops.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
	}
}

func (this *Yobot) save_account() {
	data := gtox.GetSavedata()
	ioutil.WriteFile(tox_save_file, data, 0644)
}
