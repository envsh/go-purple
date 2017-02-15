/*
   tox-prpl core code, minimal must implemention code.
*/
package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go-purple/purple"

	"github.com/kitech/colog"
	"github.com/kitech/go-toxcore"
)

type ToxPlugin struct {
	ppi      *purple.PluginProtocolInfo
	pi       *purple.PluginInfo
	p        *purple.Plugin
	_tox     *tox.Tox
	_toxav   *tox.ToxAV
	_toxopts *tox.ToxOptions
	stopch   chan struct{}

	iterTimerHandler int

	fileTransferFields
}

// plugin functions
func (this *ToxPlugin) init_tox(p *purple.Plugin) {
	log.Println("called", purple.GoID())

	var ao *purple.AccountOption

	ao = purple.NewAccountOptionString("Server", "server-ip", "")
	p.AddProtocolOption(ao)
	ao = purple.NewAccountOptionInt("Port", "server-port", 33445)
	p.AddProtocolOption(ao)

	ao = purple.NewAccountOptionBool("Use TCP", "use-tcp", true)
	p.AddProtocolOption(ao)
	ao = purple.NewAccountOptionBool("Use IPV6", "use-ipv6", false)
	p.AddProtocolOption(ao)

	ao = purple.NewAccountOptionBool("Show Contact Change", "show-contact-change", false)
	p.AddProtocolOption(ao)
	ao = purple.NewAccountOptionBool("Fake Offline Message", "fake-offline-message", true)
	p.AddProtocolOption(ao)
	ao = purple.NewAccountOptionBool("Auto accept file", "auto-accept-file", false)
	p.AddProtocolOption(ao)

	ao = purple.NewAccountOptionString("status text", "status-text", "hohohoh status text")
	p.AddProtocolOption(ao)

	ao = purple.NewAccountOptionBool("Send Typing", "send-typing", true)
	p.AddProtocolOption(ao)
	ao = purple.NewAccountOptionBool("Save Chat History", "save-chat-history", true)
	p.AddProtocolOption(ao)

	ao = purple.NewAccountOptionString("NoSpam", "nospam", "")
	p.AddProtocolOption(ao)

}

// should malloc some resource for use?
// and what resource here allocated is acceptable
func (this *ToxPlugin) load_tox(p *purple.Plugin) bool {
	log.Println("called", purple.GoID())
	rand.Seed(time.Now().UnixNano())
	this.setupModuleFields()
	return true
}

func (this *ToxPlugin) unload_tox(p *purple.Plugin) bool {
	log.Println("called")
	return true
}

func (this *ToxPlugin) destroy_tox(p *purple.Plugin) {
	log.Println("called")
}

// protocol functions, must implemented
func (this *ToxPlugin) tox_blist_icon() string {
	return "gotox"
}

func (this *ToxPlugin) tox_status_types(ac *purple.Account) []*purple.StatusType {
	stys := []*purple.StatusType{
		purple.NewStatusType(purple.STATUS_AVAILABLE, "tox_online", "Online", true),
		purple.NewStatusType(purple.STATUS_AWAY, "tox_away", "Away", true),
		purple.NewStatusType(purple.STATUS_UNAVAILABLE, "tox_busy", "Busy", true),
		purple.NewStatusType(purple.STATUS_OFFLINE, "tox_offline", "Offline", true),
	}
	return stys
}

var bsnodes = []string{
	"biribiri.org", "33445", "F404ABAA1C99A9D37D61AB54898F56793E1DEF8BD46B1038B9D822E8460FAB67",
	"178.62.250.138", "33445", "788236D34978D1D5BD822F0A5BEBD2C53C64CC31CD3149350EE27D4D9A2F9B6B",
	"205.185.116.116", "33445", "A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702",
}

func (this *ToxPlugin) tox_login(ac *purple.Account) {
	log.Println("called", purple.GoID())
	this.stopch = make(chan struct{}, 0)
	this._toxopts = tox.NewToxOptions()
	this._toxopts.Tcp_port = uint16(rand.Uint32()%55536) + 10000
	this.load_account(ac)

	// retry 50 times
	for port := 0; port < 50; port++ {
		this._toxopts.Tcp_port = uint16(rand.Uint32()%55536) + 10000
		this._tox = tox.NewTox(this._toxopts)
		if this._tox != nil {
			log.Println("TOXID:", this._tox.SelfGetAddress())
			break
		}
	}
	if this._tox == nil {
		log.Panicln("null")
	}

	for i := 0; i < len(bsnodes); i += 3 {
		port, _ := strconv.Atoi(bsnodes[i+1])
		ok1, err1 := this._tox.Bootstrap(bsnodes[i], uint16(port), bsnodes[i+2])
		ok2, err2 := this._tox.AddTcpRelay(bsnodes[i], uint16(port), bsnodes[i+2])
		if !ok1 || !ok2 || err1 != nil || err2 != nil {
			log.Println(ok1, ok2, err1, err2)
		}
	}

	// little state setup
	if true {
		conn := ac.GetConnection()
		conn.ConnSetState(purple.CONNECTING)
	}

	this.setupSelfInfo(ac)
	this.setupCallbacks(ac)
	this.setupFileCallbacks(ac)
	this.loadFriends(ac)
	this.save_account(ac.GetConnection())

	if false {
		buddy := purple.NewBuddy(ac, "onlyyou-id", "onlyyou-nick")
		ac.AddBuddy(buddy)
		group := purple.NewGroup("GOTOX")
		buddy.BlistAdd(group)
	}

	// go this.Iterate()
	log.Println(purple.GoID())
	this.iterTimerHandler = purple.TimeoutAdd(100, this, this.itercb)
}

