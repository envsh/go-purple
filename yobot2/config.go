package main

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
