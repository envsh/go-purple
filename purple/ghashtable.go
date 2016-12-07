// only support string key/value GHashTable
package purple

/*
#include <glib.h>

static GHashTable *go_g_hash_table_new_full() {
    return g_hash_table_new_full(g_str_hash, g_str_equal, g_free, g_free);
}
*/
import "C"
import _ "log"

type GHashTable struct {
	ht *C.GHashTable
}

func newGHashTableFrom(ht *C.GHashTable) *GHashTable {
	this := &GHashTable{}
	this.ht = ht
	return this
}

// TODO when to free
func NewGHashTable() *GHashTable {
	this := &GHashTable{}
	this.ht = C.go_g_hash_table_new_full()
	return this
}

func NewGHashTableFromMap(d map[string]string) *GHashTable {
	this := &GHashTable{}
	this.ht = C.go_g_hash_table_new_full()

	for k, v := range d {
		this.Insert(k, v)
	}

	return this
}

func (this *GHashTable) Destroy() {
	C.g_hash_table_destroy(this.ht)
}

func (this *GHashTable) Lookup(key string) string {
	val := C.g_hash_table_lookup(this.ht, CCString(key).Ptr)
	return C.GoString((*C.char)(val))
}

func (this *GHashTable) Insert(key string, value string) bool {
	// hash table will manange it, don't use CCString
	ret := C.g_hash_table_insert(this.ht, C.CString(key), C.CString(value))
	if ret == C.TRUE {
		return true
	} else {
		return false
	}
}

func (this *GHashTable) Replace(key string, value string) bool {
	// hash table will manange it, don't use CCString
	ret := C.g_hash_table_replace(this.ht, C.CString(key), C.CString(value))
	if ret == C.TRUE {
		return true
	} else {
		return false
	}
}

func (this *GHashTable) Add(key string) bool {
	ret := C.g_hash_table_add(this.ht, C.CString(key))
	if ret == C.TRUE {
		return true
	} else {
		return false
	}
}

func (this *GHashTable) Remove(key string) bool {
	ret := C.g_hash_table_remove(this.ht, CCString(key).Ptr)
	if ret == C.TRUE {
		return true
	} else {
		return false
	}
}

func (this *GHashTable) Contains(key string) bool {
	ret := C.g_hash_table_contains(this.ht, CCString(key).Ptr)
	if ret == C.TRUE {
		return true
	} else {
		return false
	}
}

func (this *GHashTable) Size() uint {
	ret := C.g_hash_table_size(this.ht)
	return uint(ret)
}

// need confirm type: string/string
func (this *GHashTable) GetKeys() []string {
	lst := C.g_hash_table_get_keys(this.ht)
	if lst == nil {
	} else {
		res := make([]string, 0)
		len := C.g_list_length(lst)
		for idx := 0; idx < int(len); idx++ {
			data := C.g_list_nth_data(lst, C.guint(idx))
			res = append(res, C.GoString((*C.char)(data)))
		}
		return res
	}

	return nil
}

func (this *GHashTable) GetValues() []string {
	lst := C.g_hash_table_get_values(this.ht)
	if lst == nil {
	} else {
		res := make([]string, 0)
		len := C.g_list_length(lst)
		for idx := 0; idx < int(len); idx++ {
			data := C.g_list_nth_data(lst, C.guint(idx))
			res = append(res, C.GoString((*C.char)(data)))
		}
		return res
	}
	return nil
}

func (this *GHashTable) ToMap() map[string]string {
	res := make(map[string]string, 0)
	lst := C.g_hash_table_get_keys(this.ht)
	if lst == nil {
	} else {
		len := C.g_list_length(lst)
		for idx := 0; idx < int(len); idx++ {
			key := C.g_list_nth_data(lst, C.guint(idx))
			val := C.g_hash_table_lookup(this.ht, key)
			res[C.GoString((*C.char)(key))] = C.GoString((*C.char)(val))
		}
	}
	return res
}

// TODO 参数包含C数据类型，应该改为内部使用
func (this *GHashTable) Each(functor func(C.gpointer, C.gpointer)) {
	this.Map(func(k, v C.gpointer) (interface{}, interface{}) {
		functor(k, v)
		return nil, nil
	})
}

func (this *GHashTable) Map(
	functor func(C.gpointer, C.gpointer) (interface{}, interface{})) map[interface{}]interface{} {
	res := make(map[interface{}]interface{}, 0)
	lst := C.g_hash_table_get_keys(this.ht)
	if lst != nil {
		len := C.g_list_length(lst)
		for idx := 0; idx < int(len); idx++ {
			key := C.g_list_nth_data(lst, C.guint(idx))
			val := C.g_hash_table_lookup(this.ht, key)
			gokey, goval := functor(key, val)
			res[gokey] = goval
		}
	}
	return res
}

//////////
type GList struct {
	lst *C.GList
}

func newGListFrom(lst *C.GList) *GList {
	this := &GList{}
	this.lst = lst
	return this
}

func (this *GList) ToStringArray() []string {
	res := make([]string, 0)

	len := C.g_list_length(this.lst)
	for idx := 0; idx < int(len); idx++ {
		item := C.g_list_nth_data(this.lst, C.guint(idx))
		str := C.GoString((*C.char)(item))
		res = append(res, str)
	}

	return res
}

func (this *GList) Each(functor func(C.gpointer)) {
	this.Map(func(item C.gpointer) interface{} {
		functor(item)
		return nil
	})
}

func (this *GList) Map(functor func(C.gpointer) interface{}) []interface{} {
	res := make([]interface{}, 0)
	len := C.g_list_length(this.lst)
	for idx := 0; idx < int(len); idx++ {
		item := C.g_list_nth_data(this.lst, C.guint(idx))
		goitem := functor(item)
		res = append(res, goitem)
	}
	return res
}

//////////
type GSList struct {
	lst *C.GSList
}

func newGSListFrom(lst *C.GSList) *GSList {
	this := &GSList{}
	this.lst = lst
	return this
}

func (this *GSList) ToStringArray() []string {
	res := make([]string, 0)

	len := C.g_slist_length(this.lst)
	for idx := 0; idx < int(len); idx++ {
		item := C.g_slist_nth_data(this.lst, C.guint(idx))
		str := C.GoString((*C.char)(item))
		res = append(res, str)
	}

	return res
}

func (this *GSList) Each(functor func(C.gpointer)) {
	this.Map(func(item C.gpointer) interface{} {
		functor(item)
		return nil
	})
}

func (this *GSList) Map(functor func(C.gpointer) interface{}) []interface{} {
	res := make([]interface{}, 0)
	len := C.g_slist_length(this.lst)
	for idx := 0; idx < int(len); idx++ {
		item := C.g_slist_nth_data(this.lst, C.guint(idx))
		goitem := functor(item)
		res = append(res, goitem)
	}
	return res
}
