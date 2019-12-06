package iter_test

import (
	"container/list"
	"sync/atomic"
	"testing"

	. "github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

// check if need skip the test. A test will be run count/count+1 times.
func skipAfter(t *testing.T, count int) {
	if atomic.LoadInt32(&testCounter) > int32(count) {
		t.SkipNow()
	}
}

var testCounter int32

func TestTouchSkipCounter(t *testing.T) {
	atomic.AddInt32(&testCounter, 1)
}

func TestMisc(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	it := IotaReader(100)
	assert.Equal(it.Read(), 100)
	assert.False(it.Eq(IotaReader(100)))
	it = NextInputIter(it)
	assert.Equal(it.Read(), 101)

	it = RepeatReader(100)
	assert.Equal(it.Read(), 100)
	assert.False(it.Eq(RepeatReader(100)))
	it = NextInputIter(it)
	assert.Equal(it.Read(), 100)

	g := IotaGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 101)
	g = RepeatGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 100)
}

func TestMakeString(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	bs := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!'}
	assert.Equal(MakeString(begin(bs), end(bs)), "hello world!")

	rs := []rune{'改', '革', '春', '风', '吹', '满', '地'}
	assert.Equal(MakeString(begin(rs), end(rs)), "改革春风吹满地")
}

func TestErase(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	a := []int{1, 2, 3, 4, 5, 6}
	Erase(&a, AdvanceN(begin(a), 1), AdvanceN(begin(a), 3))
	sliceEqual(assert, a, []int{1, 4, 5, 6})
	Erase(&a, AdvanceN(begin(a), 2))
	sliceEqual(assert, a, []int{1, 4})
	Erase(&a)
	sliceEqual(assert, a, []int{})
	assert.Panics(func() { Erase(&a, begin(a), begin(a), begin(a)) })

	lst := list.New()
	GenerateN(ListBackInserter(lst), 6, IotaGenerator(1))
	Erase(lst, AdvanceN(ListBegin(lst), 1), AdvanceN(ListBegin(lst), 3))
	listEq(assert, lst, 1, 4, 5, 6)
	Erase(lst, AdvanceN(ListBegin(lst), 2))
	listEq(assert, lst, 1, 4)
	Erase(lst)
	listEq(assert, lst)
	assert.Panics(func() { Erase(lst, begin(a)) })
	assert.Panics(func() { Erase(lst, begin(a), end(a)) })

	arr := [2]int{1, 2}
	assert.Panics(func() { Erase(&arr) })
}
