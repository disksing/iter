package iter

import (
	"fmt"
	"reflect"
	"strings"
)

type SliceIter struct {
	s        reflect.Value
	i        int
	backward bool
}

func SliceBegin(s interface{}) SliceIter {
	return SliceIter{
		s: reflect.ValueOf(s),
	}
}

func SliceEnd(s interface{}) SliceIter {
	v := reflect.ValueOf(s)
	return SliceIter{
		s: v,
		i: v.Len(),
	}
}

func SliceRBegin(s interface{}) SliceIter {
	v := reflect.ValueOf(s)
	return SliceIter{
		s:        v,
		i:        v.Len() - 1,
		backward: true,
	}
}

func SliceREnd(s interface{}) SliceIter {
	return SliceIter{
		s:        reflect.ValueOf(s),
		i:        -1,
		backward: true,
	}
}

func (it SliceIter) String() string {
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

func (it SliceIter) Read() Any {
	return it.s.Index(it.i).Interface()
}

func (it SliceIter) Write(v Any) {
	it.s.Index(it.i).Set(reflect.ValueOf(v))
}

func (it SliceIter) Equal(it2 Any) bool {
	return it.i == it2.(SliceIter).i
}

func (it SliceIter) Next() ForwardIter {
	return it.AdvanceN(1)
}

func (it SliceIter) Prev() BackwardIter {
	return it.AdvanceN(-1)
}

func (it SliceIter) AdvanceN(n int) RandomIter {
	if it.backward {
		n = -n
	}
	return SliceIter{
		s:        it.s,
		i:        it.i + n,
		backward: it.backward,
	}
}

func (it SliceIter) Distance(it2 RandomIter) int {
	d := it2.(SliceIter).i - it.i
	if it.backward {
		return -d
	}
	return d
}

type sliceBackInserter struct {
	s reflect.Value
}

func SliceBackInserter(s interface{}) ForwardWriter {
	return &sliceBackInserter{
		s: reflect.ValueOf(s).Elem(),
	}
}

func (bi *sliceBackInserter) Next() ForwardIter {
	return bi
}

func (bi *sliceBackInserter) Write(x Any) {
	bi.s.Set(reflect.Append(bi.s, reflect.ValueOf(x)))
}
