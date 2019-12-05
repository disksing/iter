package iter

import (
	"fmt"
	"strings"
)

// stringIter is the iterator to access a string in bytes. To travise a string
// by rune, convert the string to []rune then use SliceIter.
type stringIter struct {
	s        string
	i        int
	backward bool
}

// StringBegin returns an iterator to the front element of the string.
func StringBegin(s string) RandomReader {
	return stringIter{
		s: s,
	}
}

// StringEnd returns an iterator to the passed last element of the string.
func StringEnd(s string) RandomReader {
	return stringIter{
		s: s,
		i: len(s),
	}
}

// StringRBegin returns an iterator to the back element of the string.
func StringRBegin(s string) RandomReader {
	return stringIter{
		s:        s,
		i:        len(s) - 1,
		backward: true,
	}
}

// StringREnd returns an iterator to the passed first element of the string.
func StringREnd(s string) RandomReader {
	return stringIter{
		s:        s,
		i:        -1,
		backward: true,
	}
}

func (it stringIter) String() string {
	dir := "->"
	if it.backward {
		dir = "<-"
	}
	return fmt.Sprintf("%s@%d%s", it.s, it.i, dir)
}

func (it stringIter) Read() Any {
	return it.s[it.i]
}

func (it stringIter) Eq(it2 Iter) bool {
	return it.i == it2.(stringIter).i
}

func (it stringIter) AllowMultiplePass() {}

func (it stringIter) Next() Incrementable {
	return it.AdvanceN(1)
}

func (it stringIter) Prev() BidiIter {
	return it.AdvanceN(-1)
}

func (it stringIter) AdvanceN(n int) RandomIter {
	if it.backward {
		n = -n
	}
	return stringIter{
		s:        it.s,
		i:        it.i + n,
		backward: it.backward,
	}
}

func (it stringIter) Distance(it2 RandomIter) int {
	d := it2.(stringIter).i - it.i
	if it.backward {
		return -d
	}
	return d
}

func (it stringIter) Less(it2 Iter) bool {
	if it.backward {
		return it.i > it2.(stringIter).i
	}
	return it.i < it2.(stringIter).i
}

// StringBuilderInserter is an OutputIter that wraps a strings.Builder.
type StringBuilderInserter struct {
	strings.Builder
}

func (si *StringBuilderInserter) Write(x Any) {
	switch v := x.(type) {
	case byte:
		si.Builder.WriteByte(v)
	case rune:
		si.Builder.WriteRune(v)
	case string:
		si.Builder.WriteString(v)
	default:
		panic("unknown item type")
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
