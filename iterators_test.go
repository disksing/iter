package iter_test

import (
	"fmt"
	"testing"

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
