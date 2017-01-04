package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	// "sync"
	"time"

	"github.com/kitech/go-toxcore"
	"github.com/sasha-s/go-deadlock"
)

type ToxAgent struct {
	RelaxCallObject
	ctx    *Context
	_tox   *tox.Tox
	_toxmu deadlock.Mutex // sync.Mutex

	groupMembers map[int]map[int][]string // leave group get peer name/pubkey
	theirGroups  map[int]bool             // accepted group number => true
}

func NewToxAgent() *ToxAgent {
	this := &ToxAgent{}
	this.ctx = ctx

	this.groupMembers = make(map[int]map[int][]string)
	this.theirGroups = make(map[int]bool)

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

	this._tox.CallbackSelfConnectionStatus(func(t *tox.Tox, status int, d interface{}) {
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

	this._tox.CallbackFriendConnectionStatus(this.onFriendConnectionStatus, nil)
	this._tox.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
		log.Println(friendNumber, msg)
		if this.isGroupbot(friendNumber) {
			// maybe need skip the message here
		}
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_MESSAGE, "", msg, friendNumber)
	}, nil)
	this._tox.CallbackGroupMessage(this.onGroupMessage, nil)
	this._tox.CallbackGroupInvite(this.onGroupInvite, nil)

	this._tox.CallbackGroupTitle(func(t *tox.Tox,
		groupNumber int, peerNumber int, title string, d interface{}) {
		log.Println(groupNumber, peerNumber, title)
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber)
	}, nil)

	this._tox.CallbackGroupAction(this.onGroupAction, nil)
	this._tox.CallbackGroupNameListChange(this.onGroupNameListChange, nil)
}

func (this *ToxAgent) onFriendConnectionStatus(t *tox.Tox, friendNumber uint32, status int, d interface{}) {
	log.Println(friendNumber, status)
	this.save_account()
	pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
	if err != nil {
		log.Println(err, pubkey)
	}

	defer func() {
		if status == 0 {
			return
		}
		invcmds := []string{"0", "2", "3", "5"}
		if strings.HasPrefix(groupbot, pubkey) && status > 0 {
			// t.FriendSendMessage(friendNumber, "invite 1")
			// t.FriendSendMessage(friendNumber, "invite 2")
			for idx := 0; idx < len(invcmds); idx++ {
				if invcmds[idx] != "5" {
					continue
				}
				cmd := "invite " + invcmds[idx]
				_, err := t.FriendSendMessage(friendNumber, cmd)
				if err != nil {
					log.Println(err)
				}
				log.Println("send groupbot invite:", friendNumber, status, err, cmd)
			}
		}
	}()
	if status > 0 {
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_CONNECTED, "", friendNumber, status)
	} else {
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_FRIEND_DISCONNECTED, "", friendNumber, status)
	}
}

func (this *ToxAgent) onGroupMessage(t *tox.Tox, groupNumber int,
	peerNumber int, message string, d interface{}) {
	log.Println(groupNumber, peerNumber, message)
	groupTitle, err := t.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println(err, groupTitle)
	}
	pubkeys := t.GroupGetPeerPubkeys(groupNumber)
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
		// log.Println("omit self message forward", groupTitle)
		return
	}
	if len(this.ctx.busch) >= MAX_BUS_QUEUE_LEN {
		log.Println("busch full, will blocked")
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
}

func (this *ToxAgent) onGroupInvite(t *tox.Tox,
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
			this.theirGroups[groupNumber] = true
		}
	}
	if strings.HasPrefix(groupbot, pubkey) {
		acceptInvite(nil)
	} else if strings.HasPrefix(pubkey, "398C8") {
		acceptInvite(nil)
	}

}

func (this *ToxAgent) onGroupAction(t *tox.Tox,
	groupNumber int, peerNumber int, message string, userData interface{}) {
	log.Println(groupNumber, peerNumber, message)
	groupTitle, err := t.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println(err, groupTitle)
	}
	pubkeys := t.GroupGetPeerPubkeys(groupNumber)
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
		// log.Println("omit self message forward", groupTitle)
		return
	}
	if len(this.ctx.busch) >= MAX_BUS_QUEUE_LEN {
		log.Println("busch full, will blocked")
	}
	this.ctx.busch <- NewEvent(PROTO_TOX, EVT_GROUP_ACTION, groupTitle,
		message, groupNumber, peerNumber)

	// should be
	if groupbotIn {
	}
}

