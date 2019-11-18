package iter_test

import (
	"testing"

	. "github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

func TestAnyNoneAll(t *testing.T) {
	// https://en.cppreference.com/w/cpp/algorithm/all_any_none_of
	assert := assert.New(t)
	a := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	begin, end := SliceBegin(a), SliceEnd(a)
	assert.True(AllOf(begin, end, func(x Any) bool { return x.(int)%2 == 0 }))
	assert.True(NoneOf(begin, end, func(x Any) bool { return x.(int)%2 == 1 }))
	assert.True(AnyOf(begin, end, func(x Any) bool { return x.(int)%7 == 0 }))
}

func TestForEach(t *testing.T) {
	// https://en.cppreference.com/w/cpp/algorithm/for_each
	assert := assert.New(t)
	a := []int{3, 4, 2, 8, 15, 267}
	begin, end := SliceBegin(a), SliceEnd(a)
	ForEach(begin, end, func(it Iter) {
		it2 := it.(ReadWriter)
		it2.Write(it2.Read().(int) + 1)
	})
	assert.Equal([]int{4, 5, 3, 9, 16, 268}, a)
	var sum int
	ForEach(begin, end, func(it Iter) { sum += it.(Readable).Read().(int) })
	assert.Equal(sum, 305)
}

func TestForEachN(t *testing.T) {
	// https://en.cppreference.com/w/cpp/algorithm/for_each_n
	assert := assert.New(t)
	a := []int{1, 2, 3, 4, 5}
	ForEachN(SliceBegin(a), 3, func(it Iter) {
		it2 := it.(ReadWriter)
		it2.Write(it2.Read().(int) * 2)
	})
	assert.Equal([]int{2, 4, 6, 4, 5}, a)
}

func TestCount(t *testing.T) {
	// https://en.cppreference.com/w/cpp/algorithm/count
	assert := assert.New(t)
	a := []int{1, 2, 3, 4, 4, 3, 7, 8, 9, 10}
	begin, end := SliceBegin(a), SliceEnd(a)
	assert.Equal(Count(begin, end, 3), 2)
	assert.Equal(Count(begin, end, 5), 0)
	assert.Equal(CountIf(begin, end, func(x Any) bool { return x.(int)%3 == 0 }), 3)
}

func TestMismatch(t *testing.T) {
	// https://en.cppreference.com/w/cpp/algorithm/mismatch
	assert := assert.New(t)
	mirrorEnds := func(s string) string {
		p, _ := Mismatch(StringBegin(s), StringEnd(s), StringRBegin(s), StringREnd(s))
		return MakeString(StringBegin(s), p)
	}
	assert.Equal(mirrorEnds("abXYZba"), "ab")
	assert.Equal(mirrorEnds("abca"), "a")
	assert.Equal(mirrorEnds("aba"), "aba")
}
