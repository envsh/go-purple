package main

import (
	"sync"

	"github.com/thoj/go-ircevent"
)

const (
	PROTO_NONE = iota
	PROTO_IRC
	PROTO_TOX
)

var busch = make(chan interface{}, 123)

type Account struct {
	proto int
	ircon *irc.Connection
}

type AccountPool struct {
	mtx sync.Mutex
	acs map[string]*Account
}

func NewAccountPool() *AccountPool {
	this := &AccountPool{}
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

	ircon.AddCallback("*", func(e *irc.Event) {
		busch <- e
	})

	ac := &Account{}
	ac.ircon = ircon
	this.acs[name] = ac

	go ac.ircon.Connect(serverssl)

	return ac
}

func (this *AccountPool) remove(name string) {
	ac, _ := this.acs[name]
	delete(this.acs, name)
	ac.ircon.ClearCallback("*")
	ac.ircon.Disconnect()
}
