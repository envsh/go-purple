package bridges

import (
	"log"
	"regexp"
	"strings"
)

func IsBotUser(senderc string) bool {
	sender := &senderc
	botUsers := []string{"teleboto", "Orizon",
		"OrzGTalk", "OrzIrc2P", "xmppbot"}
	for _, u := range botUsers {
		if *sender == u || strings.TrimRight(*sender, "_^") == u {
			return true
		}
	}
	return false
}

// color: 名字的颜色
func ExtractRealUser(sender, message string) (
	new_sender, new_message string, color string) {
	// teleboto <FONT COLOR="teal">[ngkaho1234] </FONT>後來修復了 2
	// [<FONT COLOR="pink">dant mnf</FONT>] 没有
	// TODO tg2offtopic ???
	// TODO offbot ???
	new_sender = sender
	new_message = message

	// 带FONT的，可能已经是pidgin转换过了的。
	msgregs := []string{
		`^<FONT COLOR="([\w ]+)">\[(.+)\] </FONT>`, // teleboto
		`^\[<FONT COLOR="([\w ]+)">(.+)</FONT>\]`,  // Orizon
		`^(.[0-9]+)\[(.+)\] `,                      // teleboto? with 1st unprintable char
		`^()\(GTalk\) (.+):`,                       // OrzGTalk
		`^()\[(.+)\] `,                             // OrzIrc2P/tg2offtopic
	}

	for idx, msgreg := range msgregs {
		exp := regexp.MustCompile(msgreg)

		mats := exp.FindAllStringSubmatch(message, -1)
		if len(mats) > 0 {
			log.Println("match reg:", idx, len(mats), msgreg,
				len(message), "\""+message+"\"")
		}
		if len(mats) == 1 {
			new_sender = mats[0][0] // with color style
			new_sender = mats[0][2]
			new_message = message[len(mats[0][0]):]
			color = mats[0][1]
			return
		}
	}
	return
}
