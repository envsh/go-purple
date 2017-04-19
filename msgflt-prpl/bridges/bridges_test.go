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
	}
	for _, m := range msgs {
		nu, nm, color := ExtractRealUser("", m)
		log.Println(nu, "--", nm, "--", color)
	}
}