func (this *ToxPlugin) tox_close(gc *purple.Connection) {
	// this.stopch <- struct{}{}
	ok := purple.TimeoutRemove(this.iterTimerHandler)
	if !ok {
		log.Println("rm timer failed")
	}
	this.save_account(gc)
	// TODO might have pending callback???
	this.ppi.StatusText = nil
	this._tox.Kill()
	this._tox = nil
	this._toxopts = nil
}

////////
func (this *ToxPlugin) itercb(d interface{}) {
	// log.Println(purple.GoID())
	this._tox.Iterate()
}

// should block and new thread
func (this *ToxPlugin) Iterate() {
	stopped := false
	tick := time.Tick(100 * time.Millisecond)
	id := this._tox.SelfGetAddress()
	for !stopped {
		select {
		case <-tick:
			this._tox.Iterate()
		case <-this.stopch:
			stopped = true
		}
	}
	log.Println("stopped", id)
}

func (this *ToxPlugin) load_account(ac *purple.Account) {
	data64 := ac.GetString("tox_save_data")
	data, err := base64.StdEncoding.DecodeString(data64)
	if err != nil || len(data) == 0 {
		log.Println("load data error:", err, data64)
	} else {
		this._toxopts.Savedata_data = data
		this._toxopts.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
	}
}

func (this *ToxPlugin) save_account(gc *purple.Connection) {
	data := this._tox.GetSavedata()
	data64 := base64.StdEncoding.EncodeToString(data)
	ac := gc.ConnGetAccount()
	ac.SetString("tox_save_data", data64)
}

func NewToxPlugin() *ToxPlugin {
	this := &ToxPlugin{}

	pi := purple.PluginInfo{
		Id:          "prpl-gotox",
		Name:        "GoTox",
		Version:     "1.0",
		Summary:     "Tox Protocol Plugin using golang",
		Description: "Tox Protocol Plugin https://tox.chat/",
		Author:      "it's gzleo",
		Homepage:    "https://github.com/kitech/go-purple/tox-prpl/",

		Load:    this.load_tox,
		Unload:  this.unload_tox,
		Destroy: this.destroy_tox,
	}
	ppi := purple.PluginProtocolInfo{
		Options: purple.OPT_PROTO_CHAT_TOPIC |
			purple.OPT_PROTO_INVITE_MESSAGE | purple.OPT_PROTO_PASSWORD_OPTIONAL,
		IconSpec: purple.BuddyIconSpec{Format: "png,jpg,jpeg",
			MinWidth: 16, MinHeight: 16, MaxWidth: 96, MaxHeight: 96,
			MaxFilesize: 0, ScaleRules: purple.ICON_SCALE_DISPLAY | purple.ICON_SCALE_SEND},
		BlistIcon:   this.tox_blist_icon,
		Login:       this.tox_login,
		Close:       this.tox_close,
		StatusTypes: this.tox_status_types,
		SendIM:      this.SendIM,
		// group chat
		ChatInfo:           this.ChatInfo,
		ChatInfoDefaults:   this.ChatInfoDefaults,
		JoinChat:           this.JoinChat,
		RejectChat:         this.RejectChat,
		GetChatName:        this.GetChatName,
		ChatInvite:         this.ChatInvite,
		ChatLeave:          this.ChatLeave,
		ChatWhisper:        this.ChatWhisper,
		ChatSend:           this.ChatSend,
		RoomlistGetList:    this.RoomlistGetList,
		AddBuddyWithInvite: this.AddBuddyWithInvite,
		RemoveBuddy:        this.RemoveBuddy,
		GetInfo:            this.GetInfo,
		StatusText:         this.StatusText,
		SetChatTopic:       this.SetChatTopic,
		Normalize:          this.Normalize,
		// other more
		SendTyping: this.SendTyping,
		// file transfer
		CanReceiveFile: this.CanReceiveFile,
		SendFile:       this.SendFile,
		NewXfer:        this.NewXfer,
	}
	this.pi = &pi
	this.ppi = &ppi
	this.p = purple.NewPlugin(&pi, &ppi, this.init_tox)

	return this
}

func init() {
	colog.Register()
	colog.SetFlags(log.LstdFlags | log.Lshortfile | colog.Flags())

	NewToxPlugin()
}

func main() {}
