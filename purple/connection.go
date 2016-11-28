package purple

/*
#include <libpurple/purple.h>
*/
import "C"
import "unsafe"

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

func newConnectionFrom(conn *C.PurpleConnection) *Connection {
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

func (this *Connection) ConnGetAccount() *Account {
	ac := C.purple_connection_get_account(this.conn)
	return newAccountFrom(ac)
}

func (this *Connection) ConnSetDisplayName(name string) {
	C.purple_connection_set_display_name(this.conn, CCString(name).Ptr)
}

func (this *Connection) ConnFindChat(id int) *Conversation {
	conv := C.purple_find_chat(this.conn, C.int(id))
	return newConversationFrom(conv)
}

func (this *Connection) GetPrpl() *Plugin {
	plugin := C.purple_connection_get_prpl(this.conn)
	return newPluginFrom((*C.PurplePlugin)(plugin))
}
func (this *Connection) GetPrplInfo() *PluginInfo {
	plugin := C.purple_connection_get_prpl(this.conn)
	return newPluginInfoFrom((*C.PurplePluginInfo)(plugin.info))
}

// error info
type ConnectionErrorInfo struct {
	// private
	cei *C.PurpleConnectionErrorInfo
}

func newConnectionErrorInfoFrom(cei *C.PurpleConnectionErrorInfo) *ConnectionErrorInfo {
	this := &ConnectionErrorInfo{}
	this.cei = cei
	return this
}

func (this *ConnectionErrorInfo) Get() (int, string) {
	return 0, C.GoString(this.cei.description)
}

func ConnectionsGetHandle() unsafe.Pointer {
	return C.purple_connections_get_handle()
}