func (this *ToxAgent) onGroupNameListChange(t *tox.Tox,
	groupNumber int, peerNumber int, change uint8, ud interface{}) {
	groupTitle, err := this._tox.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println("wtf", err)
	}
	peerName, err := this._tox.GroupPeerName(groupNumber, peerNumber)
	if err != nil {
		if change != tox.CHAT_CHANGE_PEER_DEL {
			log.Println("wtf", err)
		}
	}
	var peerPubkey string

	switch change {
	case tox.CHAT_CHANGE_PEER_DEL:
		if _, ok := this.groupMembers[groupNumber]; !ok {
			log.Println("wtf")
		}
		peerInfo := this.groupMembers[groupNumber][peerNumber]
		if peerInfo == nil || len(peerInfo) != 2 {
			log.Println("wtf", peerInfo, peerName, groupNumber)
			break
		}
		peerName, peerPubkey = peerInfo[0], peerInfo[1]
		this.ctx.busch <- NewEvent(PROTO_TOX, EVT_LEAVE_GROUP,
			groupTitle, peerName, peerPubkey, groupNumber, peerNumber, change)

		delete(this.groupMembers[groupNumber], peerNumber)

	case tox.CHAT_CHANGE_PEER_ADD:
		peerPubkey, err := this._tox.GroupPeerPubkey(groupNumber, peerNumber)
		if err != nil {
			log.Println(err)
		}
		if _, ok := this.groupMembers[groupNumber]; !ok {
			this.groupMembers[groupNumber] = make(map[int][]string)
		}
		this.groupMembers[groupNumber][peerNumber] = []string{peerName, peerPubkey}

	case tox.CHAT_CHANGE_PEER_NAME:
		peerPubkey, err := this._tox.GroupPeerPubkey(groupNumber, peerNumber)
		if err != nil {
			log.Println(err)
		}
		if _, ok := this.groupMembers[groupNumber]; !ok {
			this.groupMembers[groupNumber] = make(map[int][]string)
		}
		this.groupMembers[groupNumber][peerNumber] = []string{peerName, peerPubkey}
	}

	// check only me left case
	if change == tox.CHAT_CHANGE_PEER_DEL {
		if pn := this._tox.GroupNumberPeers(groupNumber); pn == 1 {
			log.Println("oh, only me left:", groupNumber, groupTitle)
			// check our create group or not
			if _, ok := this.theirGroups[groupNumber]; ok {
				log.Println("their invite group matched, clean it", groupNumber, groupTitle)
				delete(this.theirGroups, groupNumber)
				this._tox.DelGroupChat(groupNumber)
			}
		}
	}
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

	toxops.Udp_enabled = false
	toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
	this.load_account(toxops)

	// retry 50 times
	for port := 0; port < 50; port++ {
		toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
		this._tox = tox.NewTox(toxops)
		if this._tox != nil {
			log.Println("TOXID:", this._tox.SelfGetAddress(), port)
			break
		}
	}
	if this._tox == nil {
		log.Panicln("null")
	}

	this._tox.SelfSetName(toxname)
	this._tox.SelfSetStatusMessage(statusMessage)

	for i := 0; i < len(bsnodes); i += 3 {
		port, _ := strconv.Atoi(bsnodes[i+1])
		ok1, err1 := this._tox.Bootstrap(bsnodes[i], uint16(port), bsnodes[i+2])
		ok2, err2 := this._tox.AddTcpRelay(bsnodes[i], uint16(port), bsnodes[i+2])
		if !ok1 || !ok2 || err1 != nil || err2 != nil {
			log.Println(ok1, ok2, err1, err2)
		}
	}

}

// TODO multiple result and reverse order search,
// for use new group, but not old unsable group
func (this *ToxAgent) getToxGroupByName(name string) int {
	chats := this._tox.GetChatList()
	log.Println(len(chats), chats, name)
	for idx, groupNumber := range chats {
		// reverse order
		groupNumber = chats[len(chats)-1-idx]
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

func (this *ToxAgent) getToxGroupNames() map[int]string {
	ret := make(map[int]string, 0)
	chats := this._tox.GetChatList()
	for idx, groupNumber := range chats {
		// reverse order
		groupNumber = chats[len(chats)-1-idx]
		groupTitle, err := this._tox.GroupGetTitle(int(groupNumber))
		if err != nil {
			log.Println(err, groupNumber, groupTitle)
		} else {
			// ret = append(ret, groupTitle)
			ret[int(groupNumber)] = groupTitle
		}
	}
	return ret
}

func (this *ToxAgent) Iterate() {
	stopped := false
	tick := time.Tick(100 * time.Millisecond)
	id := this._tox.SelfGetAddress()
	for !stopped {
		select {
		case <-tick:
			// this.Call0(func() { this._tox.Iterate() })
			this._toxmu.Lock()
			this._tox.Iterate()
			this._toxmu.Unlock()
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

func (this *ToxAgent) getOnlineFriendCount() int {
	onlineFriendCount := 0
	friendNumbers := this._tox.SelfGetFriendList()
	for _, friendNumber := range friendNumbers {
		cs, err := this.ctx.toxagt._tox.FriendGetConnectionStatus(friendNumber)
		if err != nil {
			log.Println(err)
		} else {
			if cs > tox.CONNECTION_NONE {
				onlineFriendCount += 1
			}
		}
	}
	return onlineFriendCount
}

func (this *ToxAgent) isGroupbot(friendNumber uint32) bool {
	botfn, err := this._tox.FriendByPublicKey(groupbot)
	if err != nil {
		log.Println(err)
	}
	return botfn == friendNumber
}
