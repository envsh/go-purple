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

var fixedGroups = map[string]string{
	"#tox-en-nail":     "invite 0",
	"#tox-cn-nail":     "invite 2",
	"#tox-mytest-nail": "invite 5",
}

func isFixedGroup(name string) bool {
	for n, _ := range fixedGroups {
		if name == n {
			return true
		}
	}
	return false
}

func tryJoinFixedGroups(t *tox.Tox, gc *purple.Connection, friendNumber uint32, status int) {
	if status == tox.CONNECTION_NONE {
		return
	}

	pubkey, err := t.FriendGetPublicKey(friendNumber)
	if err != nil {
	}
	if !strings.HasPrefix(groupbot, pubkey) {
		return
	}

	fixed := fixedGroups
	for name, handler := range fixed {
		chat := gc.ConnGetAccount().BlistFindChat(name)
		if chat == nil {
			log.Println("not found chat:", name)
		} else {
			log.Println("found chat:", name, chat)
			ht := chat.GetComponents()
			log.Println("chat comps:", ht.ToMap())
			log.Println("chat settings:", chat.Node().Settings().ToMap())
			if chat.Node().GetBool("gtk-autojoin") {
			}
			_, err := t.FriendSendMessage(friendNumber, handler)
			if err != nil {
				log.Println(err)
			}
		}
	}

	buddies := gc.ConnGetAccount().FindBuddies("")
	log.Println("all buddies:", len(buddies), buddies)
	for _, buddy := range buddies {
		log.Println(buddy, buddy.GetAliasOnly(), buddy.GetName())
	}

	convs := purple.GetConversations()
	for _, conv := range convs {
		log.Println(conv.GetName(), conv.GetChatData().GetUsers(), conv.GetData("GroupNumber"))
	}
	if len(convs) == 0 {
		log.Println("can not find any conv chat.")
	}
}

// 点击聊天时的事件特殊处理
func joinChatSpecialFixed(t *tox.Tox, comp *purple.GHashTable) bool {
	title := comp.Lookup("_ToxChannel")
	fixed := fixedGroups
	if cmd, ok := fixed[title]; ok {
		friendNumber, err := t.FriendByPublicKey(groupbot)
		if err != nil {
			log.Println(err, friendNumber)
		} else {
			_, err = t.FriendSendMessage(friendNumber, cmd)
			if err != nil {
				log.Println(err)
			}
		}
		return true
	}
	return false
}
