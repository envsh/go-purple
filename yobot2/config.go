package main

import (
	"go-purple/purple"

	"github.com/emirpasic/gods/maps/hashbidimap"
)

var username = "zuck03@weber.freenode.net"

type Config struct {
	accounts map[string]string // username => protocol
	chmap    map[string]string
}

var cfg = &Config{
	// root users
	accounts: map[string]string{
		"zuck03@weber.freenode.net": "irc",
		"zuck03":                    "gotox",
	},
	/*
		chmaps: map[string]map[string]string{
			"irc": map[string]string{"abc": "efg"},
			"tox": map[string]string{"abc": "efg"},
		},
	*/
}

var chmap = hashbidimap.New()

func init() {
	// irc <=> tox
	chmap.Put("#tox-cn123", "testks")
	chmap.Put("#tox-cn", "Chinese 中文")
	chmap.Put("#tox-en", "#tox")
}

func (this *Config) getIrc(from string) string {
	for n, p := range this.accounts {
		if p == "irc" {
			return n
		}
	}
	return ""
}

// 是否是config中定义的root用户。为以后自动创建的衍生用户做准备。
func isRootUser(ac *purple.Account) bool {
	return true
}
