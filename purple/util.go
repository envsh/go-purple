package purple

/*
#include <libpurple/purple.h>
*/
import "C"

func HomeDir() string {
	return C.GoString((*C.char)(C.purple_home_dir()))
}
func UserDir() string {
	return C.GoString(C.purple_user_dir())
}
func UtilSetUserDir(dir string) {
	C.purple_util_set_user_dir(CCString(dir).Ptr)
}

/*
func BuildDir() string {
	return C.GoString(C.purple_build_dir(path, mode))
}
*/
