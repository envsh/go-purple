/*
   tox-prpl core code, minimal must implemention code.
*/
package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	// "strconv"
	"time"

	"go-purple/purple"
	"go-purple/wechat-prpl/wechat"

	"github.com/kitech/colog"
)

type WechatPlugin struct {
	ppi     *purple.PluginProtocolInfo
	pi      *purple.PluginInfo
	p       *purple.Plugin
	_wechat *wechat.Wechat
	stopch  chan struct{}

	iterTimerHandler int
}

// plugin functions
func (this *WechatPlugin) init_wechat(p *purple.Plugin) {
	log.Println("called")

}

func (this *WechatPlugin) load_wechat(p *purple.Plugin) bool {
	log.Println("called")
	rand.Seed(time.Now().UnixNano())
	return true
}

func (this *WechatPlugin) unload_wechat(p *purple.Plugin) bool {
	log.Println("called")
	return true
}

func (this *WechatPlugin) destroy_wechat(p *purple.Plugin) {
	log.Println("called")
}

// protocol functions, must implemented
func (this *WechatPlugin) wechat_blist_icon() string {
	return "gowechat"
}

func (this *WechatPlugin) wechat_status_types(ac *purple.Account) []*purple.StatusType {
	stys := []*purple.StatusType{
		purple.NewStatusType(purple.STATUS_AVAILABLE, "wechat_online", "Online", true),
		purple.NewStatusType(purple.STATUS_OFFLINE, "wechat_offline", "Offline", true),
	}
	return stys
}

func (this *WechatPlugin) wechat_login(ac *purple.Account) {
	this.stopch = make(chan struct{}, 0)
	this._wechat = wechat.NewWechat()
	this._wechat.OnEvent = this.eventHandler
	// this.load_account(ac)

	// little state setup
	if true {
		conn := ac.GetConnection()
		conn.ConnSetState(purple.CONNECTING)
	}

	this.setupSelfInfo(ac)
	// this.setupCallbacks(ac)
	// this.loadFriends(ac)

	if false {
		buddy := purple.NewBuddy(ac, "onlyyou-id", "onlyyou-nick")
		ac.AddBuddy(buddy)
		group := purple.NewGroup("GOWECHAT")
		buddy.BlistAdd(group)
	}

	// go this.Iterate()
	this.iterTimerHandler = purple.TimeoutAdd(100, ac, this.itercb)
}

func (this *WechatPlugin) wechat_close(gc *purple.Connection) {
	// this.stopch <- struct{}{}
	purple.TimeoutRemove(this.iterTimerHandler)
	// this.save_account(gc)
	this._wechat.Kill()
	this._wechat = nil
}

////////
func (this *WechatPlugin) itercb(ud interface{}) {
	this._wechat.Iterate(ud)
}

// should block and new thread
func (this *WechatPlugin) Iterate() {
	stopped := false
	tick := time.Tick(100 * time.Millisecond)
	// id := this._wechat.SelfGetAddress()
	for !stopped {
		select {
		case <-tick:
			this._wechat.Iterate(this)
		case <-this.stopch:
			stopped = true
		}
	}
	log.Println("stopped" /*, id*/)
}

func (this *WechatPlugin) load_account(ac *purple.Account) {
	data64 := ac.GetString("wechat_save_data")
	data, err := base64.StdEncoding.DecodeString(data64)
	if err != nil || len(data) == 0 {
		log.Println("load data error:", err, data64)
	} else {
		// this._wechatopts.Savedata_data = data
		// this._wechatopts.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
	}
}

func (this *WechatPlugin) save_account(gc *purple.Connection) {
	/*
		data := this._wechat.GetSavedata()
		data64 := base64.StdEncoding.EncodeToString(data)
		ac := gc.ConnGetAccount()
		ac.SetString("wechat_save_data", data64)
	*/
}

func NewWechatPlugin() *WechatPlugin {
	this := &WechatPlugin{}

	pi := purple.PluginInfo{
		Id:          "prpl-wechat",
		Name:        "Wechat",
		Version:     "1.0",
		Summary:     "Wechat Protocol Plugin using golang",
		Description: "Wechat Protocol Plugin https://wx.qq.im/",
		Author:      "it's gzleo",
		Homepage:    "https://github.com/kitech/go-purple/wechat-prpl/",

		Load:    this.load_wechat,
		Unload:  this.unload_wechat,
		Destroy: this.destroy_wechat,
	}
	ppi := purple.PluginProtocolInfo{
		BlistIcon:   this.wechat_blist_icon,
		Login:       this.wechat_login,
		Close:       this.wechat_close,
		StatusTypes: this.wechat_status_types,
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
	}
	this.p = purple.NewPlugin(&pi, &ppi, this.init_wechat)

	return this
}

func init() {
	colog.Register()
	colog.SetFlags(log.LstdFlags | log.Lshortfile | colog.Flags())

	NewWechatPlugin()
}

func main() {}
