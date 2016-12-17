/*
  all implemention code except must code.
*/
package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"strings"

	"go-purple/purple"
	"go-purple/wechat-prpl/wechat"
)

func (this *WechatPlugin) eventHandler(evt *wechat.Event, ud interface{}) {
	ac := ud.(*purple.Account)
	gc := ac.GetConnection()
	if false {
		log.Println(ac, gc)
	}
	log.Println(ac.GetUserName(), int(evt.Type), evt.Type.String(), len(evt.Args))
	switch evt.Type {
	case wechat.EVT_GOT_QRCODE:
		iconData := []byte(evt.Args[0])
		purple.RequestAcceptCancelWithIconDemo(nil, gc, iconData, nil, nil)
	case wechat.EVT_SCANED_DATA:
		log.Println(int(evt.Type), evt.Type.String(), len(evt.Args), evt.Args)
	case wechat.EVT_REDIR_URL:
		log.Println(int(evt.Type), evt.Type.String(), len(evt.Args), evt.Args)
	case wechat.EVT_GOT_BASEINFO:
		log.Println(int(evt.Type), evt.Type.String(), len(evt.Args), len(evt.Args[0]))
		this.loadInitContact(ac, evt.Args[0])
		this.setupSelfInfo(ac)
		this.loadArticles(ac)
	case wechat.EVT_LOGIN_STATUS:
		switch evt.Args[0] {
		case "true":
			gc.ConnSetState(purple.CONNECTED)
		case "false":
			gc.ConnSetState(purple.DISCONNECTED)
		}
	case wechat.EVT_GOT_CONTACT:
		log.Println(int(evt.Type), evt.Type.String(), len(evt.Args), len(evt.Args[0]))
		this.loadAllContact(ac, evt.Args[0])

	case wechat.EVT_GOT_MESSAGE:
		log.Println("you have 1 new message", evt.Args[0][0:65])
		this.onWXMessage(gc, evt)

	case wechat.EVT_SAVEDATA:
		this.save_account(gc)
	}
}

func (this *WechatPlugin) loadInitContact(ac *purple.Account, initData string) {
	users := wechat.ParseWXInitData(initData)
	buddies := ac.FindBuddies("")

	group := purple.NewGroup("ACTIVES")
	for _, user := range users {
		pubkey := user.UserName
		name := user.NickName

		buddy := ac.FindBuddy(pubkey)
		if buddy == nil {
			found := false
			for _, _buddy := range buddies {
				if strings.HasPrefix(_buddy.GetName(), pubkey) {
					found = true
					buddy = _buddy
					break
				}
			}
			if !found {
				buddy = purple.NewBuddy(ac, pubkey, name)
				ac.AddBuddy(buddy)
				buddy.BlistAdd(group)
			}
		} else {
			if buddy.GetAliasOnly() != name {
				buddy.SetAlias(name)
			}
		}
		purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
	}
}

func (this *WechatPlugin) setupSelfInfo(ac *purple.Account) {
	me := this._wechat.Me()
	ac.SetUserName(me.UserName)
	ac.SetAlias(me.NickName)

	/*
		gc := ac.GetConnection()
		data, _ := ioutil.ReadFile("/home/gzleo/oss/src/go-purple/wechat-prpl/wechat/qrcode.jpg")
		purple.RequestAcceptCancelWithIconDemo(nil, gc, data, nil, nil)
	*/

	pubkey := reader_user
	name := reader_nick

	buddy := ac.FindBuddy(pubkey)
	if buddy == nil {
		buddy = purple.NewBuddy(ac, pubkey, name)
		ac.AddBuddy(buddy)
		buddy.BlistAdd(nil)
	} else {
		if buddy.GetAliasOnly() != name {
			buddy.SetAlias(name)
		}
	}
	purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
	log.Println("adding...", name, pubkey, purple.MyTid2())
}

var reader_user = "article_reader_longlonglong"
var reader_nick = "阅读&订阅"

func (this *WechatPlugin) loadArticles(ac *purple.Account) {
	pubkey := reader_user

	buddy := ac.FindBuddy(pubkey)
	if buddy == nil {
		log.Println("wtf")
	} else {
		gc := ac.GetConnection()
		mpas := this._wechat.Articles()
		for _, a := range mpas {
			msg := fmt.Sprintf("%s: %s", a.Title, a.Url)
			gc.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
		}
	}
}

func (this *WechatPlugin) loadAllContact(ac *purple.Account, contactData string) {
	users := wechat.ParseContactData(contactData)
	buddies := ac.FindBuddies("")

	group := purple.NewGroup("CONTACTS")
	for _, user := range users {
		pubkey := user.UserName
		name := user.NickName

		buddy := ac.FindBuddy(pubkey)
		if buddy == nil {
			found := false
			for _, _buddy := range buddies {
				if strings.HasPrefix(_buddy.GetName(), pubkey) {
					found = true
					buddy = _buddy
					break
				}
			}
			if !found {
				buddy = purple.NewBuddy(ac, pubkey, name)
				ac.AddBuddy(buddy)
				buddy.BlistAdd(group)
			}
		} else {
			if buddy.GetAliasOnly() != name {
				buddy.SetAlias(name)
			}
		}
		purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
	}
}

