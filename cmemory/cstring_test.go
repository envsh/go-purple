package cmemory

/*
 */
// import "C"
// use of cgo in test not supported

import (
	"fmt"
	"runtime"
	"testing"
)

// go test -v ./cmemory
func TestCString(t *testing.T) {
	cnt := 123
	freed := 0
	freeCounter = func() { freed += 1 }

	{
		for i := 0; i < cnt; i++ {
			s0 := fmt.Sprintf("abcdefgeeeeeeeeeeeeeeee.%d", i)
			s := CCString(s0)
			// t.Log(s)
			if s == nil {
				t.Error("can not nil")
			}
		}
	}
	if true {
		for i := 0; i < cnt; i++ {
			runtime.GC()
		}
	}
	if freed <= 0 {
		t.Error("not match:", freed, cnt)
	}
	// select {}
}
