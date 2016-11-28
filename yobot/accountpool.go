package main

import (
	"log"
	"sync"

	"github.com/thoj/go-ircevent"
)

const (
	PROTO_NONE = iota
	PROTO_IRC
	PROTO_TOX
)

type Account struct {
	proto int
	ircon *irc.Connection
}

type AccountPool struct {
	ctx *Context
	mtx sync.Mutex
	acs map[string]*Account
}

func NewAccountPool() *AccountPool {
	this := &AccountPool{}
	this.ctx = ctx
	this.acs = make(map[string]*Account)
	return this
}

func (this *AccountPool) has(name string) bool {
	if _, ok := this.acs[name]; ok {
		return true
	}
	return false
}

func (this *AccountPool) get(name string) *Account {
	if ac, ok := this.acs[name]; ok {
		return ac
	}
	return nil
}

func (this *AccountPool) add(name string) *Account {
	ircon := irc.IRC(name, name)
	ircon.VerboseCallbackHandler = false
	ircon.UseTLS = true
	ircon.Debug = false

	ircon.AddCallback("*", func(e *irc.Event) { this.ctx.busch <- e })

	ac := &Account{}
	ac.ircon = ircon
	this.acs[name] = ac

	go func() {
		err := ac.ircon.Connect(serverssl)
		if err != nil {
			log.Println(err)
		}
	}()

	return ac
}

func (this *AccountPool) remove(name string) {
	ac, _ := this.acs[name]
	delete(this.acs, name)
	ac.ircon.ClearCallback("*")
	if ac.ircon.Connected() {
		ac.ircon.Disconnect()
	}
}