// 因为存储的name可能是friendId，也可能是pubkey。
func (this *WechatPlugin) findBuddyEx(ac *purple.Account, pubkeyOrFriendID string) *purple.Buddy {
	name := pubkeyOrFriendID
	buddy := ac.FindBuddy(name)
	if buddy == nil {
		buddies := ac.FindBuddies("")
		for _, buddy_ := range buddies {
			if strings.HasPrefix(buddy_.GetName(), name) {
				buddy = buddy_
				break
			}
		}
	}
	return buddy
}

// poll handlers

func (this *WechatPlugin) onWXMessage(gc *purple.Connection, evt *wechat.Event) {
	ac := gc.ConnGetAccount()

	log.Println(evt)
	msg := evt.Args[0]
	msgo := wechat.ParseMessage(msg)
	pubkey := msgo.ToUserName
	buddy := this.findBuddyEx(ac, pubkey)
	if buddy == nil {
		log.Println("wtf", pubkey, msgo.MsgId)
		buddy = purple.NewBuddy(ac, msgo.ToUserName, msgo.ToUserName)
		// gc.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
		gc.ServGotIM(buddy.GetName(), msgo.Content, purple.MESSAGE_RECV)
	} else {
		log.Println(buddy, msgo)
		// gc.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
		gc.ServGotIM(buddy.GetName(), msgo.Content, purple.MESSAGE_RECV)
	}
}

// optional callbacks
func (this *WechatPlugin) ChatInfo(gc *purple.Connection) []*purple.ProtoChatEntry {
	// log.Println(gc)

	infos := []*purple.ProtoChatEntry{
		purple.NewProtoChatEntry("WechatChannel", "_WechatChannel", true),
		purple.NewProtoChatEntry("GroupNumber", "_GroupNumber", false),
	}
	return infos
}

func (this *WechatPlugin) ChatInfoDefaults(gc *purple.Connection, chatName string) map[string]string {
	log.Println(gc)
	return nil
}

func (this *WechatPlugin) SendIM(gc *purple.Connection, who string, msg string) int {
	log.Println(gc, who, msg)
	me := this._wechat.Me()
	log.Println(me)

	if who == reader_user {
		return -1
	}

	ret := this._wechat.SendMessage(me.UserName, who, msg)
	log.Println(ret)
	if !ret {
		return -1
	}
	/*
		friendNumber, _ := this._wechat.FriendByPublicKey(who)
		len, err := this._wechat.FriendSendMessage(friendNumber, msg)
		if err != nil {
			log.Println(err, len)
			return -1
		}
		return int(len)
	*/
	return len(msg)
}

func (this *WechatPlugin) JoinChat(gc *purple.Connection, comp *purple.GHashTable) {
	log.Println(gc, comp.Lookup("WechatChannel"), comp.Lookup("GroupNumber"))
	// manual join from ui
	/*
		groupNumber, err := this._wechat.AddGroupChat()
		if err != nil {
			log.Println(err)
		}
		title := comp.Lookup("WechatChannel")
		this._wechat.GroupSetTitle(groupNumber, title)
		comp.Insert("GroupNumber", fmt.Sprintf("%d", groupNumber))
		conv := gc.ServGotJoinedChat(groupNumber, comp.Lookup("WechatChannel"))
		if conv != nil {
			conv.SetLogging(true)
		}
		this.UpdateMembers(groupNumber, conv)
	*/
}
func (this *WechatPlugin) JoinChatQuite(gc *purple.Connection, title string, groupNumber uint32) {
	/*
		this._wechat.GroupSetTitle(int(groupNumber), title)
		conv := gc.ServGotJoinedChat(int(groupNumber), title)
		if conv != nil {
		}
		this.UpdateMembers(int(groupNumber), conv)
	*/
}

func (this *WechatPlugin) RejectChat(gc *purple.Connection, comp *purple.GHashTable) {
	log.Println("herhere")
	log.Println(comp.ToMap())
}
func (this *WechatPlugin) GetChatName(comp *purple.GHashTable) string {
	log.Println("herhere")
	log.Println(comp.ToMap())
	return ""
}
func (this *WechatPlugin) ChatInvite(gc *purple.Connection, id int, message string, who string) {
	log.Println("herhere")
	log.Println("herhere", id, message, who)
	/*
		friendNumber, err := this._wechat.FriendByPublicKey(who)
		if err != nil {
			log.Println(err)
		}
		rc, err := this._wechat.InviteFriend(friendNumber, id)
		if err != nil {
			log.Println(rc, err)
		}
	*/
}
func (this *WechatPlugin) ChatLeave(gc *purple.Connection, id int) {
	log.Println("herhere")
}
func (this *WechatPlugin) ChatWhisper(gc *purple.Connection, id int, who string, message string) {
	log.Println("herhere")
}
func (this *WechatPlugin) ChatSend(gc *purple.Connection, id int, message string, flags int) int {
	log.Println("herhere")
	/*
		n, err := this._wechat.GroupMessageSend(id, message)
		if err != nil {
			log.Println(err)
		}
		if n == -1 {
			// log.Println("still send ok, wtf")
		}
		log.Println(n, id, message, flags)
	*/
	return len(message)
}

