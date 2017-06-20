package main

import (
	"log"
	"sync"

	python "github.com/sbinet/go-python"
)

// SnowNLP的golang 封装
// 完全相当于用python C API调用python函数
// 复杂度有点大，以后的情况能不用就不要用了
type SnowNLP struct {
	snlp   *python.PyObject // module
	clsty  *python.PyObject // class type
	clso   *python.PyObject // class object
	pynew1 sync.Once
	doc    string
}

func NewSnowNLP(doc string) *SnowNLP {
	this := &SnowNLP{doc: doc}
	this.pynew1.Do(func() { this.pynew() })
	return this
}

func (this *SnowNLP) pynew() {
	snlp := python.PyImport_ImportModule("snownlp") // why use 3 seconds
	// log.Println(snlp)
	clsty := snlp.GetAttrString("SnowNLP")
	// log.Println(clsty, python.PyString_AS_STRING(clsty.Repr()))
	clso := clsty.CallMethod("__new__", clsty)
	// log.Println(clso)
	// log.Println(clso, python.PyErr_Occurred())
	doc := python.PyString_FromString(this.doc)
	doc = doc.CallMethod("decode", "utf8")
	_ = clso.CallMethod("__init__", doc)
	// log.Println(r)
	// log.Println(python.PyString_AS_STRING(r.Repr()))

	this.snlp, this.clsty, this.clso = snlp, clsty, clso
}

