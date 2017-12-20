package main

import (
	"fmt"
	"gopp"
	"log"
	"strings"

	tox "github.com/kitech/go-toxcore"
	"github.com/kitech/go-toxcore/xtox"
)

// 废群自动删除
// 自动加入所有群组
// 自动接受toxync的进群邀请
// 自动重连
var tyxdoCtx = xtox.NewToxContext("tyxdo.tsbin",
	"tyxdo", "shadow of toxync for try to keep group online")
var tyxdoFeat = xtox.FOTA_ACCEPT_FRIEND_REQUEST | xtox.FOTA_ACCEPT_GROUP_INVITE |
	xtox.FOTA_ADD_NET_HELP_BOTS | xtox.FOTA_REMOVE_ONLY_ME_ALL

type tyxdoTox struct {
	t *tox.Tox
}

func tyxdoToxNew() *tyxdoTox {
	this := &tyxdoTox{}

	this.t = xtox.New(tyxdoCtx)
	log.Println("ID:", this.t.SelfGetAddress())

	xtox.SetAutoBotFeatures(this.t, tyxdoFeat)
	this.initCallbacks()
	xtox.Connect(this.t)

	return this
}

// should block
func (this *tyxdoTox) run() { xtox.Run(this.t) }

var toxyncId = "96BE12EB9AF2F6851704412FA2981E03E32BBD18D40F6040D01F107A20CACC07639D8D4F2A94"

//
func (this *tyxdoTox) initCallbacks() {
	t := this.t

	newaddfrnd := false
	// add toxync as friend
	t.CallbackSelfConnectionStatusAdd(func(_ *tox.Tox, status int, userData interface{}) {
		log.Println(status, tox.ConnStatusString(status))
		_, err := t.FriendByPublicKey(toxyncId)
		gopp.ErrPrint(err)
		if err != nil {
			fn, err := t.FriendAdd(toxyncId, "I'm "+tyxdoCtx.NickName)
			gopp.ErrPrint(err, fn)
			newaddfrnd = true
			log.Println("added friend:", fn, toxyncId)
		}
	}, nil)

	//
	t.CallbackFriendConnectionStatusAdd(func(_ *tox.Tox, friendNumber uint32, status int, userData interface{}) {
		pubkey, err := t.FriendGetPublicKey(friendNumber)
		gopp.ErrPrint(err, pubkey)
		log.Println(friendNumber, status, tox.ConnStatusString(status), pubkey)
		if strings.HasPrefix(toxyncId, pubkey) && newaddfrnd {
			for i := 0; i < 32; i++ {
				_, err := t.FriendSendMessage(friendNumber, fmt.Sprintf("join %d", i))
				gopp.ErrPrint(err, i)
			}
			newaddfrnd = false
			log.Println("joined 30 groups")
		}
	}, nil)
}
