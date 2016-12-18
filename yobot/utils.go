package main

import (
	"log"
	"net/url"
	"strings"
)

func isLocalUrl(u string) bool {
	uo, err := url.Parse(u)
	if err != nil {
		log.Println(err)
	} else {
		hostAndPort := strings.Split(uo.Host, ":")
		if len(hostAndPort) > 0 {
			if hostAndPort[0] == "127.0.0.1" || hostAndPort[0] == "localhost" {

				return true
			}
		}
	}
	return false
}
