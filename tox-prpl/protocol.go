/*
  all implemention code except must code.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go-purple/purple"

	"github.com/kitech/go-toxcore"
)

func (this *ToxPlugin) setupSelfInfo(ac *purple.Account) {
	if len(ac.GetAlias()) > 0 {
		this._tox.SelfSetName(ac.GetAlias())
	} else {
		this._tox.SelfSetName(ac.GetUserName())
	}
	this._tox.SelfSetStatusMessage("It's from gotox-prpl, hoho.")

	iconName := ac.GetString("buddy_icon")
	iconFile := purple.UserDir() + "/icons/" + iconName
	_, err := os.Stat(iconFile)
	if err == os.ErrNotExist {
		log.Println(purple.UserDir(), iconName, iconFile)
	}

	if false {
		conn := ac.GetConnection()
		purple.RequestYesNoDemo(this, conn, func(ud interface{}) {
			log.Println("request yes")
		}, func(ud interface{}) {
			log.Println("request no")
		})
	}
}

func (this *ToxPlugin) setupCallbacks(ac *purple.Account) {
	conn := ac.GetConnection()
	gc := conn

	this._tox.CallbackSelfConnectionStatus(func(t *tox.Tox, status int, d interface{}) {
		if status > tox.CONNECTION_NONE {
			conn.ConnSetState(purple.CONNECTED) // 设置为已连接状态，则好友会显示。
			// a helper for help me
			tryAddFixedFriends(this._tox, conn)
		} else {
			conn.ConnSetState(purple.DISCONNECTED)
		}
		this.save_account(gc)
	}, ac)

	this._tox.CallbackFriendRequest(this.onFriendRequest, ac)

	this._tox.CallbackFriendConnectionStatus(this.onFriendConnectionStatus, ac)
	this._tox.CallbackFriendStatus(this.onFriendStatus, ac)

	this._tox.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
		log.Println(friendNumber, msg, purple.MyTid2())
		conn := ac.GetConnection()
		pubkey, err := t.FriendGetPublicKey(friendNumber)
		if err != nil {
			log.Println(err)
		} else {
			buddy := this.findBuddyEx(ac, pubkey)
			if buddy == nil {
				log.Println("wtf", friendNumber, pubkey)
			} else {
				conn.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
			}
		}

		// some extra
		fixSpecialMessage(t, friendNumber, msg)
	}, ac)

	this._tox.CallbackFriendName(func(t *tox.Tox, friendNumber uint32, name string, d interface{}) {
		// log.Println(friendNumber, name)
		pubkey, err := t.FriendGetPublicKey(friendNumber)
		if err != nil {
			log.Println(err, pubkey)
		} else {
			buddy := this.findBuddyEx(ac, pubkey)
			if buddy != nil {
				if buddy.GetAliasOnly() != name {
					buddy.SetAlias(name)
				}
			} else {
				log.Println("wtf", friendNumber, name, pubkey)
				// buddy := purple.NewBuddy(ac, pubkey, name)
				// ac.AddBuddy(buddy)
				// buddy.BlistAdd(nil)
			}
		}
	}, nil)

	this._tox.CallbackGroupNameListChange(func(t *tox.Tox, groupNumber int,
		peerNumber int, change uint8, d interface{}) {
		log.Println(groupNumber, peerNumber, change)
		conn := ac.GetConnection()
		conv := conn.ConnFindChat(groupNumber)
		switch change {
		case tox.CHAT_CHANGE_PEER_ADD:
		case tox.CHAT_CHANGE_PEER_DEL:
		case tox.CHAT_CHANGE_PEER_NAME:
		}
		this.UpdateMembers(groupNumber, conv)
	}, ac)

	this._tox.CallbackGroupMessage(this.onGroupMessage, ac)
	this._tox.CallbackGroupAction(this.onGroupAction, ac)

	// TODO should notify UI first and think groupbot's invite also
	this._tox.CallbackGroupInvite(this.onGroupInvite, ac)
	this._tox.CallbackGroupTitle(this.onGroupTitle, ac)

	this._tox.CallbackFriendTyping(this.onFriendTyping, ac)
}

func (this *ToxPlugin) loadFriends(ac *purple.Account) {
	fns := this._tox.SelfGetFriendList()
	if fns == nil || len(fns) == 0 {
		log.Println("oh, you have 0 friends")
	}
	buddies := ac.FindBuddies("")
	for _, fn := range fns {
		name, err := this._tox.FriendGetName(fn)
		pubkey, err := this._tox.FriendGetPublicKey(fn)
		if err != nil {
			log.Println(err)
		}
		if len(name) == 0 {
			name = "GoTox User"
		}
		buddy := ac.FindBuddy(pubkey)
		if buddy == nil {
			found := false
			for _, _buddy := range buddies {
				if strings.HasPrefix(_buddy.GetName(), pubkey) {
					found = true
					break
				}
			}
			if !found {
				buddy = purple.NewBuddy(ac, pubkey, name)
				ac.AddBuddy(buddy)
				buddy.BlistAdd(nil)
			}
		} else {
			if buddy.GetAliasOnly() != name {
				buddy.SetAlias(name)
			}
		}
		// purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
	}
	// set inital offline
	buddies = ac.FindBuddies("")
	for _, buddy := range buddies {
		purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_OFFLINE_STR)
	}
}

// 因为存储的name可能是friendId，也可能是pubkey。
func (this *ToxPlugin) findBuddyEx(ac *purple.Account, pubkeyOrFriendID string) *purple.Buddy {
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

// tox callbacks
type minictx struct {
	gc *purple.Connection
	s0 string
}

func (this *ToxPlugin) requestYescbForFriendRequest(d interface{}) {
	ctx := d.(*minictx)
	conn := ctx.gc
	pubkey := ctx.s0
	ac := conn.ConnGetAccount()
	gc := ctx.gc

	friendNumber, err := this._tox.FriendAddNorequest(pubkey)
	if err != nil {
		log.Println(err, friendNumber)
	}
	this.save_account(conn)
	friendName, err := this._tox.FriendGetName(friendNumber)
	if len(friendName) == 0 {
		friendName = "GoTox User"
	}
	buddy := this.findBuddyEx(ac, pubkey)
	// log.Println(buddy, buddy.GetName(), buddy.GetAliasOnly())
	if buddy == nil {
		buddy := purple.NewBuddy(ac, pubkey, friendName)
		ac.AddBuddy(buddy)
		buddy.BlistAdd(nil)
	}
	this.save_account(gc)
}

func (this *ToxPlugin) onFriendRequest(t *tox.Tox, pubkey, msg string, d interface{}) {
	log.Println("hehhe", pubkey, msg)
	ac := d.(*purple.Account)
	conn := ac.GetConnection()

	// TODO notice UI and then make desision
	ctx := &minictx{conn, pubkey}
	log.Println(purple.GoID())
	purple.RequestAcceptCancel(ctx, conn, "New Friend Request", msg, pubkey,
		this.requestYescbForFriendRequest, func(ud interface{}) {
			// reject?
			log.Println(ud)
		})
}

func (this *ToxPlugin) onFriendConnectionStatus(t *tox.Tox, friendNumber uint32, status int, d interface{}) {
	log.Println(friendNumber, status)
	ac := d.(*purple.Account)
	gc := ac.GetConnection()

	pubkey, _ := t.FriendGetPublicKey(friendNumber)
	name, _ := t.FriendGetName(friendNumber)

	buddy := this.findBuddyEx(ac, pubkey)
	if buddy == nil {
		log.Println("can not find buddy:", name, pubkey)
	} else {
		switch status {
		case tox.CONNECTION_NONE:
			purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_OFFLINE_STR)
		case tox.CONNECTION_TCP:
			fallthrough
		case tox.CONNECTION_UDP:
			fallthrough
		default:
			purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		}
	}

	//
	tryJoinFixedGroups(t, gc, friendNumber, status)
	this.save_account(gc)
}

func (this *ToxPlugin) onFriendStatus(t *tox.Tox, friendNumber uint32, status int, d interface{}) {
	ac := d.(*purple.Account)

	pubkey, _ := t.FriendGetPublicKey(friendNumber)
	name, _ := t.FriendGetName(friendNumber)

	buddy := this.findBuddyEx(ac, pubkey)
	if buddy == nil {
		log.Println("can not find buddy:", name, pubkey)
	} else {
		switch status {
		case tox.USER_STATUS_AWAY:
			purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_AWAY_STR)
		case tox.USER_STATUS_BUSY:
			purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_BUSY_STR)
		case tox.USER_STATUS_NONE:
			// purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
			fallthrough
		default:
			log.Println("what can i do?", status)
		}
	}
}

func (this *ToxPlugin) onGroupInvite(t *tox.Tox,
	friendNumber uint32, itype uint8, data []byte, d interface{}) {
	ac := d.(*purple.Account)
	conn := ac.GetConnection()

	log.Println(friendNumber, len(data), itype)
	pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
	if err != nil {
		log.Println(err)
	}
	acceptInvite := func(interface{}) {
		var groupNumber int
		var err error
		switch itype {
		case tox.GROUPCHAT_TYPE_AV:
			groupNumber, err = this._tox.JoinAVGroupChat(friendNumber, data)
			if err != nil {
				log.Println(err, groupNumber)
			}
		case tox.GROUPCHAT_TYPE_TEXT:
			groupNumber, err = this._tox.JoinGroupChat(friendNumber, data)
			if err != nil {
				log.Println(err, groupNumber)
			}
		default:
			log.Panicln("wtf")
		}
		if err == nil {
			groupTitle, err := this._tox.GroupGetTitle(groupNumber)
			if err != nil {
				log.Println(err, groupTitle)
				groupTitle = DEFAULT_GROUPCHAT_TITLE
			}
			conv := conn.ServGotJoinedChat(groupNumber, groupTitle)
			if conv == nil {
				log.Println("join chat failed:", conv, groupNumber, groupTitle)
			} else if conv != nil {
				conv.SetData("GroupNumber", fmt.Sprintf("%d", groupNumber))
				conv.SetLogging(true)
			}
		}
	}
	if strings.HasPrefix(groupbot, pubkey) {
		// go on without notify UI
		acceptInvite(nil)
	} else {
		friendName, err := this._tox.FriendGetName(friendNumber)
		if err != nil {
			log.Println("wtf")
		}
		purple.RequestAcceptCancel(nil, conn, "New Group Invite",
			fmt.Sprintf("Are you want to join %s's group?", friendName),
			friendName, acceptInvite, nil)
	}
}

func (this *ToxPlugin) onGroupTitle(t *tox.Tox,
	groupNumber, peerNumber int, title string, ud interface{}) {
	ac := ud.(*purple.Account)
	conn := ac.GetConnection()

	log.Println(groupNumber, peerNumber, title)
	conv := conn.ConnFindChat(groupNumber)
	if conv == nil {
		log.Println("can not found conv, create new chat:", groupNumber, title)
		conv = conn.ServGotJoinedChat(groupNumber, title)
	} else {
		if conv.GetName() != title {
			conv.SetName(title)
		}
	}
	if false {
		// 即使参数相同，也不可忽略的调用
		conv2 := conn.ServGotJoinedChat(groupNumber, title)
		if conv != conv2 {
			log.Println("wtf, maybe remove one")
		}
	}
}

func (this *ToxPlugin) onGroupMessage(t *tox.Tox, groupNumber int,
	peerNumber int, message string, ud interface{}) {
	log.Println(groupNumber, peerNumber, message)
	ac := ud.(*purple.Account)
	conn := ac.GetConnection()
	pubkey, err := t.GroupPeerPubkey(groupNumber, peerNumber)
	peerName, err := t.GroupPeerName(groupNumber, peerNumber)
	if err != nil {
		log.Println(err, peerName, pubkey)
	} else {
		if false {
			conn.ServGotChatIn(groupNumber, pubkey, purple.MESSAGE_RECV, message)
		}
		groupTitle, err := this._tox.GroupGetTitle(groupNumber)
		if err != nil {
			log.Println(err, groupTitle)
		}

		convs := purple.GetConversations()
		for _, conv := range convs {
			log.Println(conv.GetName(), conv.GetChatData().GetUsers(), conv.GetData("GroupNumber"))
		}
		if len(convs) == 0 {
			log.Println("can not find conv chat:", groupNumber, groupTitle)
		}
		// conversation should be created when invited
		conv := conn.FindChat(int(groupNumber))
		if conv == nil {
			log.Println("can not find conv chat:", groupNumber, groupTitle)
		} else {
			// conv.GetChatData().Send(message) // infinite message loop!!!
			conn.ServGotChatIn(groupNumber, peerName, purple.MESSAGE_RECV, message)
			// conn.ServGotChatIn(groupNumber, pubkey, purple.MESSAGE_RECV, message)
		}
	}
}

func (this *ToxPlugin) onGroupAction(t *tox.Tox, groupNumber int, peerNumber int, action string, ud interface{}) {
	message := PREFIX_ACTION + action
	this.onGroupMessage(t, groupNumber, peerNumber, message, ud)
}

func (this *ToxPlugin) onFriendTyping(t *tox.Tox, friendNumber uint32, isTyping uint8, ud interface{}) {
	ac := ud.(*purple.Account)
	gc := ac.GetConnection()

	pubkey, err := this._tox.FriendGetPublicKey(friendNumber)
	if err != nil {
		log.Println(err, pubkey)
	} else {
		buddy := this.findBuddyEx(ac, pubkey)
		if buddy == nil {
			log.Println("wtf", friendNumber, pubkey)
		} else {
			timeout := 3
			if isTyping == 1 {
				gc.ServGotTyping(buddy.GetName(), timeout, purple.TYPING)
			} else {
				gc.ServGotTypingStopped(buddy.GetName())
			}
		}
	}
}

// purple optional callbacks
func (this *ToxPlugin) ChatInfo(gc *purple.Connection) []*purple.ProtoChatEntry {
	// log.Println(gc)

	infos := []*purple.ProtoChatEntry{
		purple.NewProtoChatEntry("ToxChannel", "_ToxChannel", true),
		purple.NewProtoChatEntry("GroupNumber", "_GroupNumber", false),
	}
	return infos
}

func (this *ToxPlugin) ChatInfoDefaults(gc *purple.Connection, chatName string) map[string]string {
	log.Println(gc)
	return nil
}

func (this *ToxPlugin) SendIM(gc *purple.Connection, who string, msg string) int {
	log.Println(gc, who, msg)
	friendNumber, _ := this._tox.FriendByPublicKey(who)
	var ln uint32
	var err error
	if strings.HasPrefix(msg, PREFIX_ACTION) {
		ln, err = this._tox.FriendSendAction(friendNumber, msg[len(PREFIX_ACTION):])
	} else {
		ln, err = this._tox.FriendSendMessage(friendNumber, msg)
	}
	if err != nil {
		log.Println(err, ln)
		return -1
	}
	return int(ln)
}

func (this *ToxPlugin) JoinChat(gc *purple.Connection, comp *purple.GHashTable) {
	log.Println(gc, comp.Lookup("_ToxChannel"), comp.Lookup("GroupNumber"))
	if joinChatSpecialFixed(this._tox, comp) {
		return
	}

	// manual join from ui
	groupNumber, err := this._tox.AddGroupChat()
	if err != nil {
		log.Println(err)
	}
	title := comp.Lookup("_ToxChannel")
	this._tox.GroupSetTitle(groupNumber, title)
	comp.Insert("GroupNumber", fmt.Sprintf("%d", groupNumber))
	conv := gc.ServGotJoinedChat(groupNumber, comp.Lookup("_ToxChannel"))
	if conv != nil {
		conv.SetLogging(true)
	}
	this.UpdateMembers(groupNumber, conv)
}
func (this *ToxPlugin) JoinChatQuite(gc *purple.Connection, title string, groupNumber uint32) {
	this._tox.GroupSetTitle(int(groupNumber), title)
	conv := gc.ServGotJoinedChat(int(groupNumber), title)
	if conv != nil {
	}
	this.UpdateMembers(int(groupNumber), conv)
}

// TODO what?
func (this *ToxPlugin) RejectChat(gc *purple.Connection, comp *purple.GHashTable) {
	log.Println("herhere")
	log.Println(comp.ToMap())
}

// TODO what?
func (this *ToxPlugin) GetChatName(comp *purple.GHashTable) string {
	log.Println(comp.ToMap())
	groupName := comp.Lookup("_ToxChannel")
	if isFixedGroup(groupName) {
		return groupName
		// return DEFAULT_GROUPCHAT_TITLE
	}
	return groupName
}
func (this *ToxPlugin) ChatInvite(gc *purple.Connection, id int, message string, who string) {
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
	groupNumber := id
	title, err := this._tox.GroupGetTitle(groupNumber)
	if err != nil {
		log.Println(err)
	}
	// TODO 检查设置，是否是关闭不离开的会话
	_, err = this._tox.DelGroupChat(groupNumber)
	if err != nil {
		log.Println(err)
	}
	log.Println("proto leaved group:", groupNumber, title)
}
func (this *ToxPlugin) ChatWhisper(gc *purple.Connection, id int, who string, message string) {
	log.Println("herhere")
}
func (this *ToxPlugin) ChatSend(gc *purple.Connection, id int, message string, flags int) int {
	var n int
	var err error
	if strings.HasPrefix(message, PREFIX_ACTION) {
		n, err = this._tox.GroupActionSend(id, message[len(PREFIX_ACTION):])
	} else {
		n, err = this._tox.GroupMessageSend(id, message)
	}
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

func (this *ToxPlugin) AddBuddyWithInvite(gc *purple.Connection,
	buddy *purple.Buddy, group *purple.Group, message string) {
	log.Println(buddy, group, message)
	friendId := buddy.GetName()
	if len(message) == 0 {
		message = fmt.Sprintf("This is %s", this._tox.SelfGetName())
	}
	friendNumber, err := this._tox.FriendAdd(friendId, message)
	if err != nil {
		log.Println(err, friendNumber)
	} else {
		// gc.ConnGetAccount().AddBuddy(buddy)
		// buddy.BlistAdd(nil)
		buddy := gc.ConnGetAccount().FindBuddy(friendId)
		log.Println(buddy)
	}
}

func (this *ToxPlugin) RemoveBuddy(gc *purple.Connection, buddy *purple.Buddy, group *purple.Group) {
	friendId := buddy.GetName()
	friendNumber, err := this._tox.FriendByPublicKey(friendId)
	if err != nil {
		log.Println(err, friendNumber)
	} else {
		_, err = this._tox.FriendDelete(friendNumber)
		if err != nil {
			log.Println(err)
		} else {
			this.save_account(gc)
		}
	}
}

func (this *ToxPlugin) GetInfo(gc *purple.Connection, who string) {
	var uinfo *purple.NotifyUserInfo
	var err error

	friendNumber, err := this._tox.FriendByPublicKey(who)
	if err != nil {
		log.Println(err, friendNumber, who)
		// maybe chat peer's name, from group member list
		pubkeys := make(map[int]string, 0)
		groupNumbers := this._tox.GetChatList()
		for _, groupNumber := range groupNumbers {
			peerMap := this._tox.GroupGetPeers(int(groupNumber))
			for peerNumber, pubkey := range peerMap {
				peerName, err := this._tox.GroupPeerName(int(groupNumber), peerNumber)
				if err != nil {
					log.Println(err)
				}
				if peerName == who {
					pubkeys[int(groupNumber)] = pubkey
				}
			}
		}

		log.Println("find matched pubkeys:", len(pubkeys))
		if len(pubkeys) == 0 {
			log.Println("wtf")
		} else {
			uinfo = purple.NewNotifyUserInfo()
			for groupNumber, pubkey := range pubkeys {
				uinfo.AddPair("nickname", who)
				uinfo.AddPair("id", pubkeys[0])
				conv := gc.ConnFindChat(groupNumber)
				if conv == nil {
					log.Println("can not found coconv:", groupNumber, pubkey)
				} else {
					chat := conv.GetChatData()
					ccbuddy := chat.FindBuddy(who)
					if ccbuddy == nil {
						log.Println("can not found cobuddy:", groupNumber, pubkey, conv.GetName())
					} else {
						joinTime := ccbuddy.GetAttribute("joinTime")
						uinfo.AddPair("joinTime", joinTime)
						attrPubkey := ccbuddy.GetAttribute("pubkey")
						if attrPubkey != pubkey {
							log.Println("not match pubkey:", attrPubkey, pubkey, who, conv.GetName())
						}
					}
				}
			}
		}
	} else {
		friendName, err := this._tox.FriendGetName(friendNumber)
		if err != nil {
		}
		friendStmsg, err := this._tox.FriendGetStatusMessage(friendNumber)
		friendSt, err := this._tox.FriendGetStatus(friendNumber)
		seen, err := this._tox.FriendGetLastOnline(friendNumber)

		uinfo = purple.NewNotifyUserInfo()
		uinfo.AddPair("nickname", friendName)
		uinfo.AddPair("id", who)
		uinfo.AddPair("status message", friendStmsg)
		switch friendSt {
		case tox.USER_STATUS_AWAY:
			uinfo.AddPair("status", STATUS_AWAY_STR)
		case tox.USER_STATUS_BUSY:
			uinfo.AddPair("status", STATUS_BUSY_STR)
		default:
			connst, err := this._tox.FriendGetConnectionStatus(friendNumber)
			if err != nil {
				log.Println(err)
			}
			switch connst {
			case tox.CONNECTION_NONE:
				uinfo.AddPair("status", STATUS_OFFLINE_STR)
			default:
				uinfo.AddPair("status", STATUS_ONLINE_STR)
			}
		}
		uinfo.AddPair("seen", fmt.Sprintf("%d", seen))
		uinfo.AddPair("hehehe", "efffff")
		uinfo.AddPair("hehehe12", "efffff456")
	}

	gc.NotifyUserInfo(who, uinfo, func(ud interface{}) {
		log.Println("closed", ud)
	}, 123)
}

func (this *ToxPlugin) StatusText(buddy *purple.Buddy) string {
	who := buddy.GetName()
	if this._tox == nil {
		log.Println("already closed tox instance")
		// return ""
	}
	friendNumber, err := this._tox.FriendByPublicKey(who)
	if err != nil {
		log.Println(err, friendNumber, who)
	}
	friendStmsg, err := this._tox.FriendGetStatusMessage(friendNumber)
	return friendStmsg
}

func (this *ToxPlugin) SetChatTopic(gc *purple.Connection, id int, topic string) {
	n, err := this._tox.GroupSetTitle(id, topic)
	if err != nil {
		log.Println(err, n)
	} else {
		conv := gc.ConnFindChat(id)
		if conv != nil {
			if conv.GetName() != topic {
				conv.SetName(topic)
			}
		}
	}
}

func (this *ToxPlugin) Normalize(gc *purple.Connection, who string) string {
	return strings.ToUpper(who)
}

func (this *ToxPlugin) SendTyping(gc *purple.Connection, name string, state int) uint {
	switch state {
	case purple.NOT_TYPING:
		fn, err := this._tox.FriendByPublicKey(name)
		if err != nil {
			log.Println(err)
		} else {
			this._tox.SelfSetTyping(fn, false)
		}
	case purple.TYPING:
		fn, err := this._tox.FriendByPublicKey(name)
		if err != nil {
			log.Println(err)
		} else {
			this._tox.SelfSetTyping(fn, true)
		}
		return 2
	case purple.TYPED:
		fn, err := this._tox.FriendByPublicKey(name)
		if err != nil {
			log.Println(err)
		} else {
			this._tox.SelfSetTyping(fn, false)
		}
	}

	return 0
}

// utils
func (this *ToxPlugin) UpdateMembers(groupNumber int, conv *purple.Conversation) {
	chat := conv.GetChatData()
	// ac := conv.GetAccount()

	// TODO member list diff and clean, so it is member list sync
	t := this._tox
	memList := chat.GetUsers()
	nameList := t.GroupGetNames(groupNumber)
	pubkeyList := t.GroupGetPeerPubkeys(groupNumber)
	peerMap := t.GroupGetPeers(groupNumber)
	peerCount := t.GroupNumberPeers(groupNumber)
	if len(nameList) != peerCount {
		log.Println("wtf")
	}

	if false {
		log.Println("need sync names...")
		log.Println("purple list:", memList)
		log.Println("tox list:", nameList)
		log.Println("pubkey list:", pubkeyList)
		log.Println("peer list:", peerMap)
	}

	// still use nick for chatroom
	if true {
		// remove not existed
		for _, ccbuddy := range memList {
			found := false
			for _, tname := range nameList {
				if tname == ccbuddy.GetName() {
					found = true
				}
			}
			if found == false {
				// should already destroy the ConvChatBuddy here
				chat.RemoveUser(ccbuddy.GetName())
				// ccbuddy.Destroy() // the remove with call this too
			}
		}

		memList = chat.GetUsers() // reget
		// add new
		for peerNumber, pubkey := range peerMap {
			peerName, err := t.GroupPeerName(groupNumber, peerNumber)
			if err != nil {
				log.Println(err)
			}
			found := false
			for _, ccbuddy := range memList {
				if ccbuddy.GetName() == peerName {
					found = true
				}
			}
			if found == false {
				isours := t.GroupPeerNumberIsOurs(groupNumber, peerNumber)
				if isours == true {
				}
				if true {
					chat.AddUser(peerName)
					ccbuddy := chat.FindBuddy(peerName)
					ccbuddy.SetAttribute(chat, "pubkey", pubkey)
					ccbuddy.SetAttribute(chat, "joinTime", time.Now().String())
					// log.Println(peerName, ccbuddy.GetAlias(), ccbuddy.GetName(), ccbuddy.GetAttribute("pubkey"))
				}
			}
		}
	}

}
