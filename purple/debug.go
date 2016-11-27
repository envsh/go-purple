package purple

/*
#include <libpurple/purple.h>

static void gopurple_debug(PurpleDebugLevel level, const char *category, const char *format)
{ purple_debug(level, category, format); }
static void gopurple_debug_misc(const char *category, const char *format)
{ purple_debug_misc(category, format); }
static void gopurple_debug_info(const char *category, const char *format)
{ purple_debug_info(category, format); }
static void gopurple_debug_warning(const char *category, const char *format)
{ purple_debug_warning(category, format); }
static void gopurple_debug_error(const char *category, const char *format)
{ purple_debug_error(category, format); }
static void gopurple_debug_fatal(const char *category, const char *format)
{ purple_debug_fatal(category, format); }

*/
import "C"

func Debug(level int, cat string, log string) {
	C.gopurple_debug(C.PurpleDebugLevel(level), CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugMisc(cat string, log string) {
	C.gopurple_debug_misc(CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugInfo(cat string, log string) {
	C.gopurple_debug_info(CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugWarning(cat string, log string) {
	C.gopurple_debug_warning(CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugError(cat string, log string) {
	C.gopurple_debug_error(CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugFatal(cat string, log string) {
	C.gopurple_debug_fatal(CCString(cat).Ptr, CCString(log).Ptr)
}

func DebugSetEnabled(enabled bool) {
	if enabled {
		C.purple_debug_set_enabled(C.TRUE)
	} else {
		C.purple_debug_set_enabled(C.FALSE)
	}
}

func DebugIsEnabled() bool {
	bret := C.purple_debug_is_enabled()
	if bret == C.TRUE {
		return true
	} else {
		return false
	}
}

func DebugSetVerbose(verbose bool) {
	if verbose {
		C.purple_debug_set_verbose(C.TRUE)
	} else {
		C.purple_debug_set_verbose(C.FALSE)
	}
}

func DebugIsVerbose() bool {
	bret := C.purple_debug_is_verbose()
	if bret == C.TRUE {
		return true
	} else {
		return false
	}
}

func DebugSetUnsafe(unsafe bool) {
	if unsafe {
		C.purple_debug_set_unsafe(C.TRUE)
	} else {
		C.purple_debug_set_unsafe(C.FALSE)
	}
}

func DebugIsUnsafe() bool {
	bret := C.purple_debug_is_unsafe()
	if bret == C.TRUE {
		return true
	} else {
		return false
	}
}
