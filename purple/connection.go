package purple

/*
#include <libpurple/purple.h>
*/
import "C"

import (
	"log"
)

const (
	DISCONNECTED = int(C.PURPLE_DISCONNECTED)
	CONNECTED    = int(C.PURPLE_CONNECTED)
	CONNECTING   = int(C.PURPLE_CONNECTING)

	MESSAGE_SEND   = int(C.PURPLE_MESSAGE_SEND)
	MESSAGE_RECV   = int(C.PURPLE_MESSAGE_RECV)
	MESSAGE_SYSTEM = int(C.PURPLE_MESSAGE_SYSTEM)
)

type Connection struct {
	conn *C.PurpleConnection
}

func newConnectWrapper(conn *C.PurpleConnection) *Connection {
	this := &Connection{conn}
	return this
}

func (this *Connection) ConnSetState(state int) {
	switch state {
	case DISCONNECTED:
		C.purple_connection_set_state(this.conn, C.PURPLE_DISCONNECTED)
	case CONNECTED:
		C.purple_connection_set_state(this.conn, C.PURPLE_CONNECTED)
	case CONNECTING:
		C.purple_connection_set_state(this.conn, C.PURPLE_CONNECTING)
	default:
		log.Panicln("not supported", state)
	}
}

func (this *Connection) ConnGetState() int {
	state := C.purple_connection_get_state(this.conn)
	return int(state)
}

func (this *Connection) ConnSetAccount(ac *Account) {
	C.purple_connection_set_account(this.conn, ac.account)
}

func (this *Connection) ConnSetDisplayName(name string) {
	C.purple_connection_set_display_name(this.conn, C.CString(name))
}
