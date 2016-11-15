package main

import (
	"flag"
	"log"

	"github.com/kitech/colog"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Flags())
	colog.Register()
}

func main() {
	flag.Parse()
	newPurpleServer().serve()
}
