package slices

import (
	"fmt"
	"strings"
)

// Iterator is a random-access iterator over a slice.
type Iterator[T any] struct {
	s    []T
	i    int
	step int
}

// Begin returns an iterator to the front element of the slice.
func Begin[T any](s []T) Iterator[T] {
	return Iterator[T]{s: s, i: 0, step: 1}
}

// End returns an iterator to the passed last element of the slice.
func End[T any](s []T) Iterator[T] {
	return Iterator[T]{s: s, i: len(s), step: 1}
}

// RBegin returns an iterator to the back element of the slice.
func RBegin[T any](s []T) Iterator[T] {
	return Iterator[T]{s: s, i: len(s) - 1, step: -1}
}

// REnd returns an iterator to the passed first element of the slice.
func REnd[T any](s []T) Iterator[T] {
	return Iterator[T]{s: s, i: -1, step: -1}
}

func (it Iterator[T]) Read() T {
	return it.s[it.i]
}

func (it Iterator[T]) Write(v T) {
	it.s[it.i] = v
}

func (it Iterator[T]) Eq(it2 Iterator[T]) bool {
	return it.i == it2.i
}

func (it Iterator[T]) Less(it2 Iterator[T]) bool {
	if it.step < 0 {
		return it.i > it2.i
	}
	return it.i < it2.i
}

func (it Iterator[T]) Next() Iterator[T] {
	return it.AdvanceN(1)
}

func (it Iterator[T]) Prev() Iterator[T] {
	return it.AdvanceN(-1)
}

func (it Iterator[T]) AdvanceN(n int) Iterator[T] {
	return Iterator[T]{
		s:    it.s,
		i:    it.i + n*it.step,
		step: it.step,
	}
}

func (it Iterator[T]) String() string {
	dir := "->"
	if it.step < 0 {
		dir = "<-"
	}
	var buf []string
	for i := 0; i < 64 && i < len(it.s); i++ {
		buf = append(buf, fmt.Sprintf("%v", it.s[i]))
	}
	if len(it.s) > 64 {
		buf = append(buf, "...")
	}
	return fmt.Sprintf("[%v](len=%d,cap=%d)@%d%s", strings.Join(buf, ","), len(it.s), cap(it.s), it.i, dir)
}

func (it Iterator[T]) AllowMultiplePass() {}

func (it Iterator[T]) Distance(it2 Iterator[T]) int {
	return (it2.i - it.i) * it.step
}

// BackInserter is an output iterator that appends values to a slice.
type BackInserter[T any] struct {
	s *[]T
}

// Appender returns an OutputIter to append elements to the back of the
// slice.
func Appender[T any](s *[]T) BackInserter[T] {
	return BackInserter[T]{
		s: s,
	}
}

func (bi BackInserter[T]) Write(x T) {
	*bi.s = append(*bi.s, x)
}
