package bridges

import (
	"log"
	"testing"
)

//
func TestFF(t *testing.T) {
	// TODO xmppbot
	msgs := []string{
		`<FONT COLOR="teal laet">[ngkaho1234] </FONT>後來修復了`, // teleboto
		`[<FONT COLOR="pink knip">dant mnf</FONT>] 没有`,      // Orizon
		`(GTalk) niconiconi: 无误`,                            // OrzGTalk
		`[FsckGoF] 群主女装吼不吼啊？`,                               // OrzIrc2P
		`[Lisa] \h: Lisa is here :-)`,                       // xmppbot
		`7[Miyamizu_Mitsuha] 厉害了`,                           // teleboto?
		`4[Abel_Abel] 这算不算父进程？`,                             //??
		`13[Universebenzene] 缺poppler-data？`,                //??
		`6[KireinaHoro_] 15「Re farseerfc: wow 麗狼加油...」謝謝`,   //??
		`[tg2offtopic@irc] [FQEgg] loop`,
		`[erhandsoME[m]@irc] 。。。`, // riot.im
	}
	for _, m := range msgs {
		nu, nm, color := ExtractRealUser("", m)
		log.Println("color:", color, "realuser:", nu, "realmsg:", nm)
		if true {
			nu, nm, color := ExtractRealUserMD("", m)
			log.Println("color:", color, "realuser:", nu, "realmsg:", nm)
		}
	}
}
