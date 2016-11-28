package purple

/*
#include <libpurple/purple.h>
*/
import "C"

type CoreUiOps struct {
	uiops *C.PurpleCoreUiOps
}

func newCoreUiOpsFrom(uiops *C.PurpleCoreUiOps) *CoreUiOps {
	this := &CoreUiOps{}
	this.uiops = uiops
	return this
}

func NewCoreUiOps() *CoreUiOps {
	this := &CoreUiOps{}
	return this
}

func CoreInit(ui string) bool {
	return c2goBool(C.purple_core_init(CCString(ui).Ptr))
}
func CoreQuit() {
	C.purple_core_quit()
}

func CoreGetVersion() string {
	ret := C.purple_core_get_version()
	return C.GoString(ret)
}
