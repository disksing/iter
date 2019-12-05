package iter

import (
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
		r = r.Next().(InputIter)
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
	val := reflect.ValueOf(c)
	switch val.Elem().Type().Kind() {
	case reflect.Slice:
		switch len(it) {
		case 0:
			sliceErase(c, SliceBegin(val.Elem().Interface()), SliceEnd(val.Elem().Interface()))
		case 1:
			sliceErase(c, it[0], SliceEnd(val.Elem().Interface()))
		case 2:
			sliceErase(c, it[0], it[1])
		}
	}
}

func sliceErase(s interface{}, first, last Iter) {
	v := reflect.ValueOf(s).Elem()
	begin := SliceBegin(s)
	m, h := Distance(begin, first), Distance(begin, last)
	v.Set(reflect.AppendSlice(v.Slice(0, m), v.Slice(h, v.Len())))
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
