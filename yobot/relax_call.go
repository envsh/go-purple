package main

import (
	"log"
	"reflect"
	"runtime"
	"sync"
)

type Any interface{}

// sync and ordered call object
// 该实例的call实现同步的顺序的调用，但在执行在不同的routine中
// 类似于进程内的rpc，减少锁的显式调用用（从代码编写的角度）。
type RelaxCallObject struct {
	__chkchmu  sync.Mutex
	__chcall   chan interface{}
	__chreturn chan interface{}
	__chstop   chan struct{}
}

func (this *RelaxCallObject) startRoutine() {
	this.__chkchmu.Lock()
	defer this.__chkchmu.Unlock()

	if this.__chcall != nil {
		return
	}

	this.__chcall = make(chan interface{}, 123)
	this.__chreturn = make(chan interface{}, 123)
	this.__chstop = make(chan struct{}, 0)

	go func(chcall chan interface{}, chreturn chan interface{}, chstop chan struct{}) {
		log.Printf("goroutine started:\n")
		for {
			select {
			case fi := <-chcall:
				fi.(func())()
				// fn := fi.(func() []reflect.Value)
				// res := fn()
				// chreturn <- res
			case <-chstop:
				goto endfor
			}
		}
	endfor:
		log.Printf("goroutine stopped: \n")
	}(this.__chcall, this.__chreturn, this.__chstop)

	log.Printf("install final...%p\n", this)
	runtime.SetFinalizer(this, func(this_ *RelaxCallObject) {
		log.Printf("final...%p\n", this_)
		this_.stopRoutine()
		log.Printf("final done...%p\n", this_)
	})
}

func (this *RelaxCallObject) stopRoutine() { this.__chstop <- struct{}{} }

// x.Call0(func(){xxx})
func (this *RelaxCallObject) Call0(f func()) []Any { return this.callImpl(f, nil) }

// x.Call1(func()Any{return xxx})
func (this *RelaxCallObject) Call1(f func() Any) []Any { return this.callImpl(f, nil) }

// x.Call2(func()(Any,Any){return xxx})
func (this *RelaxCallObject) Call2(f func() (Any, Any)) []Any { return this.callImpl(f, nil) }

// x.Call(funcname, 1,2,3)
func (this *RelaxCallObject) Call(f interface{}, args ...interface{}) []Any {
	return this.callImpl(f, args)
}

func (this *RelaxCallObject) callImpl(f interface{}, args []interface{}) []Any {
	// 双层变量检测，还能速度快一点，整体提高30%
	if this.__chcall == nil {
		this.startRoutine()
	}

	fv := reflect.ValueOf(f)
	_chreturn := make(chan []reflect.Value, 1)
	this.__chcall <- func() { _chreturn <- this.invoke(f, args) }
	// outValuesIf := <-this.__chreturn // 共用返回值chan可能导致时序问题，如一个调用抢到另一个调用的返回值。
	// outValues := outValuesIf.([]reflect.Value)
	outValues := <-_chreturn
	outAnys := this.value2any(outValues)
	if fv.Type().NumOut() != len(outAnys) {
		log.Println("wtf", fv.Type().NumOut(), len(outAnys), outAnys)
	}
	return outAnys
}

func (this *RelaxCallObject) invoke(f interface{}, args []interface{}) []reflect.Value {
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		log.Panicln(fv.Kind().String())
	}
	in := make([]reflect.Value, len(args))
	for idx := 0; idx < len(args); idx++ {
		in[idx] = reflect.ValueOf(args[idx])
	}
	out := fv.Call(in)
	if fv.Type().NumOut() != len(out) {
		log.Println("wtf", fv.Type().NumOut(), len(out), out)
	}
	return out
}

func (this *RelaxCallObject) value2any(outValues []reflect.Value) []Any {
	outAnys := make([]Any, len(outValues))
	for idx := 0; idx < len(outValues); idx++ {
		outAnys[idx] = outValues[idx].Interface()
	}
	return outAnys
}
