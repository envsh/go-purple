/*
  all implemention code except must code.
*/
package main

import (
	"fmt"
	"log"

	"yobot/purple"

	"github.com/kitech/go-toxcore"
)

func (this *ToxPlugin) setupSelfInfo(ac *purple.Account) {
	if len(ac.GetAlias()) > 0 {
		this._tox.SelfSetName(ac.GetAlias())
	} else {
		this._tox.SelfSetName(ac.GetUserName())
	}
	this._tox.SelfSetStatusMessage("It's from gotox-prpl, hoho.")
}

func (this *ToxPlugin) setupCallbacks(ac *purple.Account) {
	conn := ac.GetConnection()
	this._tox.CallbackSelfConnectionStatus(func(t *tox.Tox, status uint32, d interface{}) {
		log.Println("hehhe", status)
		if status == 2 {
			conn.ConnSetState(purple.CONNECTED) // 设置为已连接状态，则好友会显示。
		} else {
			conn.ConnSetState(purple.DISCONNECTED)
		}
	}, ac)

	this._tox.CallbackFriendRequest(func(t *tox.Tox, pubkey, msg string, d interface{}) {
		log.Println("hehhe", pubkey, msg)
		this._tox.FriendAddNorequest(pubkey)
		this.save_account()
	}, ac)

	this._tox.CallbackFriendConnectionStatus(func(t *tox.Tox, friendNumber uint32, status uint32, d interface{}) {
		log.Println(friendNumber, status)
		pubkey, _ := t.FriendGetPublicKey(friendNumber)
		buddy := ac.FindBuddy(pubkey)
		if buddy == nil {
			n, _ := t.FriendGetName(friendNumber)
			log.Println("can not find buddy:", n)
		} else {
			switch status {
			case 0:
				purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_OFFLINE_STR)
			case 1:
				purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_BUSY_STR)
			case 2:
				purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
			}
		}
	}, ac)

	this._tox.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
		log.Println(friendNumber, msg, purple.MyTid2())
		conn := ac.GetConnection()
		pubkey, err := t.FriendGetPublicKey(friendNumber)
		if err != nil {
			log.Println(err)
		} else {
			conn.ServGotIM(pubkey, msg, purple.MESSAGE_RECV)
		}
	}, ac)

	this._tox.CallbackGroupNameListChange(func(t *tox.Tox, groupNumber int,
		peerNumber int, change uint8, d interface{}) {
		log.Println(groupNumber, peerNumber, change)
		conn := ac.GetConnection()
		conv := conn.ConnFindChat(groupNumber)
		chat := conv.GetChatData()
		peerName, err := t.GroupPeerName(groupNumber, peerNumber)
		if err != nil {
			log.Println(err)
		}
		peerPubkey, err := t.GroupPeerPubkey(groupNumber, peerNumber)
		log.Println(peerName, peerPubkey, chat)
		switch change {
		case tox.CHAT_CHANGE_PEER_ADD:
			chat.AddUser(peerName)
		case tox.CHAT_CHANGE_PEER_DEL:
			chat.RemoveUser(peerName)
		case tox.CHAT_CHANGE_PEER_NAME:
			chat.AddUser(peerName)
		}
		// TODO member list diff and clean, so it is member list sync
		plst := chat.GetUsers()
		tlst := t.GroupGetNames(groupNumber)
		peerCount := t.GroupNumberPeers(groupNumber)
		if len(tlst) != peerCount {
			log.Println("wtf")
		}
		if len(plst) != len(tlst) {
			log.Println("need sync names...")
			log.Println("purple list:", plst)
			log.Println("tox list:", tlst)
			for _, pname := range plst {
				found := false
				for _, tname := range tlst {
					if tname == pname {
						found = true
					}
				}
				if found == false {
					chat.RemoveUser(pname)
				}
			}
		}
	}, ac)

	this._tox.CallbackGroupMessage(func(t *tox.Tox, groupNumber int,
		peerNumber int, message string, d interface{}) {
		conn := ac.GetConnection()
		pubkey, err := t.GroupPeerPubkey(groupNumber, peerNumber)
		if err != nil {
			log.Println(err)
		} else {
			conn.ServGotChatIn(groupNumber, pubkey, purple.MESSAGE_RECV, message)
		}
	}, ac)
}

