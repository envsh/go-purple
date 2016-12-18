package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type MessageFilter struct {
}

func NewMessageFilter() *MessageFilter {
	return &MessageFilter{}
}

func (this *MessageFilter) Filter(message string) string {
	return this.filterColor(message)
}

func (this *MessageFilter) filterColor(message string) string {
	if false {
		chs := strings.Split(message, "")
		for idx, ch := range chs {
			log.Println(idx, len(ch), ch)
			if len(ch) == 1 && !strconv.IsPrint(rune(ch[0])) {
				chs[idx] = ""

			}
		}
		return strings.Join(chs, "")
	}

	exp := `([^[:print:]]([0-9]{0,2}))`
	reg, err := regexp.Compile(exp)
	if err != nil {
		log.Println(err, reg)
	}

	// log.Println(reg.MatchString(message))
	mats := reg.FindAllStringSubmatch(message, -1)
	// log.Println(len(mats), mats)
	for _, mat := range mats {
		if mat[0][0] >= 0 && mat[0][0] <= 15 {
			// log.Println(mat[0], len(mat[0]), []byte(mat[0]))
			message = strings.Replace(message, mat[0], "", -1)
		}
	}
	return message
}
