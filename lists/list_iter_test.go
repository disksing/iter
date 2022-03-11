package lists_test

import (
	"container/list"
	"testing"

	"github.com/disksing/iter/v2"
	. "github.com/disksing/iter/v2/algo"
	. "github.com/disksing/iter/v2/lists"
	"github.com/disksing/iter/v2/slices"
	"github.com/stretchr/testify/assert"
)

func listEq[T comparable](assert *assert.Assertions, lst *list.List, v ...T) {
	end := slices.End(v)
	assert.True(Equal[T](Begin[T](lst), End[T](lst), slices.Begin(v), &end))
}

func TestListIterator(t *testing.T) {
	assert := assert.New(t)
	lst := list.New()
	listEq[int](assert, lst)
	GenerateN(ListBackInserter[int](lst), 3, iter.IotaGenerator(1))
	listEq(assert, lst, 1, 2, 3)

	b := Begin[int](lst)
	b.AllowMultiplePass()
	b1 := b.Next()
	assert.Equal(b, b1.Prev())

	rb := RBegin[int](lst)
	assert.Equal(rb.Read(), 3)
	rb1 := rb.Next()
	assert.Equal(rb1.Read(), 2)
	assert.True(rb.Next().Eq(rb1))
	assert.True(rb1.Prev().Eq(rb))
	assert.True(End[int](lst).Prev().Prev().Eq(b1))
	assert.Equal(REnd[int](lst).Prev().Read(), 1)

	assert.True(iter.AdvanceN[int](b, 3).Eq(End[int](lst)))
	assert.True(iter.AdvanceN[int](End[int](lst), -3).Eq(b))

	b.Write(2)
	listEq(assert, lst, 2, 2, 3)
	ListInserter[int](lst, lst.Back()).Write(4)
	listEq(assert, lst, 2, 2, 4, 3)
}
