package main

import (
	"log"

	"go-purple/purple"
)

// 管理后端的网络连接账号
// 由于这种使用libpurple的版本，基本要实现完整的libpurple客户端，复杂度非常高
// 暂停使用这种方式，还是改用原始的协议拼装好了。

type AccountManager struct {
	pc  *purple.PurpleCore
	acs map[int64]*purple.Account
	acn int64
}

func NewAccountManager() *AccountManager {
	this := &AccountManager{}
	this.pc = gbot.pc
	this.acn = 0
	this.acs = make(map[int64]*purple.Account, 0)
	purple.SignedOn = this.onSingedOn

	return this
}

func (this *AccountManager) nextNo() int64 {
	this.acn += 1
	return this.acn
}

func (this *AccountManager) init() {
	purple.BlistInit()
	purple.BlistLoad()
	purple.BlistGet()

	blst := purple.BlistNew()
	purple.BlistSet(blst)

	if false {
		ac := this.addAccount("yobot", "prpl-irc")
		this.acs[this.nextNo()] = ac
	}
	{
		// ac := this.addAccount("yobot", "prpl-gotox")
		ac := this.addAccount("gotox-01", "prpl-gotox")
		this.acs[this.nextNo()] = ac
	}
	purple.BlistScheduleSave()
}

func (this *AccountManager) addAccount(username, protocol string) *purple.Account {
	switch protocol {
	case "prpl-irc":
		return this.addIrc(username)
	case "prpl-gotox":
		return this.addTox(username)
	}

	log.Fatalln("unsupported protocol:", protocol)
	return nil
}

func (this *AccountManager) addIrc(username string) *purple.Account {
	acc := this.pc.AccountsFind(username, "prpl-irc")
	log.Println(acc)
	if acc == nil {
		acc = purple.NewAccountCreate(username, "prpl-irc", "")
		log.Println(acc)
	}
	acc.SetBool("ssl", true)
	acc.SetInt("port", 6697)

	acc.SetEnabled(true)
	acc.Connect()
	return acc
}

func (this *AccountManager) addTox(username string) *purple.Account {
	acc := this.pc.AccountsFind(username, "prpl-gotox")
	log.Println(acc)
	if acc == nil {
		if false {
			acc = purple.NewAccountCreate(username, "prpl-gotox", "")
			log.Println(acc)
		}
		log.Fatalln("need precreate one")
	}

	acc.SetEnabled(true)
	if false {
		ht := purple.NewGHashTable()
		ht.Insert("ToxChannel", "#tox-toxmytest-helper")
		ht.Insert("GroupNumber", "0")
		chat := acc.ChatNew("#tox-toxmytest-helper", ht)
		chat.BlistAddChat(nil, nil)
		log.Println(chat)
	}

	acc.Connect()
	return acc
}

func (this *AccountManager) onSingedOn(gc *purple.Connection) {
	ac := gc.ConnGetAccount()
	log.Println(gc, ac.GetUserName(), ac.GetAlias())

	prplName := gc.GetPrplInfo().Id
	log.Println(gc, ac.GetUserName(), ac.GetAlias(), prplName)

	switch prplName {
	case "prpl-irc":
		hash := purple.NewGHashTable()
		hash.Insert("channel", "#tox-cn123")
		hash.Insert("password", "")
		gc.ServJoinChat(hash)
	case "prpl-gotox":
		if false {
			conv := purple.NewConversation(purple.CONV_TYPE_IM, ac, "hehhe")
			conv.GetChatData().Write("asfsdfsdf", "eaewfefweafewf", 0)
		}
	}
}
