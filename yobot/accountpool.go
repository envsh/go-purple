package main

import (
	"log"
	"sync"
	// "github.com/thoj/go-ircevent"
)

const (
	PROTO_NONE = "none"
	PROTO_IRC  = "irc"
	PROTO_TOX  = "tox"
)

type Account struct {
	proto int
	// ircon  *irc.Connection
	becon  Backend
	conque chan *Event
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
	be := NewIrcBackend(this.ctx, name)

	ac := &Account{}
	// ac.ircon = ircon
	ac.becon = be
	ac.conque = make(chan *Event, 123)
	this.acs[name] = ac

	be.connect()

	return ac
}

func (this *AccountPool) remove(name string) {
	if ac, ok := this.acs[name]; ok {
		delete(this.acs, name)
	} else {
		log.Println("not found:", name, ac)
	}
	/*
		ac.ircon.ClearCallback("*")
		if ac.ircon.Connected() {
			ac.ircon.Disconnect()
		}
	*/
}