func (this *WechatPlugin) RoomlistGetList(gc *purple.Connection) {
	log.Println("herere")
}

func (this *WechatPlugin) AddBuddyWithInvite(gc *purple.Connection,
	buddy *purple.Buddy, group *purple.Group, message string) {
	log.Println(buddy, group, message)
	/*
		friendId := buddy.GetName()
		if len(message) == 0 {
			message = fmt.Sprintf("This is %s", this._wechat.SelfGetName())
		}
		friendNumber, err := this._wechat.FriendAdd(friendId, message)
		if err != nil {
			log.Println(err, friendNumber)
		} else {
			// gc.ConnGetAccount().AddBuddy(buddy)
			// buddy.BlistAdd(nil)
			buddy := gc.ConnGetAccount().FindBuddy(friendId)
			log.Println(buddy)
		}
	*/
}

func (this *WechatPlugin) RemoveBuddy(gc *purple.Connection, buddy *purple.Buddy, group *purple.Group) {
	/*
		friendId := buddy.GetName()
		friendNumber, err := this._wechat.FriendByPublicKey(friendId)
		if err != nil {
			log.Println(err, friendNumber)
		} else {
			_, err = this._wechat.FriendDelete(friendNumber)
			if err != nil {
				log.Println(err)
			} else {
				this.save_account(gc)
			}
		}
	*/
}

func (this *WechatPlugin) GetInfo(gc *purple.Connection, who string) {
	var friendName string
	var friendStmsg string
	u := this._wechat.GetUser(who)
	if u != nil {
		friendName = u.NickName
		friendStmsg = u.Signature
	} else {
		log.Println("user not found:", who)
		buddy := gc.ConnGetAccount().FindBuddy(who)
		if buddy != nil {
			friendName = buddy.GetAliasOnly()
		} else {
			log.Println("buddy not found:", who)
		}
	}

	uinfo := purple.NewNotifyUserInfo()
	uinfo.AddPair("nickname", friendName)
	uinfo.AddPair("status message", friendStmsg)
	// uinfo.AddPair("seen", fmt.Sprintf("%d", seen))
	uinfo.AddPair("hehehe", "efffff")
	uinfo.AddPair("hehehe12", "efffff456")
	uinfo.AddPair("Uin", "0")

	gc.NotifyUserInfo(who, uinfo, func(ud interface{}) {
		log.Println("closed", ud)
	}, 123)

}

func (this *WechatPlugin) StatusText(buddy *purple.Buddy) string {
	/*
		who := buddy.GetName()
		friendNumber, err := this._wechat.FriendByPublicKey(who)
		if err != nil {
			log.Println(err, friendNumber, who)
		}
		friendStmsg, err := this._wechat.FriendGetStatusMessage(friendNumber)
		return friendStmsg
	*/
	return ""
}

// utils
func (this *WechatPlugin) UpdateMembers(groupNumber int, conv *purple.Conversation) {
	/*
		chat := conv.GetChatData()
		// TODO member list diff and clean, so it is member list sync
		t := this._wechat
		plst := chat.GetUsers()
		tlst := t.GroupGetNames(groupNumber)
		klst := t.GroupGetPeerPubkeys(groupNumber)
		mlst := t.GroupGetPeers(groupNumber)
		peerCount := t.GroupNumberPeers(groupNumber)
		if len(tlst) != peerCount {
			log.Println("wtf")
		}

		if true {
			log.Println("need sync names...")
			log.Println("purple list:", plst)
			log.Println("wechat list:", tlst)
			log.Println("pubkey list:", klst)
			log.Println("peer list:", mlst)

			// remove not existed
			for _, pname := range plst {
				found := false
				for _, tname := range tlst {
					if tname == pname {
						found = true
					}
				}
				if found == false {
					chat.RemoveUser(pname) // should already destroy the ConvChatBuddy here
					cbbuddy := chat.FindBuddy(pname)
					cbbuddy.Destroy()
				}
			}

			// add new
			for peerNumber, pubkey := range mlst {
				found := false
				peerName, err := t.GroupPeerName(groupNumber, peerNumber)
				if err != nil {
				}
				for _, pname := range plst {
					if pname == peerName {
						found = true
					}
				}
				if found == false {
					isours := t.GroupPeerNumberIsOurs(groupNumber, peerNumber)
					if isours == true {
					}
					if true {
						chat.AddUser(peerName)
						cbbudy := chat.FindBuddy(peerName)
						cbbudy.SetAlias(pubkey)
					}
				}
			}
		}
	*/
}
