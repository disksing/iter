package slices_test

import (
	"fmt"
	"testing"

	"github.com/disksing/iter/v2"
	"github.com/disksing/iter/v2/slices"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator(t *testing.T) {
	assert := assert.New(t)
	a := make([]int, 100)
	b := slices.Begin(a)
	b.AllowMultiplePass()
	b1 := b.Next()
	rb := slices.RBegin(a)
	e := slices.End(a)
	assert.True(b.AdvanceN(100).Eq(e))
	assert.True(e.AdvanceN(-100).Eq(b))
	assert.Equal(100, iter.Distance[int](b, e))
	assert.Equal(-100, iter.Distance[int](e, b))

	re := slices.REnd(a)
	assert.True(rb.AdvanceN(100).Eq(re))
	assert.True(re.AdvanceN(-100).Eq(rb))
	assert.Equal(100, iter.Distance[int](rb, re))
	assert.Equal(-100, iter.Distance[int](re, rb))

	assert.True(b.Less(e))
	assert.True(rb.Less(re))
	assert.True(b.Less(b1))
	assert.False(b1.Less(b))

	e1 := e.Prev()
	assert.NotEqual(fmt.Sprintf("%s", e1), fmt.Sprintf("%s", rb))
}
