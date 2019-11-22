package iter_test

import (
	"math/rand"
	"time"

	"github.com/disksing/iter"
)

func randInt() int {
	return rand.Intn(100)
}

func randIntSlice() []int {
	l, h := randInt(), randInt()
	if l > h {
		l, h = h, l
	}
	s := make([]int, randInt())
	for i := range s {
		s[i] = l + rand.Intn(h-l+1)
	}
	return s
}

var (
	begin = iter.SliceBegin
	end   = iter.SliceEnd
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
