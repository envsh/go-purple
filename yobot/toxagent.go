package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	// "sync"
	"runtime"
	"time"

	"github.com/kitech/go-toxcore"
)

type ToxAgent struct {
	RelaxCallObject
	ctx       *Context
	_tox      *tox.Tox
	stopC     chan bool
	reconnC   chan time.Duration
	delGroupC chan int

	groupMembers map[int]map[int][]string // leave group get peer name/pubkey
	theirGroups  map[int]bool             // accepted group number => true

}

func NewToxAgent() *ToxAgent {
	this := &ToxAgent{}
	this.ctx = ctx
	this.stopC = make(chan bool, 0)
	this.reconnC = make(chan time.Duration, 1)
	this.delGroupC = make(chan int, 16)

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
	this.save_account()
	this.stopC <- true
	runtime.Gosched()
	this._tox.Kill()
}

func (this *ToxAgent) setupCallbacks() {

	this._tox.CallbackSelfConnectionStatus(func(t *tox.Tox, status int, d interface{}) {
		log.Println(status, tox.ConnStatusString(status))
		friendNumber, err := t.FriendByPublicKey(groupbot)
		if err != nil && status > tox.CONNECTION_NONE {
			t.FriendAdd(groupbot, fmt.Sprintf("Hey %d, me here", friendNumber))
		}
		t.SelfSetName(t.SelfGetName()) // deadlock test
		// TODO fixed enter groupbot group

		this.save_account()
		if status > tox.CONNECTION_NONE {
			this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_CONNECTED, "", status))
		} else {
			this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_DISCONNECTED, "", status))
		}
		if status == tox.CONNECTION_NONE {
			// this.tryReconnect(3 * time.Second)
			this.reconnC <- (3 * time.Second)
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
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_FRIEND_MESSAGE, "", msg, friendNumber))
	}, nil)
	this._tox.CallbackGroupMessage(this.onGroupMessage, nil)
	this._tox.CallbackGroupInvite(this.onGroupInvite, nil)

	this._tox.CallbackGroupTitle(func(t *tox.Tox,
		groupNumber int, peerNumber int, title string, d interface{}) {
		log.Println(groupNumber, peerNumber, title)
		if strings.HasPrefix(title, DELETED_INVITED_GROUPCHAT_P) {
			// TODO do what?
			return
		}
		peerName, err := this._tox.GroupPeerName(groupNumber, peerNumber)
		if err != nil {
			log.Println(lerrorp, err)
		}
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_JOIN_GROUP, title, groupNumber, peerNumber, peerName))
	}, nil)

	this._tox.CallbackGroupAction(this.onGroupAction, nil)
	this._tox.CallbackGroupNameListChange(this.onGroupNameListChange, nil)
}

const DELETED_INVITED_GROUPCHAT_P = "#deleted_invited_groupchat_"

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
	if status > tox.CONNECTION_NONE {
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_FRIEND_CONNECTED, "", friendNumber, status))
	} else {
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_FRIEND_DISCONNECTED, "", friendNumber, status))
	}
}

// TODO save invite data, test rejoin if leave
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

func (this *ToxAgent) onGroupMessage(t *tox.Tox, groupNumber int,
	peerNumber int, message string, d interface{}) {
	log.Println(groupNumber, peerNumber, message)
	defer log.Println(groupNumber, peerNumber, message)

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
	this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_GROUP_MESSAGE, groupTitle,
		message, groupNumber, peerNumber))

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
	this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_GROUP_ACTION, groupTitle,
		message, groupNumber, peerNumber))

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
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_LEAVE_GROUP, groupTitle, peerName, peerPubkey, groupNumber, peerNumber, change))

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
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_JOIN_GROUP, groupTitle, groupNumber, peerNumber, peerName))
	case tox.CHAT_CHANGE_PEER_NAME:
		peerPubkey, err := this._tox.GroupPeerPubkey(groupNumber, peerNumber)
		if err != nil {
			log.Println(err)
		}
		if _, ok := this.groupMembers[groupNumber]; !ok {
			this.groupMembers[groupNumber] = make(map[int][]string)
		}
		this.groupMembers[groupNumber][peerNumber] = []string{peerName, peerPubkey}
		this.ctx.sendBusEvent(NewEvent(PROTO_TOX, EVT_JOIN_GROUP, groupTitle, groupNumber, peerNumber, peerName))

		// TODO change name event
		log.Println(lwarningp, "TODO", "tox rename event", peerName)
	}

	// check only me left case
	this.checkOnlyMeLeftGroup(groupNumber, peerNumber, change)
}

