package main

import (
	"fmt"
	"gopp"
	"log"
	"strings"
	"time"

	tox "github.com/kitech/go-toxcore"
)

type ToxVM struct {
	t    *tox.Tox
	opts *tox.ToxOptions
	oilC chan interface{} // out iterate loop
}

type NamePeerListChanged struct {
	t      *tox.Tox
	ud     interface{}
	gn     uint32
	pn     uint32
	change uint8
}

var tsname = "grp2mob.tsbin"
var instname = "Aegaeon"
var inststmsg = "Forward tox message between tox cn group to/from bot friend."

func newToxVM() *ToxVM {
	this := &ToxVM{}
	opts := tox.NewToxOptions()
	opts.ThreadSafe = true

	tscc, err := tox.LoadSavedata(tsname)
	if err == nil {
		opts.Savedata_data = tscc
		opts.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
		opts.Tcp_port = 55555
		opts.Udp_enabled = true
	} else {
		gopp.ErrPrint(err)
	}

	t := tox.NewTox(opts)
	t.SelfSetName(instname)
	t.SelfSetStatusMessage(inststmsg)
	this.t = t

	this.t.WriteSavedata(tsname)
	log.Println(this.t.SelfGetAddress())

	this.initCallbacks()
	this.bootstrap()
	this.oilC = make(chan interface{}, 0)
	go this.outIterateLoopHandler()
	return this
}

func (this *ToxVM) initCallbacks() {
	this.t.CallbackSelfConnectionStatus(func(_ *tox.Tox, status int, userData interface{}) {
		log.Println(status, tox.ConnStatusString(status))
		this.t.WriteSavedata(tsname)
	}, nil)

	this.t.CallbackConferenceAction(func(_ *tox.Tox, groupNumber uint32, peerNumber uint32, action string, userData interface{}) {
		log.Println(groupNumber, peerNumber, action)
	}, nil)

	this.t.CallbackFriendStatus(func(_ *tox.Tox, friendNumber uint32, status int, userData interface{}) {
		log.Println(friendNumber, status)
	}, nil)
	this.t.CallbackFriendConnectionStatus(func(_ *tox.Tox, friendNumber uint32, status int, userData interface{}) {
		log.Println(friendNumber, status, tox.ConnStatusString(status))
		this.t.WriteSavedata(tsname)
	}, nil)

	this.t.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, msg string, userData interface{}) {
		log.Println(friendNumber, msg)
		pubkey, err := t.FriendGetPublicKey(friendNumber)
		gopp.ErrPrint(err)
		_ = pubkey
		friendName, err := t.FriendGetName(friendNumber)
		gopp.ErrPrint(err)

		gns := t.ConferenceGetChatlist()
		for _, gn := range gns {
			t.ConferenceSendMessage(gn, tox.MESSAGE_TYPE_NORMAL,
				fmt.Sprintf("%s: %s", friendName, msg))
		}

		for i := uint32(0); i < 128; i++ {
			if i == friendNumber {
				continue
			}
			if !t.FriendExists(i) {
				break
			}
			st, err := t.FriendGetConnectionStatus(i)
			gopp.ErrPrint(err)
			if st > tox.CONNECTION_NONE {
				_, err := t.FriendSendMessage(i, fmt.Sprintf("%s: %s", friendName, msg))
				gopp.ErrPrint(err)
			}
		}
	}, nil)

	this.t.CallbackConferenceMessage(func(t *tox.Tox, groupNumber uint32, peerNumber uint32, message string, userData interface{}) {
		log.Println(groupNumber, peerNumber, message)
		peerPubkey, err := t.ConferencePeerGetPublicKey(groupNumber, peerNumber)
		gopp.ErrPrint(err)
		if strings.HasPrefix(t.SelfGetAddress(), peerPubkey) {
			return
		}
		peerName, err := t.ConferencePeerGetName(groupNumber, peerNumber)
		gopp.ErrPrint(err)
		for i := uint32(0); i < 128; i++ {
			if !t.FriendExists(i) {
				break
			}
			st, err := t.FriendGetConnectionStatus(i)
			gopp.ErrPrint(err)
			if st > tox.CONNECTION_NONE {
				_, err := t.FriendSendMessage(i, fmt.Sprintf("%s: %s", peerName, message))
				gopp.ErrPrint(err)
			}
		}
	}, nil)

	this.t.CallbackFriendRequest(func(t *tox.Tox, pubkey string, message string, userData interface{}) {
		_, err := t.FriendAddNorequest(pubkey)
		gopp.ErrPrint(err)
		if err != nil {
			t.WriteSavedata(tsname)
		}
	}, nil)

	this.t.CallbackConferenceInvite(func(t *tox.Tox, friendNumber uint32, itype uint8, data []byte, userData interface{}) {
		switch int(itype) {
		case tox.CONFERENCE_TYPE_TEXT:
			_, err := t.ConferenceJoin(friendNumber, data)
			gopp.ErrPrint(err)
		case tox.CONFERENCE_TYPE_AV:
			t.JoinAVGroupChat(friendNumber, data)
		}
		t.WriteSavedata(tsname)
	}, nil)

	this.t.CallbackConferenceNameListChange(func(t *tox.Tox, groupNumber uint32, peerNumber uint32, change uint8, userData interface{}) {
		this.oilC <- &NamePeerListChanged{t, userData, groupNumber, peerNumber, change}
	}, nil)
}

func (this *ToxVM) outIterateLoopHandler() {
	stop := false
	for !stop {
		select {
		case evtx := <-this.oilC:
			switch evt := evtx.(type) {
			case *NamePeerListChanged:
				t := evt.t
				groupNumber := evt.gn
				change := evt.change
				switch int(change) {
				case tox.CONFERENCE_STATE_CHANGE_PEER_EXIT:
					if t.ConferencePeerCount(groupNumber) == 1 {
						grpTitle, err := t.ConferenceGetTitle(groupNumber)
						gopp.ErrPrint(err)
						log.Println("only me left, leave now:", groupNumber, grpTitle)
						t.ConferenceDelete(groupNumber)
					}
				}
			}
		}
	}
}

func (this *ToxVM) bootstrap() {
	this.t.Bootstrap("194.249.212.109", 33445, "3CEE1F054081E7A011234883BC4FC39F661A55B73637A5AC293DDF1251D9432B")
	this.t.Bootstrap("130.133.110.14", 33445, "461FA3776EF0FA655F1A05477DF1B3B614F7D6B124F7DB1DD4FE3C08B03B640F")
}

func (this *ToxVM) isconnected() bool {
	return this.t.SelfGetConnectionStatus() > tox.CONNECTION_NONE
}

func (this *ToxVM) run() {
	tmer := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-tmer.C:
			this.t.Iterate2(nil)
		}
	}
}