func (this *SnowNLP) Words() (rets []string) {
	r := this.clso.GetAttrString("words")
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		// log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

func (this *SnowNLP) Sentences() (rets []string) {
	r := this.clso.GetAttrString("words")
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		// log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

func (this *SnowNLP) Han() string {
	r := this.clso.GetAttrString("han")
	return python.PyString_AsString(r.CallMethod("encode", "utf8"))
}

func (this *SnowNLP) PinYin() (rets []string) {
	r := this.clso.GetAttrString("pinyin")
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

func (this *SnowNLP) Sentiments() float64 {
	r := this.clso.GetAttrString("sentiments")
	log.Println(r, pyvartype(r), python.PyString_AsString(r.Repr()))
	return python.PyFloat_AsDouble(r)
}

// TODO [(u'\u8fd9\u4e2a', u'r'), (u'\u4e1c\u897f', u'n'), (u'\u771f\u5fc3', u'd'), (u'\u5f88', u'd'), (u'\u8d5e', u'Vg')]
func (this *SnowNLP) Tags() (rets []string) {
	r := this.clso.GetAttrString("tags")
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

func (this *SnowNLP) Tf() (rets []map[string]int) {
	r := this.clso.GetAttrString("tf")
	// log.Println(r, pyvartype(r), python.PyString_AsString(r.Repr()))
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		keys := python.PyDict_Keys(item)
		key := python.PyList_GetItem(keys, 0)
		value := python.PyDict_GetItem(item, key)
		key = key.CallMethod("encode", "utf8")
		rets = append(rets, map[string]int{
			python.PyString_AsString(key): python.PyInt_AsLong(value),
		})
	}
	return
}

func (this *SnowNLP) Idf() (rets map[string]float64) {
	rets = make(map[string]float64)
	r := this.clso.GetAttrString("idf")
	// log.Println(r, pyvartype(r), python.PyString_AsString(r.Repr()))
	python.PyDict_Size(r)
	keys := python.PyDict_Keys(r)

	for i := 0; i < python.PyList_Size(keys); i++ {
		item := python.PyList_GetItem(keys, i)
		value := python.PyDict_GetItem(r, item)
		item = item.CallMethod("encode", "utf8")
		rets[python.PyString_AsString(item)] = python.PyFloat_AsDouble(value)
	}
	return
}

// method
func (this *SnowNLP) Sim(s string) (rets []float64) {
	r := this.clso.CallMethod("sim", s)
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		// log.Println(item, pyvartype(item))
		switch pyvartype(item) {
		case PY_TYPE_INT:
			rets = append(rets, float64(python.PyInt_AsLong(item)))
		case PY_TYPE_FLOAT:
			rets = append(rets, python.PyFloat_AsDouble(item))
		default:
			log.Fatalln(pyvartype(item))
		}
	}
	return
}

// method
func (this *SnowNLP) Summary(n int) (rets []string) {
	r := this.clso.CallMethod("summary", n)
	// log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

// method
func (this *SnowNLP) Keywords(n int) (rets []string) {
	r := this.clso.CallMethod("keywords", n) // panic: runtime error: cgo argument has Go pointer to Go pointer
	log.Println(r, python.PyList_Size(r))
	pyvar2go(r)
	for i := 0; i < python.PyList_Size(r); i++ {
		item := python.PyList_GetItem(r, i)
		item = item.CallMethod("encode", "utf8")
		log.Println(item)
		rets = append(rets, python.PyString_AsString(item))
	}
	return
}

// support list, int, string
func pyvar2go(v *python.PyObject) interface{} {
	switch pyvartype(v) {
	case PY_TYPE_LIST:
		// 想起来了，就怕没法转啊，list的元素可能还是其他复杂的数据结构
	}
	return nil
}

func pyvartype(v *python.PyObject) int {
	vty := v.Type()
	if t, ok := pytypesm[vty]; ok {
		return t
	}
	for to, t := range pytypesm {
		if vty.Compare(to) == 0 {
			return t
		}
	}
	log.Println(vty, pytypesm, python.PyString_AsString(vty.Repr()))
	return PY_TYPE_UNKNOWN
}

const (
	PY_TYPE_UNKNOWN = iota
	PY_TYPE_LIST
	PY_TYPE_DICT
	PY_TYPE_INT
	PY_TYPE_STR
	PY_TYPE_FLOAT
)

var pytypesm = make(map[*python.PyObject]int)

func init() {
	err := python.Initialize()
	if err != nil {
		panic(err.Error())
	}
	// log.Println("1")

	if true {
		pytypesm[python.PyList_New(0).Type()] = PY_TYPE_LIST
		pytypesm[python.PyDict_New().Type()] = PY_TYPE_DICT
		pytypesm[python.PyInt_FromLong(0).Type()] = PY_TYPE_INT
		pytypesm[python.PyString_FromString("").Type()] = PY_TYPE_STR
		pytypesm[python.PyFloat_FromDouble(0.1).Type()] = PY_TYPE_FLOAT
		pytypesm[python.PyFloat_FromString(python.PyString_FromString("0.123")).Type()] = PY_TYPE_FLOAT
	}

	if false {
		nlp := NewSnowNLP("这个东西真心很赞")
		log.Println(nlp.Words())
		log.Println(nlp.Sentiments())
		log.Println(nlp.Sentences())
		log.Println(nlp.Han())
		// log.Println(nlp.Tags())
		log.Println(nlp.Sim("真"))
		// log.Println(nlp.Keywords(3))
		log.Println(nlp.Tf())
		log.Println(nlp.Idf())
	}

	if false {
		snlp := python.PyImport_ImportModule("snownlp") // why use 3 seconds
		log.Println(snlp)
		clsa := snlp.GetAttrString("SnowNLP")
		log.Println(clsa, python.PyString_AS_STRING(clsa.Repr()))
		clso := clsa.CallMethod("__new__", clsa)
		log.Println(clso)
		log.Println(clso, python.PyErr_Occurred())
		r := clso.CallMethod("__init__", python.PyString_FromString("呵呵"))
		log.Println(r)
		log.Println(python.PyString_AS_STRING(r.Repr()))
		r = clso.GetAttrString("words")
		log.Println(r)
		log.Println(python.PyString_AS_STRING(r.Type().Str()))
		log.Println(python.PyString_AS_STRING(r.Bytes()))
	}
}
