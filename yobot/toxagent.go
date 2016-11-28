package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/kitech/go-toxcore"
)

type ToxAgent struct {
	ctx  *Context
	_tox *tox.Tox
}

func NewToxAgent() *ToxAgent {
	this := &ToxAgent{}
	this.ctx = ctx
	return this
}

func (this *ToxAgent) start() {
	this.setupTox()
	this.save_account()
	this.setupCallbacks()
	go this.Iterate()
}

func (this *ToxAgent) stop() {

}

func (this *ToxAgent) setupCallbacks() {

	this._tox.CallbackSelfConnectionStatus(func(t *tox.Tox, status uint32, d interface{}) {
		log.Println(status)
		fn, err := t.FriendByPublicKey(groupbot)
		log.Println(fn, err)
		if err != nil {
			t.FriendAdd(groupbot, "me here")
		}
		this.save_account()
		if status > 0 {
			this.ctx.busch <- NewEvent(PROTO_TOX, EVT_CONNECTED, "", status)
		} else {
			this.ctx.busch <- NewEvent(PROTO_TOX, EVT_DISCONNECTED, "", status)
		}
	}, nil)
	this._tox.CallbackFriendRequest(func(t *tox.Tox, pubkey, msg string, d interface{}) {
		log.Println("hehhe", pubkey, msg)
		friendNumber, err := t.FriendAddNorequest(pubkey)
		if err != nil {
			log.Println(err, friendNumber)
		}
		this.save_account()
	}, nil)

	this._tox.CallbackFriendConnectionStatus(func(t *tox.Tox, friendNumber uint32, status uint32, d interface{}) {
		log.Println(friendNumber, status)
		this.save_account()
		pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
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
		if status > 0 {
			this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_CONNECTED, "", friendNumber, status)
		} else {
			this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_DISCONNECTED, "", friendNumber, status)
		}
	}, nil)
	this._tox.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
		log.Println(friendNumber, msg)
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_MESSAGE, "", msg, friendNumber)
	}, nil)
	this._tox.CallbackGroupMessage(func(t *tox.Tox, groupNumber int,
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
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_GROUP_MESSAGE, groupTitle,
			message, groupNumber, peerNumber)

		// should be
		if groupbotIn {
		}
		/*
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
		*/
	}, nil)

	this._tox.CallbackGroupInvite(func(t *tox.Tox,
		friendNumber uint32, itype uint8, data []byte, d interface{}) {
		log.Println(friendNumber, len(data), itype)
		pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
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

	this._tox.CallbackGroupTitle(func(t *tox.Tox,
		groupNumber int, peerNumber int, title string, d interface{}) {
		log.Println(groupNumber, peerNumber, title)
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber)
	}, nil)

}

var bsnodes = []string{
	"biribiri.org", "33445", "F404ABAA1C99A9D37D61AB54898F56793E1DEF8BD46B1038B9D822E8460FAB67",
	"178.62.250.138", "33445", "788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B",
	"205.185.116.116", "33445", "A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702",
}

var groupbot = "56A1ADE4B65B86BCD51CC73E2CD4E542179F47959FE3E0E21B4B0ACDADE51855D34D34D37CB5"

func (this *ToxAgent) setupTox() {
	toxops := tox.NewToxOptions()
	this._tox = tox.NewTox(toxops)

	toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
	this.load_account(toxops)

	// retry 50 times
	for port := 0; port < 50; port++ {
		toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
		this._tox = tox.NewTox(toxops)
		if this._tox != nil {
			log.Println("TOXID:", this._tox.SelfGetAddress())
			break
		}
	}
	if this._tox == nil {
		log.Panicln("null")
	}

	this._tox.SelfSetName(toxname)

	for i := 0; i < len(bsnodes); i += 3 {
		port, _ := strconv.Atoi(bsnodes[i+1])
		ok1, err1 := this._tox.Bootstrap(bsnodes[i], uint16(port), bsnodes[i+2])
		ok2, err2 := this._tox.AddTcpRelay(bsnodes[i], uint16(port), bsnodes[i+2])
		if !ok1 || !ok2 || err1 != nil || err2 != nil {
			log.Println(ok1, ok2, err1, err2)
		}
	}

}

func (this *ToxAgent) getToxGroupByName(name string) int {
	chats := this._tox.GetChatList()
	log.Println(len(chats), chats)
	for _, groupNumber := range chats {
		groupTitle, err := this._tox.GroupGetTitle(int(groupNumber))
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

func (this *ToxAgent) Iterate() {
	stopped := false
	tick := time.Tick(100 * time.Millisecond)
	id := this._tox.SelfGetAddress()
	for !stopped {
		select {
		case <-tick:
			this._tox.Iterate()
		}
	}
	log.Println("stopped", id)
}

var tox_save_file = "./tox.save"

func (this *ToxAgent) load_account(toxops *tox.ToxOptions) {
	data, err := ioutil.ReadFile(tox_save_file)
	if err != nil || len(data) == 0 {
		log.Println("load data error:", err)
	} else {
		toxops.Savedata_data = data
		toxops.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
	}
}

func (this *ToxAgent) save_account() {
	data := this._tox.GetSavedata()
	ioutil.WriteFile(tox_save_file, data, 0644)
}
