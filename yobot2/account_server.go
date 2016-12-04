package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"go-purple/purple"
)

type AccountServer struct {
	pc       *purple.PurpleCore
	cbs      purple.CoreCallbacks
	csigs    purple.CoreSignals
	requiops *purple.RequestUiOps
}

func NewAccountServer(pc *purple.PurpleCore) *AccountServer {
	this := &AccountServer{}
	this.pc = pc
	this.requiops = purple.NewRequestUiOpts()
	this.requiops.RequestAction = this.RequestAction
	purple.RequestSetUiOps(this.requiops)

	this.init()
	return this
}

func (this *AccountServer) init() {
	this.fillCallbacks()
	this.pc.SetCallbacks(this.cbs)
	this.pc.SetSignals(this.csigs)

	// account
	acs := purple.AccountsGetAll()
	log.Println(len(acs))
	for _, ac := range acs {
		ac.SetEnabled(false)
	}

	plugs := purple.PluginsGetProtocols()
	var pids = make(map[string]bool)
	for _, plug := range plugs {
		// log.Println(plug.GetName(), plug.GetId())
		pids[plug.GetId()] = true
	}

	// 删除配置中没有的账号
	for _, ac := range acs {
		userSig := fmt.Sprintf("%s&%s", ac.GetUserName(), ac.GetProtocolId())
		found := false
		for username, proto := range cfg.accounts {
			userSig2 := fmt.Sprintf("%s&prpl-%s", username, proto)
			if userSig2 == userSig {
				found = true
				break
			}
		}
		if !found {
			log.Println("remove unused account:", userSig)
			ac.AccountsRemove()
		}
	}
	acs = purple.AccountsGetAll()

	// 检查协议支持情况
	for username, proto := range cfg.accounts {
		protoName := fmt.Sprintf("prpl-%s", proto)
		if _, ok := pids[protoName]; !ok {
			log.Println("protocol not found:", protoName)
		}
		ac := purple.AccountsFind(username, protoName)
		if ac == nil {
			ac = purple.NewAccountCreate(username, protoName, "")
		}
		ac.SetEnabled(false) // too late, why?

		// more settings
		switch proto {
		case "irc":
			ac.SetInt("port", 6697) // ssl 6697/7000, no-ssl 8001
			ac.SetBool("ssl", true)
		}
	}

	// 连接账号
	purple.SavedStatusSetIdleAway(false)
	acs = purple.AccountsGetAll()
	for _, ac := range acs {
		ac.SetEnabled(true)
		ac.SetStatus("ONLINE", true)
		log.Println(ac.IsConnecting(), ac.IsConnected(), ac.IsDisconnected())
	}

	// 设置账号auto-login = false
	for _, ac := range acs {
		ac.SetUiBool("auto-login", false)
		// TODO
		log.Println(ac.IsConnecting(), ac.IsConnected(), ac.IsDisconnected())
	}

	// prefs set auto reply off. 为啥不管用
	purple.PrefsSetString("/purple/away/idle_reporting", "none")
	purple.PrefsSetString("/purple/away/auto_reply", "never")
	// purple.PrefsRemove("/purple/away/auto_reply")
	purple.PrefsSetBool("/purple/away/away_when_idle", false)
	purple.PrefsSetInt("/purple/away/mins_before_away", math.MaxInt32/2)

	// check all plugins
	allplugs := purple.PluginsGetAll()
	for _, plug := range allplugs {
		if false {
			log.Println(plug.GetId(), plug.GetName(), plug.GetSummary())
		}
	}
}

func (this *AccountServer) fillCallbacks() {
	this.csigs.SignedOn = this.onSignedOn
	this.csigs.SignedOff = this.onSignedOff
	this.csigs.ReceivedIMMsg = this.onReceivedImMsg
	this.csigs.ReceivedChatMsg = this.onReceivedChatMsg
	this.csigs.ChatJoined = this.onChatJoined
	this.csigs.ChatLeft = this.onChatLeft
}

func (this *AccountServer) RequestAction(title string, primary string, secondary string,
	default_action int, account *purple.Account, who string, conv *purple.Conversation,
	user_data interface{}, action_count int) interface{} {
	log.Println(title, primary, who)
	// purple.RequestActionCB(user_data, 1)
	return 1
}

func (this *AccountServer) onSignedOn(gc *purple.Connection) {
	pid := gc.GetPrplInfo().Id[5:]
	log.Println(pid)

	switch pid {
	case "irc":
		msg := fmt.Sprintf("hello you @ %s", time.Now())
		rc := gc.ServSendIM("kitech1", msg, 0)
		log.Println(rc)
	case "gotox":

	}
}

func (this *AccountServer) onSignedOff(gc *purple.Connection) {
	pid := gc.GetPrplInfo().Id[5:]
	ac := gc.ConnGetAccount()
	log.Println(pid, ac, ac.GetAlias())
	switch pid {
	case "irc":
		/*
			if ac.IsConnected() {
				ac.SetEnabled(false)
				ac.Disconnect()
			}
			if !ac.GetEnabled() {
				ac.SetEnabled(true)
			}
		*/
		log.Println(ac.IsConnecting(), ac.IsConnected(), ac.IsDisconnected(), ac.GetEnabled())
		ac.SetEnabled(true)
		// ac.Connect()
	}
}

