package main

import (
	"log"
	"testing"
)

func TestFilter(t *testing.T) {
	f := NewMessageFilter()
	msgs := []string{
		"abc",
		"[teleboto] 3[PeterCxy] = = yubico-pam中Peter",
		"[teleboto] 10[ouyangjun1999] 感觉你们好rich",
	}

	for _, m := range msgs {
		nm := f.Filter(m)
		log.Println(m, "=>", nm)
	}
}
