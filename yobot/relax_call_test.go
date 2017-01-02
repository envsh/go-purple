package main

import (
	"errors"
	"log"
	"runtime"
	"testing"
	"time"
)

type TestObject struct {
	RelaxCallObject
	a string
	b [1234567]int
}

// test nil error return?
func TestRc3(t *testing.T) {
	to := &TestObject{}
	f1 := func() (int, error) { return 0, nil }
	f2 := func() (int, error) { return 0, errors.New("aaaa") }
	rets1 := to.Call2(func() (Any, Any) { return f1() })
	log.Println(len(rets1), rets1)
	rets2 := to.Call2(func() (Any, Any) { return f2() })
	log.Println(len(rets2), rets2)
}

func TestRc2(t *testing.T) {
	log.Println("goroutine num:", runtime.NumGoroutine())
	gorn := runtime.NumGoroutine()

	if true {
		to := &TestObject{}

		fraw := func() (int, string) { return 3, "345" }
		f := func() (Any, Any) { return fraw() }

		if to != nil {
			log.Println("calling...")
			r := to.Call2(f)
			log.Println(r)
			if gorn+1 != runtime.NumGoroutine() {
				t.Error(gorn, runtime.NumGoroutine())
			}

		}

		to = nil
	}

	for i := 0; i < 3; i++ {
		log.Println("gc...")
		runtime.GC()
		log.Println("sleep...", runtime.NumGoroutine())
		if runtime.NumGoroutine() == gorn {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if gorn != runtime.NumGoroutine() {
		t.Error(gorn, runtime.NumGoroutine())
	}
}

func TestRc1(t *testing.T) {
	log.Println("goroutine num:", runtime.NumGoroutine())
	gorn := runtime.NumGoroutine()

	if true {
		to := &TestObject{}
		// runtime.SetFinalizer(to, func(o *TestObject) { o.stopRoutine() })

		f := func(a int) int { return 1 }
		if to != nil {
			log.Println("calling...")
			r := to.Call(f, 123)
			log.Println(r)
			if gorn+1 != runtime.NumGoroutine() {
				t.Error(gorn, runtime.NumGoroutine())
			}

		}

		n := 5000000
		if true {
			btime := time.Now()
			for i := 0; i < n; i++ {
				to.Call(f, i)
			}
			etime := time.Now()
			dtime := etime.Sub(btime)
			log.Println(dtime, dtime/time.Duration(n))
		}

		if true {
			btime := time.Now()
			for i := 0; i < n; i++ {
				f(i)
			}
			etime := time.Now()
			dtime := etime.Sub(btime)
			log.Println(dtime, dtime/time.Duration(n))
		}
		// 这两种调用大概差300*2.5=750倍，更有可能达到1000倍，果然好慢啊

		to = nil
	}

	for i := 0; i < 3; i++ {
		log.Println("gc...")
		runtime.GC()
		log.Println("sleep...", runtime.NumGoroutine())
		if runtime.NumGoroutine() == gorn {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if gorn != runtime.NumGoroutine() {
		t.Error(gorn, runtime.NumGoroutine())
	}
}
