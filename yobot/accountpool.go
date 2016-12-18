package main

import (
	"log"
	"sync"
	// "github.com/thoj/go-ircevent"
)

const (
	PROTO_NONE  = "none"
	PROTO_IRC   = "irc"
	PROTO_TOX   = "tox"
	PROTO_TABLE = "table"
)

type Account struct {
	proto int
	// ircon  *irc.Connection
	uid    string
	name   string
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

func (this *AccountPool) has(name string, id string) bool {
	if _, ok := this.acs[id]; ok {
		return true
	}
	return false
}

func (this *AccountPool) get(name string, uid string) *Account {
	if ac, ok := this.acs[uid]; ok {
		return ac
	}
	return nil
}

func (this *AccountPool) add(name string, uid string) *Account {
	be := NewIrcBackend2(this.ctx, name)
	be.uid = uid

	ac := &Account{name: name, uid: uid}
	// ac.ircon = ircon
	ac.becon = be
	ac.conque = make(chan *Event, MAX_BUS_QUEUE_LEN)
	this.acs[uid] = ac

	be.connect()

	return ac
}

func (this *AccountPool) remove(name string, uid string) {
	if ac, ok := this.acs[uid]; ok {
		delete(this.acs, uid)
		if ac.becon.isconnected() {
			ac.becon.disconnect()
		}
	} else {
		log.Println("not found:", name, ac, uid)
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
	for _, ac := range this.acs {
		names = append(names, ac.name)
		if len(names) >= max && max > 0 {
			break
		}
	}
	return names
}
