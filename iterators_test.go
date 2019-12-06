package iter_test

import (
	"container/list"
	"fmt"
	"testing"

	"github.com/disksing/iter"
	. "github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	a := make([]int, 100)
	b := SliceBegin(a)
	b.AllowMultiplePass()
	assert.Implements((*ReadWriter)(nil), b)
	assert.Implements((*InputIter)(nil), b)
	assert.Implements((*OutputIter)(nil), b)
	assert.Implements((*ForwardReadWriter)(nil), b)
	assert.Implements((*BidiReadWriter)(nil), b)
	assert.Implements((*RandomReadWriter)(nil), b)

	b1 := b.Next().(RandomReadWriter)
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
	rb1 := SliceRBegin(a).Next().(RandomReadWriter)
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

	e := SliceEnd(a)
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

func listEq(assert *assert.Assertions, lst *list.List, v ...int) {
	assert.True(Equal(ListBegin(lst), ListEnd(lst), begin(v), end(v)))
}

func TestListIterator(t *testing.T) {
	skipAfter(t, 1)

	assert := assert.New(t)
	lst := list.New()
	listEq(assert, lst)
	GenerateN(ListBackInserter(lst), 3, IotaGenerator(1))
	listEq(assert, lst, 1, 2, 3)

	b := ListBegin(lst)
	assert.Implements((*ReadWriter)(nil), b)
	assert.Implements((*InputIter)(nil), b)
	assert.Implements((*OutputIter)(nil), b)
	assert.Implements((*ForwardReadWriter)(nil), b)
	assert.Implements((*BidiReadWriter)(nil), b)
	_, ok := b.(RandomIter)
	assert.False(ok)
	b.AllowMultiplePass()
	b1 := b.Next().(BidiReadWriter)
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
	rb1 := rb.Next().(BidiReadWriter)
	assert.Equal(rb1.Read(), 2)
	assert.True(_eq(rb.Next(), rb1))
	assert.True(_eq(rb1.Prev(), rb))
	assert.True(ListEnd(lst).Prev().Prev().Eq(b1))
	assert.Equal(ListREnd(lst).Prev().(Reader).Read(), 1)

	assert.True(AdvanceN(b, 3).(BidiReader).Eq(ListEnd(lst)))
	assert.True(AdvanceN(ListEnd(lst), -3).(BidiReader).Eq(b))

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

	assert.Contains(fmt.Sprint(StringBegin(s)), "->")
	assert.Contains(fmt.Sprint(StringRBegin(s)), "<-")

	assert.Equal(PrevBidiReader(StringEnd(s)).Read(), byte('g'))
	assert.Equal(PrevBidiReader(StringREnd(s)).Read(), byte('a'))

	assert.Equal(
		AdvanceNReader(StringBegin(s), 3).Read(),
		byte('d'),
	)
	assert.Equal(
		AdvanceNReader(StringRBegin(s), 3).Read(),
		byte('d'),
	)
	assert.Equal(
		AdvanceNReader(StringEnd(s), -3).Read(),
		byte('e'),
	)
	assert.Equal(
		AdvanceNReader(StringREnd(s), -3).Read(),
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
	assert.Equal(Accumulate(ChanReader(ch), iter.ChanEOF, 0), 55)

	ch = make(chan int)
	close(ch)
	assert.True(ChanReader(ch).Eq(iter.ChanEOF))

	ch = make(chan int)
	close(ch)
	assert.True(iter.ChanEOF.Eq(ChanReader(ch)))

	assert.Nil(iter.ChanEOF.Read())
	assert.Equal(iter.ChanEOF.Next(), iter.ChanEOF)
	assert.True(iter.ChanEOF.Eq(iter.ChanEOF))

	ch = make(chan int, 1)
	ch <- 100
	assert.Equal(ChanReader(ch).Read(), 100)

	ch = make(chan int, 2)
	ch <- 42
	ch <- 43
	it := ChanReader(ch)
	it.Next()
	assert.Equal(it.Read(), 43)
}
