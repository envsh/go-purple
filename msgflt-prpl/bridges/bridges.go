package bridges

import (
	"log"
	"regexp"
)

func ExtractRealUser(sender, message string) (new_sender, new_message string, color string) {
	// teleboto <FONT COLOR="teal">[ngkaho1234] </FONT>後來修復了 2
	// [<FONT COLOR="pink">dant mnf</FONT>] 没有

	new_sender = sender
	new_message = message

	msgregs := []string{
		`^<FONT COLOR="([\w ]+)">\[(.+)\] </FONT>`, // teleboto
		`^\[<FONT COLOR="([\w ]+)">(.+)</FONT>\]`,  // Orizon
		`^()\(GTalk\) (.+):`,                       // OrzGTalk
		`^()\[(.+)\] `,                             // OrzIrc2P
	}

	for _, msgreg := range msgregs {
		exp := regexp.MustCompile(msgreg)

		mats := exp.FindAllStringSubmatch(message, -1)
		if false {
			log.Println(mats, len(mats))
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
