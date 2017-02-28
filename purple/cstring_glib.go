package purple

/*
#include <stdlib.h>
#include <glib.h>
*/
import "C"

func (this *CString) gpointer() C.gpointer {
	return (C.gpointer)(this.Ptr)
}

func (this *CString) gconstpointer() C.gconstpointer {
	return (C.gconstpointer)(this.Ptr)
}
