package iter_test

import (
	"sync/atomic"
	"testing"
)

// check if need skip the test. A test will be run count/count+1 times.
func skipAfter(t *testing.T, count int) {
	if atomic.LoadInt32(&testCounter) > int32(count) {
		t.SkipNow()
	}
}

var testCounter int32

func TestTouchSkip(t *testing.T) {
	atomic.AddInt32(&testCounter, 1)
}
