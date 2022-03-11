package iter_test

import (
	"container/list"
	"fmt"
	"testing"

	. "github.com/disksing/iter/v2"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	a := make([]int, 100)
	b := SliceBegin(a)
	b.AllowMultiplePass()
	b1 := b.Next()
	rb := SliceRBegin(a)
	e := SliceEnd(a)
	assert.True(b.AdvanceN(100).Eq(e))
	assert.True(e.AdvanceN(-100).Eq(b))
	assert.Equal(100, Distance[int](b, e))
	assert.Equal(-100, Distance[int](e, b))

	re := SliceREnd(a)
	assert.True(rb.AdvanceN(100).Eq(re))
	assert.True(re.AdvanceN(-100).Eq(rb))
	assert.Equal(100, Distance[int](rb, re))
	assert.Equal(-100, Distance[int](re, rb))

	assert.True(b.Less(e))
	assert.True(rb.Less(re))
	assert.True(b.Less(b1))
	assert.False(b1.Less(b))

	e1 := e.Prev()
	assert.NotEqual(fmt.Sprintf("%s", e1), fmt.Sprintf("%s", rb))
}

func listEq[T comparable](assert *assert.Assertions, lst *list.List, v ...T) {
	end := SliceEnd(v)
	assert.True(Equal[T](ListBegin[T](lst), ListEnd[T](lst), SliceBegin(v), &end))
}

func TestListIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	lst := list.New()
	listEq[int](assert, lst)
	GenerateN(ListBackInserter[int](lst), 3, IotaGenerator(1))
	listEq(assert, lst, 1, 2, 3)

	b := ListBegin[int](lst)
	b.AllowMultiplePass()
	b1 := b.Next()
	assert.Equal(b, b1.Prev())

	rb := ListRBegin[int](lst)
	assert.Equal(rb.Read(), 3)
	rb1 := rb.Next()
	assert.Equal(rb1.Read(), 2)
	assert.True(_eq(rb.Next(), rb1))
	assert.True(_eq(rb1.Prev(), rb))
	assert.True(ListEnd[int](lst).Prev().Prev().Eq(b1))
	assert.Equal(ListREnd[int](lst).Prev().Read(), 1)

	assert.True(AdvanceN[int](b, 3).Eq(ListEnd[int](lst)))
	assert.True(AdvanceN[int](ListEnd[int](lst), -3).Eq(b))

	b.Write(2)
	listEq(assert, lst, 2, 2, 3)
	ListInserter[int](lst, lst.Back()).Write(4)
	listEq(assert, lst, 2, 2, 4, 3)
}

func TestStringIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	s := "abcdefg"
	assert.Equal(MakeString[byte](StringRBegin(s), StringREnd(s)), "gfedcba")

	StringBegin(s).AllowMultiplePass()

	assert.Contains(fmt.Sprint(StringBegin(s)), "->")
	assert.Contains(fmt.Sprint(StringRBegin(s)), "<-")

	assert.Equal(StringEnd(s).Prev().Read(), byte('g'))
	assert.Equal(StringREnd(s).Prev().Read(), byte('a'))

	assert.Equal(
		StringBegin(s).AdvanceN(3).Read(),
		byte('d'),
	)
	assert.Equal(
		StringRBegin(s).AdvanceN(3).Read(),
		byte('d'),
	)
	assert.Equal(
		StringEnd(s).AdvanceN(-3).Read(),
		byte('e'),
	)
	assert.Equal(
		StringREnd(s).AdvanceN(-3).Read(),
		byte('c'),
	)

	assert.Equal(StringBegin(s).Distance(StringEnd(s)), len(s))
	assert.Equal(StringEnd(s).Distance(StringBegin(s)), -len(s))
	assert.Equal(StringRBegin(s).Distance(StringREnd(s)), len(s))
	assert.Equal(StringREnd(s).Distance(StringRBegin(s)), -len(s))

	assert.True(StringBegin(s).Less(StringEnd(s)))
	assert.False(StringEnd(s).Less(StringBegin(s)))
	assert.True(StringRBegin(s).Less(StringREnd(s)))
	assert.False(StringREnd(s).Less(StringRBegin(s)))
}

func TestStringBuilder(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	var bs StringBuilderInserter[any]
	FillN[any](&bs, 3, 'a')
	assert.Equal(bs.String(), "aaa")

	bs.Reset()
	FillN[any](&bs, 3, rune('香'))
	assert.Equal(bs.String(), "香香香")

	bs.Reset()
	bs.Delimiter = ","
	FillN[any](&bs, 3, []byte("abc"))
	assert.Equal(bs.String(), "abc,abc,abc")

	bs.Reset()
	bs.Delimiter = "-->"
	FillN[any](&bs, 3, "abc")
	assert.Equal(bs.String(), "abc-->abc-->abc")

	var bsi StringBuilderInserter[int]
	bs.Delimiter = ""
	CopyN[int](IotaReader(1), 5, &bsi)
	assert.Equal(bsi.String(), "12345")
}

// func TestChanIterator(t *testing.T) {
// 	skipAfter(t, 1)

// 	assert := assert.New(t)
// 	ch := make(chan int)
// 	go func() {
// 		CopyN(IotaReader(1), 10, ChanWriter(ch))
// 		close(ch)
// 	}()
// 	assert.Equal(Accumulate(ChanReader(ch), ChanEOF, 0), 55)

// 	ch = make(chan int)
// 	close(ch)
// 	assert.True(ChanReader(ch).Eq(ChanEOF))

// 	ch = make(chan int)
// 	close(ch)
// 	assert.True(ChanEOF.Eq(ChanReader(ch)))

// 	assert.Nil(ChanEOF.Read())
// 	assert.Equal(ChanEOF.Next(), ChanEOF)
// 	assert.True(ChanEOF.Eq(ChanEOF))

// 	ch = make(chan int, 1)
// 	ch <- 100
// 	assert.Equal(ChanReader(ch).Read(), 100)

// 	ch = make(chan int, 2)
// 	ch <- 42
// 	ch <- 43
// 	it := ChanReader(ch)
// 	it.Next()
// 	assert.Equal(it.Read(), 43)

// 	ch = make(chan int, 10)
// 	CopyN(IotaReader(1), 5, ChanWriter(ch))
// 	close(ch)
// 	it = ChanReader(ch)
// 	it = AdvanceN(it, 3).(InputIter)
// 	assert.Equal(it.Read(), 4)
// 	assert.Equal(Distance(it, ChanEOF), 2)

// 	ch = make(chan int)
// 	w := ChanWriter(ch)
// 	assert.Panics(func() { AdvanceN(w, 1) })
// 	assert.Panics(func() { Distance(w, ChanEOF) })
// }

// type dummyWriter struct {
// 	tick int
// }

// func (w *dummyWriter) Write(b []byte) (int, error) {
// 	w.tick--
// 	if w.tick < 0 {
// 		return 0, errors.New("boom!")
// 	}
// 	return len(b), nil
// }

// func TestIOWriterPanics(t *testing.T) {
// 	assert := assert.New(t)
// 	assert.Panics(func() {
// 		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{}, ","))
// 	})
// 	assert.Panics(func() {
// 		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{tick: 1}, ","))
// 	})
// }