func (this *ToxPlugin) loadFriends(ac *purple.Account) {
	fns := this._tox.SelfGetFriendList()
	if fns == nil || len(fns) == 0 {
		log.Println("oh, you have 0 friends")
	}
	for _, fn := range fns {
		name, err := this._tox.FriendGetName(fn)
		pubkey, err := this._tox.FriendGetPublicKey(fn)
		if err != nil {
			log.Println(err)
		}

		buddy := ac.FindBuddy(pubkey)
		if buddy == nil {
			buddy = purple.NewBuddy(ac, pubkey, name)
			ac.AddBuddy(buddy)
			buddy.BlistAdd(nil)
		}
		// purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
	}
}

// optional callbacks
func (this *ToxPlugin) ChatInfo(gc *purple.Connection) []string {
	log.Println(gc)
	infos := []string{"ToxChannel"}
	return infos
}

func (this *ToxPlugin) ChatInfoDefaults(gc *purple.Connection, chatName string) map[string]string {
	log.Println(gc)
	return nil
}

func (this *ToxPlugin) SendIM(gc *purple.Connection, who string, msg string) int {
	log.Println(gc, who, msg)
	friendNumber, _ := this._tox.FriendByPublicKey(who)
	len, err := this._tox.FriendSendMessage(friendNumber, msg)
	if err != nil {
		log.Println(err, len)
		return -1
	}
	return int(len)
}

func (this *ToxPlugin) JoinChat(gc *purple.Connection, comp *purple.GHashTable) {
	groupNumber, err := this._tox.AddGroupChat()
	if err != nil {
		log.Println(err)
	}
	title := comp.Lookup("ToxChannel")
	this._tox.GroupSetTitle(groupNumber, title)
	comp.Insert("GroupNumber", fmt.Sprintf("%d", groupNumber))
	conv := gc.ServGotJoinedChat(groupNumber, comp.Lookup("ToxChannel"))
	if conv != nil {
	}
	peerCount := this._tox.GroupNumberPeers(groupNumber)
	if peerCount != 1 {
	}
	peerNumber := 0 // i create the groupchat and the 0th peer
	peerName, err := this._tox.GroupPeerName(groupNumber, peerNumber)
	if err != nil {
		log.Println(err)
	}
	chat := conv.GetChatData()
	chat.AddUser(peerName)
	selfName, err := this._tox.SelfGetName()
	if selfName != peerName {
		log.Println(selfName, peerName, selfName == peerName)
		log.Panicln("wtf")
	}
}
func (this *ToxPlugin) RejectChat(gc *purple.Connection, comp *purple.GHashTable) {
	log.Println("herhere")
	log.Println(comp.ToMap())
}
func (this *ToxPlugin) GetChatName(comp *purple.GHashTable) string {
	log.Println("herhere")
	log.Println(comp.ToMap())
	return ""
}
func (this *ToxPlugin) ChatInvite(gc *purple.Connection, id int, message string, who string) {
	log.Println("herhere")
	log.Println("herhere", id, message, who)
	friendNumber, err := this._tox.FriendByPublicKey(who)
	if err != nil {
		log.Println(err)
	}
	rc, err := this._tox.InviteFriend(friendNumber, id)
	if err != nil {
		log.Println(rc, err)
	}
}
func (this *ToxPlugin) ChatLeave(gc *purple.Connection, id int) {
	log.Println("herhere")
}
func (this *ToxPlugin) ChatWhisper(gc *purple.Connection, id int, who string, message string) {
	log.Println("herhere")
}
func (this *ToxPlugin) ChatSend(gc *purple.Connection, id int, message string, flags int) int {
	log.Println("herhere")
	n, err := this._tox.GroupMessageSend(id, message)
	if err != nil {
		log.Println(err)
	}
	if n == -1 {
		// log.Println("still send ok, wtf")
	}
	log.Println(n, id, message, flags)
	return len(message)
}

func (this *ToxPlugin) RoomlistGetList(gc *purple.Connection) {
	log.Println("herere")
}
