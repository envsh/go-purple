package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"go-purple/purple"

	"github.com/kitech/go-toxcore"
)

// TODO 也许这些附加功能可以做成一种插件式的，脚本式的呢？
var groupbot = "56A1ADE4B65B86BCD51CC73E2CD4E542179F47959FE3E0E21B4B0ACDADE51855D34D34D37CB5"

var addFixedFriendOnce sync.Once

func tryAddFixedFriends(t *tox.Tox, gc *purple.Connection) {
	addFixedFriendOnce.Do(func() {
		message := fmt.Sprintf("hello from gotox-prpl.%s", t.SelfGetAddress()[0:5])
		friendId := groupbot

		// check friend existance first
		friendNumber, err := t.FriendByPublicKey(friendId)
		if err == nil {
			return
		}

		friendNumber, err = t.FriendAdd(friendId, message)
		// log.Println(friendNumber, err)
		if err != nil {
			if false {
				log.Println(err, friendNumber)
			}
		} else {
			buddy := gc.ConnGetAccount().FindBuddy(friendId)
			// log.Println(buddy, friendNumber, err)
			if buddy == nil {
				// 由于是后端添加，所以需要显式加到好友表中
				// 会不会这个AddBudy就会触发调用ToxPlugin.AddBuddyWithInvite呢？
				// 如果是的话，本函数前面的添加好友步骤就是多余
				buddy := purple.NewBuddy(gc.ConnGetAccount(), friendId, "")
				gc.ConnGetAccount().AddBuddy(buddy)
				buddy.BlistAdd(nil)
			}
		}
	})
}

func tryJoinFixedGroups(t *tox.Tox, gc *purple.Connection, friendNumber uint32, status uint32) {
	if status == tox.CONNECTION_NONE {
		return
	}

	pubkey, err := t.FriendGetPublicKey(friendNumber)
	if err != nil {
	}
	if !strings.HasPrefix(groupbot, pubkey) {
		return
	}

	fixed := map[string]string{
		"#tox-toxen-helper":     "invite 0",
		"#tox-toxcn-helper":     "invite 2",
		"#tox-toxmytest-helper": "invite 5",
	}
	for name, handler := range fixed {
		chat := gc.ConnGetAccount().BlistFindChat(name)
		if chat == nil {
			log.Println("not found", name)
		} else {
			log.Println("found", name, chat)
			ht := chat.GetComponents()
			log.Println(ht.ToMap())
			log.Println(chat.Node().Settings().ToMap())
			if chat.Node().GetBool("gtk-autojoin") {
			}
			_, err := t.FriendSendMessage(friendNumber, handler)
			if err != nil {
				log.Println(err)
			}
		}
	}

	buddies := gc.ConnGetAccount().FindBuddies("")
	log.Println(len(buddies), buddies)
	for _, buddy := range buddies {
		log.Println(buddy, buddy.GetName(), buddy.GetAliasOnly())
	}

	convs := purple.GetConversations()
	for _, conv := range convs {
		log.Println(conv.GetName(), conv.GetChatData().GetUsers(), conv.GetData("GroupNumber"))
	}
	if len(convs) == 0 {
		log.Println("can not find conv chat:")
	}
}

func joinChatSpecialFixed(t *tox.Tox, comp *purple.GHashTable) bool {
	title := comp.Lookup("ToxChannel")
	fixed := map[string]bool{"#tox-toxmytest-helper": true, "#tox-toxen-helper": true}
	if _, ok := fixed[title]; ok {
		friendNumber, err := t.FriendByPublicKey(groupbot)
		if err != nil {
			log.Println(err, friendNumber)
		} else {
			_, err = t.FriendSendMessage(friendNumber, "invite 5")
			if err != nil {
				log.Println(err)
			}
		}
		return true
	}
	return false
}