func (this *AccountServer) onReceivedImMsg(ac *purple.Account, sender, msg string,
	conv *purple.Conversation, flags int) {
	gc := ac.GetConnection()
	pid := gc.GetPrplInfo().Id[5:]
	log.Println(pid, ac, ac.GetAlias())
	log.Println(ac, sender, msg, conv, flags, conv.GetName())

	switch pid {
	case "irc":
	case "gotox":

	}
}

func (this *AccountServer) onReceivedChatMsg(ac *purple.Account, sender, msg string,
	conv *purple.Conversation, flags int) {
	gc := ac.GetConnection()
	pid := gc.GetPrplInfo().Id[5:]
	log.Println(pid, ac, ac.GetAlias(), ac.GetUserName())
	log.Println(ac, sender, msg, conv, flags, conv.GetName())

	nmsg := fmt.Sprintf("%s: %s", sender, msg)
	switch pid {
	case "irc":
		if sender == strings.Split(ac.GetUserName(), "@")[0] {
			// log.Println("self msg, break")
			break // self msg, break
		}
		realConvName := conv.GetName()
		if v, ok := chmap.Get(conv.GetName()); ok {
			realConvName = v.(string)
		}
		convs := purple.GetConversations()
		for _, c := range convs {
			log.Println(c.GetName(), c.GetConnection().GetPrplInfo().Id)
			if c.GetConnection().GetPrplInfo().Id[5:] == "gotox" && c.GetName() == realConvName {
				log.Println("found", c)
				c.GetChatData().Send(nmsg)
				break
			}
		}

	case "gotox":
		if sender == ac.GetUserName() {
			// log.Println("self msg, break")
			break
		}
		acdst := purple.AccountsFind(cfg.getIrc(""), "prpl-irc")
		condst := acdst.GetConnection()
		if acdst == nil {
			log.Println("can't find:", cfg.getIrc(""))
		}
		if condst == nil {
			log.Println("conv dest nil")
		} else {
			log.Println(acdst.IsConnected(), acdst.IsDisconnected(), acdst.IsConnecting(), acdst.GetEnabled())
		}
		ht := purple.NewGHashTable()
		ht.Insert("channel", conv.GetName())
		if k, ok := chmap.GetKey(conv.GetName()); ok {
			ht.Insert("channel", k.(string))
		}
		condst.ServJoinChat(ht)

		realConvName := ht.Lookup("channel")
		convdst := purple.FindConversationWithAccount(purple.CONV_TYPE_CHAT, realConvName, acdst)
		if convdst == nil {
			convs := purple.GetConversations()
			for _, c := range convs {
				log.Println(c.GetName())
			}
		} else {
			// 不同的发送消息方式，区别在吗呢？
			if rand.Int()%2 == 0 {
				convdst.GetChatData().Send(nmsg + " from chat send")
			} else {
				chatid := convdst.GetChatData().GetId()
				condst.ServChatSend(chatid, nmsg+" from serv chat send", 0)
			}
			convdst.GetChatData().Write(sender, nmsg+" from chat write", 0) // ??? 发送不了消息？？？
		}
	}
}

func (this *AccountServer) onChatJoined(conv *purple.Conversation) {
	gc := conv.GetConnection()
	pid := gc.GetPrplInfo().Id[5:]
	log.Println(conv, conv.GetName(), pid)

	switch pid {
	case "irc":
	case "gotox":
		acdst := purple.AccountsFind(cfg.getIrc(""), "prpl-irc")
		condst := acdst.GetConnection()
		if acdst == nil {
			log.Println("can't find:", cfg.getIrc(""))
		}
		if condst == nil {
			log.Println("conv dest nil")
		} else {
			log.Println(acdst.IsConnected(), acdst.IsDisconnected(), acdst.IsConnecting(), acdst.GetEnabled())
		}
		ht := purple.NewGHashTable()
		ht.Insert("channel", conv.GetName())

		if k, ok := chmap.GetKey(conv.GetName()); ok {
			ht.Insert("channel", k.(string))
		}
		condst.ServJoinChat(ht)
	}
}
func (this *AccountServer) onChatLeft(conv *purple.Conversation) {
	gc := conv.GetConnection()
	pid := gc.GetPrplInfo().Id[5:]
	log.Println(conv, conv.GetName(), pid)
}

func (this *AccountServer) run() {
	select {}
}

func (this *AccountServer) hehe() {
	go func() {
		time.Sleep(2)
		acs := purple.AccountsGetAll()
		log.Println(len(acs))
		ac := purple.AccountsFind(username, "prpl-irc")
		if ac == nil {
			ac = purple.NewAccountCreate(username, "prpl-irc", "")
		}
		go func() {
			time.Sleep(3 * time.Second)
			ac.SetEnabled(false)
			ac.SetEnabled(true)
			ac.Connect()
		}()
		log.Println(ac, ac.GetUserName(), ac.GetAlias(), ac.GetEnabled())
		/*
			for _, ac := range acs {
				log.Println(ac.GetAlias(), ac.GetProtocolName(), ac.GetEnabled())
				ac.Connect()
			}
		*/
		go func() {
			time.Sleep(6000 * time.Second)
			gc := ac.GetConnection()
			log.Println(gc)
			msg := fmt.Sprintf("hello you @ %s", time.Now())
			rc := gc.ServSendIM("kitech1", msg, 0)
			log.Println(rc)
		}()
	}()
}
