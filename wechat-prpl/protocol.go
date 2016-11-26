/*
  all implemention code except must code.
*/
package main

import (
	// "fmt"
	// "io/ioutil"
	"log"
	"strings"

	"go-purple/purple"
	"go-purple/wechat-prpl/wechat"
)

func (this *WechatPlugin) setupSelfInfo(ac *purple.Account) {
	/*
		gc := ac.GetConnection()
		data, _ := ioutil.ReadFile("/home/gzleo/oss/src/go-purple/wechat-prpl/wechat/qrcode.jpg")
		purple.RequestAcceptCancelWithIconDemo(nil, gc, data, nil, nil)
	*/
	/*
		if len(ac.GetAlias()) > 0 {
			this._wechat.SelfSetName(ac.GetAlias())
		} else {
			this._wechat.SelfSetName(ac.GetUserName())
		}
		this._wechat.SelfSetStatusMessage("It's from gowechat-prpl, hoho.")

		if false {
			conn := ac.GetConnection()
			purple.RequestYesNoDemo(this, conn, func(ud interface{}) {
				log.Println("request yes")
			}, func(ud interface{}) {
				log.Println("request no")
			})
		}
	*/
}

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
		msgo := wechat.ParseMessage(evt.Args[0])
		pubkey := msgo.FromUserName
		msg := evt.Args[0]
		buddy := this.findBuddyEx(ac, pubkey)
		if buddy == nil {
			log.Println("wtf", pubkey, msgo.MsgId)
			buddy = purple.NewBuddy(ac, msgo.FromUserName, msgo.FromUserName)
			gc.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
		} else {
			gc.ServGotIM(buddy.GetName(), msg, purple.MESSAGE_RECV)
		}

	}
}

