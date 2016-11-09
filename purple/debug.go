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
	C.gopurple_debug(C.PurpleDebugLevel(level), C.CString(cat), C.CString(log))
}

func DebugMisc(cat string, log string) {
	C.gopurple_debug_misc(C.CString(cat), C.CString(log))
}

func DebugInfo(cat string, log string) {
	C.gopurple_debug_info(C.CString(cat), C.CString(log))
}

func DebugWarning(cat string, log string) {
	C.gopurple_debug_warning(C.CString(cat), C.CString(log))
}

func DebugError(cat string, log string) {
	C.gopurple_debug_error(C.CString(cat), C.CString(log))
}

func DebugFatal(cat string, log string) {
	C.gopurple_debug_fatal(C.CString(cat), C.CString(log))
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
