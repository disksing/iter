package slices

import (
	"fmt"
	"strings"
)

type sliceIter[T any] struct {
	s    []T
	i    int
	step int
}

// Begin returns an iterator to the front element of the slice.
func Begin[T any, It sliceIter[T]](s []T) It {
	return It{s, 0, 1}
}

// End returns an iterator to the passed last element of the slice.
func End[T any, It sliceIter[T]](s []T) It {
	return It{s, len(s), 1}
}

// RBegin returns an iterator to the back element of the slice.
func RBegin[T any, It sliceIter[T]](s []T) It {
	return It{s, len(s) - 1, -1}
}

// REnd returns an iterator to the passed first element of the slice.
func REnd[T any, It sliceIter[T]](s []T) It {
	return It{s, -1, -1}
}

func (it sliceIter[T]) Read() T {
	return it.s[it.i]
}

func (it sliceIter[T]) Write(v T) {
	it.s[it.i] = v
}

func (it sliceIter[T]) Eq(it2 sliceIter[T]) bool {
	return it.i == it2.i
}

func (it sliceIter[T]) Less(it2 sliceIter[T]) bool {
	if it.step < 0 {
		return it.i > it2.i
	}
	return it.i < it2.i
}

func (it sliceIter[T]) Next() sliceIter[T] {
	return it.AdvanceN(1)
}

func (it sliceIter[T]) Prev() sliceIter[T] {
	return it.AdvanceN(-1)
}

func (it sliceIter[T]) AdvanceN(n int) sliceIter[T] {
	return sliceIter[T]{
		s:    it.s,
		i:    it.i + n*it.step,
		step: it.step,
	}
}

func (it sliceIter[T]) String() string {
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

func (it sliceIter[T]) AllowMultiplePass() {}

func (it sliceIter[T]) Distance(it2 sliceIter[T]) int {
	return (it2.i - it.i) * it.step
}

type sliceBackInserter[T any] struct {
	s *[]T
}

// Appender returns an OutputIter to append elements to the back of the
// slice.
func Appender[T any](s *[]T) sliceBackInserter[T] {
	return sliceBackInserter[T]{
		s: s,
	}
}

func (bi sliceBackInserter[T]) Write(x T) {
	*bi.s = append(*bi.s, x)
}