func (this *WechatPlugin) setupCallbacks(ac *purple.Account) {
	// conn := ac.GetConnection()

	/*
		this._wechat.CallbackSelfConnectionStatus(func(t *tox.Tox, status uint32, d interface{}) {
			if status > tox.CONNECTION_NONE {
				conn.ConnSetState(purple.CONNECTED) // 设置为已连接状态，则好友会显示。
				// a helper for help me
			} else {
				conn.ConnSetState(purple.DISCONNECTED)
			}
		}, ac)

		this._wechat.CallbackFriendRequest(func(t *tox.Tox, pubkey, msg string, d interface{}) {
			log.Println("hehhe", pubkey, msg)
			// TODO notice UI and then make desision
			purple.RequestAcceptCancel(this, conn, "New Friend Request", msg,
				func(ud interface{}) {
					friendNumber, err := this._wechat.FriendAddNorequest(pubkey)
					if err != nil {
						log.Println(err, friendNumber)
					}
					this.save_account(conn)
					friendName, err := this._wechat.FriendGetName(friendNumber)
					if len(friendName) == 0 {
						friendName = "GoWechat User"
					}
					buddy := this.findBuddyEx(ac, pubkey)
					// log.Println(buddy, buddy.GetName(), buddy.GetAliasOnly())
					if buddy == nil {
						buddy := purple.NewBuddy(ac, pubkey, friendName)
						ac.AddBuddy(buddy)
						buddy.BlistAdd(nil)
					}
				}, func(ud interface{}) {
					// reject?
				})
		}, ac)

		this._wechat.CallbackFriendConnectionStatus(func(t *tox.Tox, friendNumber uint32, status uint32, d interface{}) {
			log.Println(friendNumber, status)
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
					purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_BUSY_STR)
				case tox.CONNECTION_UDP:
					purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
				}
			}

		}, ac)

		this._wechat.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, d interface{}) {
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
		}, ac)

		this._wechat.CallbackFriendName(func(t *tox.Tox, friendNumber uint32, name string, d interface{}) {
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

		this._wechat.CallbackGroupNameListChange(func(t *tox.Tox, groupNumber int,
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

		this._wechat.CallbackGroupMessage(func(t *tox.Tox, groupNumber int,
			peerNumber int, message string, d interface{}) {
			log.Println(groupNumber, peerNumber, message)
			conn := ac.GetConnection()
			pubkey, err := t.GroupPeerPubkey(groupNumber, peerNumber)
			peerName, err := t.GroupPeerName(groupNumber, peerNumber)
			if err != nil {
				log.Println(err, peerName, pubkey)
			} else {
				if false {
					conn.ServGotChatIn(groupNumber, pubkey, purple.MESSAGE_RECV, message)
				}
				groupTitle, err := this._wechat.GroupGetTitle(groupNumber)
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
		}, ac)

		// TODO should notify UI first
		this._wechat.CallbackGroupInvite(func(t *tox.Tox,
			friendNumber uint32, itype uint8, data []byte, d interface{}) {
			log.Println(friendNumber, len(data), itype)
			var groupNumber int
			var err error
			switch itype {
			case tox.GROUPCHAT_TYPE_AV:
				groupNumber, err = this._wechat.JoinAVGroupChat(friendNumber, data)
				if err != nil {
					log.Println(err, groupNumber)
				}
			case tox.GROUPCHAT_TYPE_TEXT:
				groupNumber, err = this._wechat.JoinGroupChat(friendNumber, data)
				if err != nil {
					log.Println(err, groupNumber)
				}
			default:
				log.Panicln("wtf")
			}
			if err == nil {
				groupTitle, err := this._wechat.GroupGetTitle(groupNumber)
				if err != nil {
					log.Println(err, groupTitle)
					groupTitle = DEFAULT_GROUPCHAT_TITLE
				}
				conv := conn.ServGotJoinedChat(groupNumber, groupTitle)
				if conv == nil {
					log.Println("join chat failed:", conv, groupNumber, groupTitle)
				} else if conv != nil {
					conv.SetData("GroupNumber", fmt.Sprintf("%d", groupNumber))
				}
			}
		}, ac)

		this._wechat.CallbackGroupTitle(func(t *tox.Tox,
			groupNumber int, peerNumber int, title string, d interface{}) {
			log.Println(groupNumber, peerNumber, title)
			conv := conn.ConnFindChat(groupNumber)
			if conv != nil {
				if conv.GetName() != title {
					conv.SetName(title)
				}
			}
		}, ac)

	*/
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
		// purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
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
		// purple.PrplGotUserStatus(ac, buddy.GetName(), STATUS_ONLINE_STR)
		log.Println("adding...", name, pubkey, purple.MyTid2())
	}
}

func (this *WechatPlugin) loadFriends(ac *purple.Account) {
	/*
		fns := this._wechat.SelfGetFriendList()
		if fns == nil || len(fns) == 0 {
			log.Println("oh, you have 0 friends")
		}
		buddies := ac.FindBuddies("")
		for _, fn := range fns {
			name, err := this._wechat.FriendGetName(fn)
			pubkey, err := this._wechat.FriendGetPublicKey(fn)
			if err != nil {
				log.Println(err)
			}
			if len(name) == 0 {
				name = "GoWechat User"
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
	*/
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
	/*
		friendNumber, _ := this._wechat.FriendByPublicKey(who)
		len, err := this._wechat.FriendSendMessage(friendNumber, msg)
		if err != nil {
			log.Println(err, len)
			return -1
		}
		return int(len)
	*/
	return 0
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
	/*
		friendNumber, err := this._wechat.FriendByPublicKey(who)
		if err != nil {
			log.Println(err, friendNumber, who)
		}
		friendName, err := this._wechat.FriendGetName(friendNumber)
		friendStmsg, err := this._wechat.FriendGetStatusMessage(friendNumber)
		seen, err := this._wechat.FriendGetLastOnline(friendNumber)

		uinfo := purple.NewNotifyUserInfo()
		uinfo.AddPair("nickname", friendName)
		uinfo.AddPair("status message", friendStmsg)
		uinfo.AddPair("seen", fmt.Sprintf("%d", seen))
		uinfo.AddPair("hehehe", "efffff")
		uinfo.AddPair("hehehe12", "efffff456")

		gc.NotifyUserInfo(who, uinfo, func(ud interface{}) {
			log.Println("closed", ud)
		}, 123)
	*/
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
