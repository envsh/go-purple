/*
  all implemention code except must code.
*/
package main

import (
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
	log.Println("herhere")
	log.Println(comp.ToMap())
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
}
func (this *ToxPlugin) ChatLeave(gc *purple.Connection, id int) {
	log.Println("herhere")
}
func (this *ToxPlugin) ChatWhisper(gc *purple.Connection, id int, who string, message string) {
	log.Println("herhere")
}
func (this *ToxPlugin) ChatSend(gc *purple.Connection, id int, message string, flags int) int {
	log.Println("herhere")
	return 0
}
