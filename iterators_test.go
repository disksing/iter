package iter_test

import (
	"container/list"
	"errors"
	"fmt"
	"testing"

	. "github.com/disksing/iter/v2"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	a := make([]int, 100)
	b := begin(a)
	b.AllowMultiplePass()
	assert.Implements((*ReadWriter)(nil), b)
	assert.Implements((*InputIter)(nil), b)
	assert.Implements((*OutputIter)(nil), b)
	assert.Implements((*ForwardReadWriter)(nil), b)
	assert.Implements((*BidiReadWriter)(nil), b)
	assert.Implements((*RandomReadWriter)(nil), b)

	b1 := NextRandomReadWriter(b)
	assert.True(NextBidiIter(b).Eq(b1))
	assert.True(PrevBidiIter(b1).Eq(b))
	assert.True(NextBidiReader(b).Eq(b1))
	assert.True(PrevBidiReader(b1).Eq(b))
	assert.True(NextBidiWriter(b).Eq(b1))
	assert.True(PrevBidiWriter(b1).Eq(b))
	assert.True(NextBidiReadWriter(b).Eq(b1))
	assert.True(PrevBidiReadWriter(b1).Eq(b))
	assert.True(NextRandomIter(b).Eq(b1))
	assert.True(PrevRandomIter(b1).Eq(b))
	assert.True(NextRandomReader(b).Eq(b1))
	assert.True(PrevRandomReader(b1).Eq(b))
	assert.True(NextRandomWriter(b).Eq(b1))
	assert.True(PrevRandomWriter(b1).Eq(b))
	assert.True(NextRandomReadWriter(b).Eq(b1))
	assert.True(PrevRandomReadWriter(b1).Eq(b))

	rb := SliceRBegin(a)
	rb1 := NextRandomReadWriter(SliceRBegin(a))
	assert.True(NextBidiIter(rb).Eq(rb1))
	assert.True(PrevBidiIter(rb1).Eq(rb))
	assert.True(NextBidiReader(rb).Eq(rb1))
	assert.True(PrevBidiReader(rb1).Eq(rb))
	assert.True(NextBidiWriter(rb).Eq(rb1))
	assert.True(PrevBidiWriter(rb1).Eq(rb))
	assert.True(NextBidiReadWriter(rb).Eq(rb1))
	assert.True(PrevBidiReadWriter(rb1).Eq(rb))
	assert.True(NextRandomIter(rb).Eq(rb1))
	assert.True(PrevRandomIter(rb1).Eq(rb))
	assert.True(NextRandomReader(rb).Eq(rb1))
	assert.True(PrevRandomReader(rb1).Eq(rb))
	assert.True(NextRandomWriter(rb).Eq(rb1))
	assert.True(PrevRandomWriter(rb1).Eq(rb))
	assert.True(NextRandomReadWriter(rb).Eq(rb1))
	assert.True(PrevRandomReadWriter(rb1).Eq(rb))

	e := end(a)
	assert.True(AdvanceN(b, 100).(ForwardIter).Eq(e))
	assert.True(AdvanceN(e, -100).(ForwardIter).Eq(b))
	assert.Equal(Distance(b, e), 100)
	assert.Equal(Distance(e, b), -100)

	re := SliceREnd(a)
	assert.True(AdvanceN(rb, 100).(ForwardIter).Eq(re))
	assert.True(AdvanceN(re, -100).(ForwardIter).Eq(rb))
	assert.Equal(Distance(rb, re), 100)
	assert.Equal(Distance(re, rb), -100)

	assert.True(b.Less(e))
	assert.True(rb.Less(re))
	assert.True(b.Less(b1))
	assert.False(b1.Less(b))

	e1 := e.Prev().(RandomReadWriter)
	assert.NotEqual(fmt.Sprintf("%s", e1), fmt.Sprintf("%s", rb))
}

func listEq(assert *assert.Assertions, lst *list.List, v ...Any) {
	assert.True(Equal(lBegin(lst), lEnd(lst), begin(v), end(v)))
}

func TestListIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	lst := list.New()
	listEq(assert, lst)
	GenerateN(ListBackInserter(lst), 3, IotaGenerator(1))
	listEq(assert, lst, 1, 2, 3)

	b := lBegin(lst)
	assert.Implements((*ReadWriter)(nil), b)
	assert.Implements((*InputIter)(nil), b)
	assert.Implements((*OutputIter)(nil), b)
	assert.Implements((*ForwardReadWriter)(nil), b)
	assert.Implements((*BidiReadWriter)(nil), b)
	_, ok := b.(RandomIter)
	assert.False(ok)
	b.AllowMultiplePass()
	b1 := NextBidiReadWriter(b)
	assert.True(NextBidiIter(b).Eq(b1))
	assert.True(PrevBidiIter(b1).Eq(b))
	assert.True(NextBidiReader(b).Eq(b1))
	assert.True(PrevBidiReader(b1).Eq(b))
	assert.True(NextBidiWriter(b).Eq(b1))
	assert.True(PrevBidiWriter(b1).Eq(b))
	assert.True(NextBidiReadWriter(b).Eq(b1))
	assert.True(PrevBidiReadWriter(b1).Eq(b))

	rb := ListRBegin(lst)
	assert.Equal(rb.Read(), 3)
	rb1 := NextBidiReadWriter(rb)
	assert.Equal(rb1.Read(), 2)
	assert.True(_eq(rb.Next(), rb1))
	assert.True(_eq(rb1.Prev(), rb))
	assert.True(lEnd(lst).Prev().Prev().Eq(b1))
	assert.Equal(ListREnd(lst).Prev().(Reader).Read(), 1)

	assert.True(AdvanceN(b, 3).(BidiReader).Eq(lEnd(lst)))
	assert.True(AdvanceN(lEnd(lst), -3).(BidiReader).Eq(b))

	b.Write(2)
	listEq(assert, lst, 2, 2, 3)
	ListInserter(lst, lst.Back()).Write(4)
	listEq(assert, lst, 2, 2, 4, 3)
}

func TestStringIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	s := "abcdefg"
	assert.Equal(MakeString(StringRBegin(s), StringREnd(s)), "gfedcba")

	sBegin(s).AllowMultiplePass()

	assert.Contains(fmt.Sprint(sBegin(s)), "->")
	assert.Contains(fmt.Sprint(StringRBegin(s)), "<-")

	assert.Equal(PrevBidiReader(sEnd(s)).Read(), byte('g'))
	assert.Equal(PrevBidiReader(StringREnd(s)).Read(), byte('a'))

	assert.Equal(
		AdvanceNReader(sBegin(s), 3).Read(),
		byte('d'),
	)
	assert.Equal(
		AdvanceNReader(StringRBegin(s), 3).Read(),
		byte('d'),
	)
	assert.Equal(
		AdvanceNReader(sEnd(s), -3).Read(),
		byte('e'),
	)
	assert.Equal(
		AdvanceNReader(StringREnd(s), -3).Read(),
		byte('c'),
	)

	assert.Equal(sBegin(s).Distance(sEnd(s)), len(s))
	assert.Equal(sEnd(s).Distance(sBegin(s)), -len(s))
	assert.Equal(StringRBegin(s).Distance(StringREnd(s)), len(s))
	assert.Equal(StringREnd(s).Distance(StringRBegin(s)), -len(s))

	assert.True(sBegin(s).Less(sEnd(s)))
	assert.False(sEnd(s).Less(sBegin(s)))
	assert.True(StringRBegin(s).Less(StringREnd(s)))
	assert.False(StringREnd(s).Less(StringRBegin(s)))
}

func TestStringBuilder(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	var bs StringBuilderInserter
	FillN(&bs, 3, 'a')
	assert.Equal(bs.String(), "aaa")

	bs.Reset()
	FillN(&bs, 3, rune('香'))
	assert.Equal(bs.String(), "香香香")

	bs.Reset()
	bs.Delimiter = ","
	FillN(&bs, 3, []byte("abc"))
	assert.Equal(bs.String(), "abc,abc,abc")

	bs.Reset()
	bs.Delimiter = "-->"
	FillN(&bs, 3, "abc")
	assert.Equal(bs.String(), "abc-->abc-->abc")

	bs.Reset()
	bs.Delimiter = ""
	CopyN(IotaReader(1), 5, &bs)
	assert.Equal(bs.String(), "12345")
}

func TestChanIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	ch := make(chan int)
	go func() {
		CopyN(IotaReader(1), 10, ChanWriter(ch))
		close(ch)
	}()
	assert.Equal(Accumulate(ChanReader(ch), ChanEOF, 0), 55)

	ch = make(chan int)
	close(ch)
	assert.True(ChanReader(ch).Eq(ChanEOF))

	ch = make(chan int)
	close(ch)
	assert.True(ChanEOF.Eq(ChanReader(ch)))

	assert.Nil(ChanEOF.Read())
	assert.Equal(ChanEOF.Next(), ChanEOF)
	assert.True(ChanEOF.Eq(ChanEOF))

	ch = make(chan int, 1)
	ch <- 100
	assert.Equal(ChanReader(ch).Read(), 100)

	ch = make(chan int, 2)
	ch <- 42
	ch <- 43
	it := ChanReader(ch)
	it.Next()
	assert.Equal(it.Read(), 43)

	ch = make(chan int, 10)
	CopyN(IotaReader(1), 5, ChanWriter(ch))
	close(ch)
	it = ChanReader(ch)
	it = AdvanceN(it, 3).(InputIter)
	assert.Equal(it.Read(), 4)
	assert.Equal(Distance(it, ChanEOF), 2)

	ch = make(chan int)
	w := ChanWriter(ch)
	assert.Panics(func() { AdvanceN(w, 1) })
	assert.Panics(func() { Distance(w, ChanEOF) })
}

type dummyWriter struct {
	tick int
}

func (w *dummyWriter) Write(b []byte) (int, error) {
	w.tick--
	if w.tick < 0 {
		return 0, errors.New("boom!")
	}
	return len(b), nil
}

func TestIOWriterPanics(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{}, ","))
	})
	assert.Panics(func() {
		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{tick: 1}, ","))
	})
}
