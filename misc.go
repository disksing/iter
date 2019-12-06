package iter

import (
	"container/list"
	"math/rand"
	"reflect"
	"strings"
)

type iotaReader struct {
	x Any
}

func (r iotaReader) Read() Any {
	return r.x
}

func (r iotaReader) Next() Incrementable {
	return iotaReader{x: _inc(r.x)}
}

func (r iotaReader) Eq(Iter) bool {
	return false
}

// IotaReader creates an InputIter that returns [x, x+1, x+2...).
func IotaReader(x Any) InputIter {
	return iotaReader{x: x}
}

// IotaGenerator creates a Generator that returns [x, x+1, x+2...).
func IotaGenerator(x Any) Generator {
	r := IotaReader(x)
	return func() Any {
		v := r.Read()
		r = NextInputIter(r)
		return v
	}
}

type repeatReader struct {
	x Any
}

func (r repeatReader) Read() Any { return r.x }

func (r repeatReader) Next() Incrementable { return r }

func (r repeatReader) Eq(Iter) bool { return false }

// RepeatReader creates an InputIter that returns [x, x, x...).
func RepeatReader(x Any) InputIter {
	return repeatReader{x: x}
}

// RepeatGenerator creates an Generator that returns [x, x, x...).
func RepeatGenerator(x Any) Generator {
	return func() Any { return x }
}

type randomGenerator struct {
	candidates []interface{}
}

// RandomGenerator creates a generator that returns random item of a slice.
func RandomGenerator(s interface{}, r *rand.Rand) Generator {
	v := reflect.ValueOf(s)
	l := v.Len()
	return func() Any { return v.Index(r.Intn(l)).Interface() }
}

// Erase removes a range from a container.
// c should be settable (*[]T or *list.List).
// With 2 iterators arguments, it removes [it1, it2).
// With 1 iterator argument, it removes [it, end).
// With no iterator argument, it remvoes [begin, end).
func Erase(c interface{}, it ...Iter) {
	if len(it) > 2 {
		panic("too many iterators, expect <=2")
	}
	if val := reflect.ValueOf(c); val.Elem().Type().Kind() == reflect.Slice {
		v := val.Elem()
		switch len(it) {
		case 0:
			v.Set(v.Slice(0, 0))
		case 1:
			l := Distance(SliceBegin(c), it[0])
			v.Set(v.Slice(0, l))
		case 2:
			l, h := Distance(SliceBegin(c), it[0]), Distance(SliceBegin(c), it[1])
			v.Set(reflect.AppendSlice(v.Slice(0, l), v.Slice(h, v.Len())))
		}
	} else if lst, ok := c.(*list.List); ok {
		switch len(it) {
		case 0:
			for lst.Len() > 0 {
				lst.Remove(lst.Front())
			}
		case 1:
			if it, ok := it[0].(listIter); ok {
				for e := it.e; e != nil; {
					next := e.Next()
					it.l.Remove(e)
					e = next
				}
				return
			}
			panic("iterator is not listIter")
		case 2:
			if it1, ok := it[0].(listIter); ok {
				if it2, ok := it[1].(listIter); ok {
					for e := it1.e; e != it2.e; {
						next := e.Next()
						it1.l.Remove(e)
						e = next
					}
					return
				}
			}
			panic("iterator is not listIter")
		}
	} else {
		panic("c is not *[]T or *list.List")
	}
}

// MakeString creates a string by range spesified by [first, last). The value
// type should be byte or rune.
func MakeString(first, last ForwardReader) string {
	var s strings.Builder
	for ; _ne(first, last); first = NextForwardReader(first) {
		switch v := first.Read().(type) {
		case byte:
			s.WriteByte(v)
		case rune:
			s.WriteRune(v)
		}
	}
	return s.String()
}
