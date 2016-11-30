/*
wrapper purple dbus module
*/
package purple

/*
#define HAVE_DBUS 1
#include <libpurple/purple.h>
#include <libpurple/dbus-purple.h>
*/
import "C"

import (
	"strings"
)

const (
	DBUS_SERVICE_PURPLE   = C.DBUS_SERVICE_PURPLE   //      "im.pidgin.purple.PurpleService"
	DBUS_PATH_PURPLE      = C.DBUS_PATH_PURPLE      //         "/im/pidgin/purple/PurpleObject"
	DBUS_INTERFACE_PURPLE = C.DBUS_INTERFACE_PURPLE //    "im.pidgin.purple.PurpleInterface"
)

func GetDBusService() string {
	return strings.Replace(DBUS_SERVICE_PURPLE, ".purple.", ".gopurple.", -1)
}

func GetDBusPath() string {
	return strings.Replace(DBUS_PATH_PURPLE, "/purple/", "/gopurple/", -1)
}

func GetDBusInterface() string {
	return strings.Replace(DBUS_INTERFACE_PURPLE, ".purple.", ".gopurple.", -1)
}