// check only me left case
func (this *ToxAgent) checkOnlyMeLeftGroup(groupNumber int, peerNumber int, change uint8) {
	groupTitle, err := this._tox.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println("wtf", err)
	}
	peerName, err := this._tox.GroupPeerName(groupNumber, peerNumber)
	if err != nil {
		if change != tox.CHAT_CHANGE_PEER_DEL {
			log.Println("wtf", err, peerName)
		}
	}
	// var peerPubkey string

	switch change {
	case tox.CHAT_CHANGE_PEER_DEL:
	case tox.CHAT_CHANGE_PEER_ADD:
	case tox.CHAT_CHANGE_PEER_NAME:
	}

	// check only me left case
	if change == tox.CHAT_CHANGE_PEER_DEL {
		if pn := this._tox.GroupNumberPeers(groupNumber); pn == 1 {
			log.Println("oh, only me left:", groupNumber, groupTitle)
			// check our create group or not
			// 即使不是自己创建的群组，在只剩下自己之后，也可以不删除。因为这个群的所有人就是自己了。
			// 这里找一下为什么程序会崩溃吧
			if _, ok := this.theirGroups[groupNumber]; ok {
				log.Println("their invite group matched, clean it", groupNumber, groupTitle)
				delete(this.theirGroups, groupNumber)
				grptype, err := this._tox.GroupGetType(uint32(groupNumber))
				log.Println("before delete group chat", groupNumber, grptype, err)
				switch uint8(grptype) {
				case tox.GROUPCHAT_TYPE_AV:
					// log.Println("dont delete av groupchat for a try", groupNumber, ok, err)
				case tox.GROUPCHAT_TYPE_TEXT:
					// ok, err := this._tox.DelGroupChat(groupNumber)
					// log.Println("after delete group chat", groupNumber, ok, err)
				default:
					log.Fatal("wtf")
				}
				time.AfterFunc(1*time.Second, func() {
					this.delGroupC <- groupNumber
				})
				log.Println("hehhe----------------------------")
				this._tox.GroupSetTitle(groupNumber, fmt.Sprintf("#deleted_invited_groupchat_%s", time.Now().Format("20060102_150405")))
				log.Println("dont delete invited groupchat for a try", groupNumber, ok, err)
			}
		}
	}

}

var bsnodes = []string{
	"biribiri.org", "33445", "F404ABAA1C99A9D37D61AB54898F56793E1DEF8BD46B1038B9D822E8460FAB67",
	// "178.62.250.138", "33445", "788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B",
	"108.61.165.198", "33445", "8E7D0B859922EF569298B4D261A8CCB5FEA14FB91ED412A7603A585A25698832",
	"205.185.116.116", "33445", "A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702",
}

var groupbot = "56A1ADE4B65B86BCD51CC73E2CD4E542179F47959FE3E0E21B4B0ACDADE51855D34D34D37CB5"

func (this *ToxAgent) setupTox() {
	toxops := tox.NewToxOptions()
	toxops.ThreadSafe = true
	this._tox = tox.NewTox(toxops)

	toxops.Udp_enabled = false
	toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
	this.load_account(toxops)

	// retry 50 times
	for port := 0; port < 50; port++ {
		toxops.Tcp_port = uint16(rand.Uint32()%55536) + 10000
		this._tox = tox.NewTox(toxops)
		if this._tox != nil {
			log.Println("TOXID:", this._tox.SelfGetAddress(), port, toxops.Tcp_port)
			break
		}
	}
	if this._tox == nil {
		log.Panicln("null")
	}

	this._tox.SelfSetName(toxname)
	this._tox.SelfSetStatusMessage(statusMessage)

	this.bootstrap()
}

// maybe block seconds
func (this *ToxAgent) bootstrap() {
	for i := 0; i < len(bsnodes); i += 3 {
		port, _ := strconv.Atoi(bsnodes[i+1])
		ok1, err1 := this._tox.Bootstrap(bsnodes[i], uint16(port), bsnodes[i+2])
		ok2, err2 := this._tox.AddTcpRelay(bsnodes[i], uint16(port), bsnodes[i+2])
		if !ok1 || !ok2 || err1 != nil || err2 != nil {
			log.Println(ok1, ok2, err1, err2)
		}
	}
	log.Println("bootstrap done:", len(bsnodes)/3)
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

func (this *ToxAgent) tryReconnect(d time.Duration) {
	go func() {
		for {
			log.Println("try reconnect now")
			status := this._tox.SelfGetConnectionStatus()
			if status == tox.CONNECTION_NONE {
				this.bootstrap()
			}
			log.Println("check connection after 3 seconds...")
			time.Sleep(3 * time.Second)
			if this._tox.SelfGetConnectionStatus() > tox.CONNECTION_NONE {
				log.Println("reconnect ok")
				break
			}
		}
	}()
}

func (this *ToxAgent) tryReconnectBlock(d time.Duration) {
	log.Println("try reconnect now")
	status := this._tox.SelfGetConnectionStatus()
	if status == tox.CONNECTION_NONE {
		this.bootstrap()
	}
	log.Println("check connection after 3 seconds...")
	time.AfterFunc(3*time.Second, func() {
		if this._tox.SelfGetConnectionStatus() > tox.CONNECTION_NONE {
			log.Println("reconnect ok")
		} else {
			this.reconnC <- d
		}
	})
}

func (this *ToxAgent) Iterate() {

	tick := time.Tick(100 * time.Millisecond)
	id := this._tox.SelfGetAddress()
	for {
		select {
		case <-tick:
			// this.Call0(func() { this._tox.Iterate() })
			// this._toxmu.Lock()
			this._tox.Iterate()
			// this._toxmu.Unlock()
		case dur := <-this.reconnC:
			this.tryReconnectBlock(dur)
		case groupNumber := <-this.delGroupC:
			log.Println("before real delete group:", groupNumber)
			r, ok := this._tox.DelGroupChat(groupNumber)
			log.Println(r, ok)
			log.Println("after real delete group:", groupNumber)
		case <-this.stopC:
			goto endfor
		}
	}

endfor:
	log.Println("stopped", id)
}

var tox_save_file = "./tox.save.yobot"

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
