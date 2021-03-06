package main

import (
	"log"
	// "sync"

	"github.com/sasha-s/go-deadlock"
)

const (
	PROTO_NONE  = "none"
	PROTO_SYS   = "sys"
	PROTO_IRC   = "irc"
	PROTO_TOX   = "tox"
	PROTO_TABLE = "table"
)

type AccountBase struct {
	Account
}

type AccountI interface {
}

// TODO Account 与 BackendBase 功能重叠
type Account struct {
	proto int
	// ircon  *irc.Connection
	uid    string
	name   string
	becon  Backend
	conque chan *Event
}

type AccountPool struct {
	RelaxCallObject
	ctx *Context
	mu  deadlock.RWMutex // sync.Mutex
	acs map[string]*Account
}

func NewAccountPool() *AccountPool {
	this := &AccountPool{}
	this.ctx = ctx
	this.acs = make(map[string]*Account)
	return this
}

func (this *AccountPool) has(name string, id string) bool {
	this.mu.RLock()
	defer this.mu.RUnlock()
	if _, ok := this.acs[id]; ok {
		return true
	}
	return false
}

func (this *AccountPool) get(name string, uid string) *Account {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if ac, ok := this.acs[uid]; ok {
		return ac
	}
	return nil
}

func (this *AccountPool) add(name string, uid string) *Account {
	this.mu.Lock()
	defer this.mu.Unlock()

	be := NewIrcBackend2(this.ctx, name, uid)
	be.uid = uid

	ac := &Account{name: name, uid: uid}
	// ac.ircon = ircon
	ac.becon = be
	ac.conque = make(chan *Event, MAX_BUS_QUEUE_LEN)
	this.acs[uid] = ac

	be.connect() // maybe block

	return ac
}

func (this *AccountPool) remove(name string, uid string) {
	this.mu.Lock()
	defer this.mu.Unlock()

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
func (this *AccountPool) disconnectAll() {
	this.mu.RLock()
	uids := make([]string, 0)
	for uid, _ := range this.acs {
		uids = append(uids, uid)
	}
	this.mu.RUnlock()

	for _, uid := range uids {
		ac := this.acs[uid]
		this.remove(ac.name, uid)
	}
}
func (this *AccountPool) count() int {
	this.mu.RLock()
	defer this.mu.RUnlock()

	return len(this.acs)
}

func (this *AccountPool) getNames(max int) []string {
	this.mu.RLock()
	defer this.mu.RUnlock()

	names := make([]string, 0)
	for _, ac := range this.acs {
		names = append(names, ac.name)
		if len(names) >= max && max > 0 {
			break
		}
	}
	return names
}

// 检测消息的nick是否是我们自己创建的
func (this *AccountPool) isOurs(nick string) bool {
	for _, ac := range this.acs {
		be := ac.becon.(*IrcBackend2)
		beme := be.ircon.Me()
		if nick == beme.Nick {
			// if suffix with ^^^, the beme.Nick contains it, and nick contains it too.
			// log.Printf("drop by my ourcon:%s, %s, %v\n", nick, be.getName(), beme)
			return true
		}
		if len(beme.Nick) > 16 && nick == beme.Nick[0:16] {
			return true
		}
	}
	return false
}
