package iter

import "reflect"

type Iter interface {
	Next() Iter
	Get() interface{}
	Set(interface{})
}

type SliceIter struct {
	s reflect.Value
	i int
}

func SliceBegin(s interface{}) SliceIter {
	return SliceIter{
		s: reflect.ValueOf(s),
	}
}

func SliceEnd(s interface{}) SliceIter {
	v := reflect.ValueOf(s)
	return SliceIter {
		s: s,
		i: v.Len(),
	}
}

func (it SliceIter) Next() Iter {
	return SliceIter{
		s: it.s,
		i:it.i+1,
	}
}

func (it SliceIter) Get() interface{} {
	return it.s.At(it.i)
}

func (it SliceIter) Set(v interface{}) {
	it.s.At(it.i).Set(v)
}