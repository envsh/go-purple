package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type Connection struct {
	conn *C.PurpleConnection
}

func newConnectWrapper(conn *C.PurpleConnection) *Connection {
	this := &Connection{conn}
	return this
}
