package gloox

/*
#cgo pkg-config: gloox
#cgo CFLAGS: -std=c99
#cgo CXXFLAGS: -std=c++14
#cgo LDFLAGS: -std=c++14
*/
import "C"

import (
	"log"

	"github.com/kitech/colog"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	colog.SetFlags(log.Flags())
	colog.Register()
}
