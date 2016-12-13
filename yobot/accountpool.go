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
	be := NewIrcBackend2(this.ctx, name)

	ac := &Account{}
	// ac.ircon = ircon
	ac.becon = be
	ac.conque = make(chan *Event, MAX_BUS_QUEUE_LEN)
	this.acs[name] = ac

	be.connect()

	return ac
}

func (this *AccountPool) remove(name string) {
	if ac, ok := this.acs[name]; ok {
		delete(this.acs, name)
		if ac.becon.isconnected() {
			ac.becon.disconnect()
		}
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

func (this *AccountPool) count() int {
	return len(this.acs)
}

func (this *AccountPool) getNames(max int) []string {
	names := make([]string, 0)
	for n, _ := range this.acs {
		names = append(names, n)
		if len(names) >= max && max > 0 {
			break
		}
	}
	return names
}
