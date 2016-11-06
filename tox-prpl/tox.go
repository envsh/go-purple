package main

import (
	"log"

	"yobot/purpleplugin"

	"github.com/kitech/colog"
)

type ToxPlugin struct {
	ppi *plugin.PluginProtocolInfo
	pi  *plugin.PluginInfo
	p   *plugin.Plugin
}

// plugin functions
func (this *ToxPlugin) init_tox(p *plugin.Plugin) {
	log.Println("called")
}

func (this *ToxPlugin) load_tox(p *plugin.Plugin) bool {
	log.Println("called")
	return true
}

func (this *ToxPlugin) unload_tox(p *plugin.Plugin) bool {
	log.Println("called")
	return true
}

func (this *ToxPlugin) destroy_tox(p *plugin.Plugin) {
	log.Println("called")
}

// protocol functions
func (this *ToxPlugin) tox_blist_icon() string {
	log.Println("called")
	return "gotox"
}

func (this *ToxPlugin) tox_login() {
	log.Println("called")
}

func (this *ToxPlugin) tox_close() {
	log.Println("called")
}

func (this *ToxPlugin) tox_status_types() {
	log.Println("called")
}

func NewToxPlugin() *ToxPlugin {
	this := &ToxPlugin{}

	pi := plugin.PluginInfo{
		Id:          "prpl-gotox",
		Name:        "GoTox",
		Version:     "1.0",
		Summary:     "it's summary",
		Description: "it's description",
		Author:      "it's gzleo",
		Homepage:    "https://fixlan.net/",

		Load:    this.load_tox,
		Unload:  this.unload_tox,
		Destroy: this.destroy_tox,
	}
	ppi := plugin.PluginProtocolInfo{
		BlistIcon: this.tox_blist_icon,
		Login:     this.tox_login,
		Close:     this.tox_close,
	}
	this.p = plugin.NewPlugin(&pi, &ppi, this.init_tox)

	return this
}

func init() {
	colog.Register()
	colog.SetFlags(log.LstdFlags | log.Lshortfile | colog.Flags())

	NewToxPlugin()
}

func main() {}
