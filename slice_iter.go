package iter

import (
	"fmt"
	"reflect"
	"strings"
)

type sliceIter struct {
	s        reflect.Value
	i        int
	backward bool
}

// SliceBegin returns an iterator to the front element of the slice.
func SliceBegin(s interface{}) RandomReadWriter {
	return sliceIter{
		s: reflect.ValueOf(s),
	}
}

// SliceEnd returns an iterator to the passed last element of the slice.
func SliceEnd(s interface{}) RandomReadWriter {
	v := reflect.ValueOf(s)
	return sliceIter{
		s: v,
		i: v.Len(),
	}
}

// SliceRBegin returns an iterator to the back element of the slice.
func SliceRBegin(s interface{}) RandomReadWriter {
	v := reflect.ValueOf(s)
	return sliceIter{
		s:        v,
		i:        v.Len() - 1,
		backward: true,
	}
}

// SliceREnd returns an iterator to the passed first element of the slice.
func SliceREnd(s interface{}) RandomReadWriter {
	return sliceIter{
		s:        reflect.ValueOf(s),
		i:        -1,
		backward: true,
	}
}

func (it sliceIter) String() string {
	dir := "->"
	if it.backward {
		dir = "<-"
	}
	var buf []string
	for i := 0; i < 64 && i < it.s.Len(); i++ {
		buf = append(buf, fmt.Sprintf("%v", it.s.Index(i)))
	}
	if it.s.Len() > 64 {
		buf = append(buf, "...")
	}
	return fmt.Sprintf("[%v](len=%d,cap=%d)@%d%s", strings.Join(buf, ","), it.s.Len(), it.s.Cap(), it.i, dir)
}

func (it sliceIter) Read() Any {
	return it.s.Index(it.i).Interface()
}

func (it sliceIter) Write(v Any) {
	it.s.Index(it.i).Set(reflect.ValueOf(v))
}

func (it sliceIter) Eq(it2 Iter) bool {
	return it.i == it2.(sliceIter).i
}

func (it sliceIter) AllowMultiplePass() {}

func (it sliceIter) Less(it2 Iter) bool {
	if it.backward {
		return it.i > it2.(sliceIter).i
	}
	return it.i < it2.(sliceIter).i
}

func (it sliceIter) Next() Incrementable {
	return it.AdvanceN(1)
}

func (it sliceIter) Prev() BidiIter {
	return it.AdvanceN(-1)
}

func (it sliceIter) AdvanceN(n int) RandomIter {
	if it.backward {
		n = -n
	}
	return sliceIter{
		s:        it.s,
		i:        it.i + n,
		backward: it.backward,
	}
}

func (it sliceIter) Distance(it2 RandomIter) int {
	d := it2.(sliceIter).i - it.i
	if it.backward {
		return -d
	}
	return d
}

type sliceBackInserter struct {
	s reflect.Value
}

// SliceBackInserter returns an OutputIter to append elements to the back of the
// slice.
func SliceBackInserter(s interface{}) OutputIter {
	return &sliceBackInserter{
		s: reflect.ValueOf(s).Elem(),
	}
}

func (bi *sliceBackInserter) Write(x Any) {
	bi.s.Set(reflect.Append(bi.s, reflect.ValueOf(x)))
}

// SliceErase removes elements from the end of a slice, starting from it.
func SliceErase(s interface{}, it Iter) {
	v := reflect.ValueOf(s).Elem()
	v.Set(v.Slice(0, Distance(SliceBegin(v), it)))
}
