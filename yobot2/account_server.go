package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"go-purple/purple"
)

type AccountServer struct {
	pc    *purple.PurpleCore
	cbs   purple.CoreCallbacks
	csigs purple.CoreSignals
}

func NewAccountServer(pc *purple.PurpleCore) *AccountServer {
	this := &AccountServer{}
	this.pc = pc

	this.init()
	return this
}

func (this *AccountServer) init() {
	this.csigs.SignedOn = func(gc *purple.Connection) {
		msg := fmt.Sprintf("hello you @ %s", time.Now())
		rc := gc.ServSendIM("kitech1", msg, 0)
		log.Println(rc)
	}
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
		log.Println(plug.GetId(), plug.GetName(), plug.GetSummary())
	}
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
