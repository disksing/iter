package strs

import (
	"fmt"
	"strings"

	"github.com/disksing/iter/v2"
)

// Iterator accesses the bytes in a string. To traverse runes, convert the
// string to []rune and use slices.Iterator.
type Iterator struct {
	s    string
	i    int
	step int
}

// Begin returns an iterator to the front element of the string.
func Begin(s string) Iterator {
	return Iterator{
		s:    s,
		step: 1,
	}
}

// End returns an iterator to the passed last element of the string.
func End(s string) Iterator {
	return Iterator{
		s:    s,
		i:    len(s),
		step: 1,
	}
}

// RBegin returns an iterator to the back element of the string.
func RBegin(s string) Iterator {
	return Iterator{
		s:    s,
		i:    len(s) - 1,
		step: -1,
	}
}

// REnd returns an iterator to the passed first element of the string.
func REnd(s string) Iterator {
	return Iterator{
		s:    s,
		i:    -1,
		step: -1,
	}
}

func (it Iterator) String() string {
	dir := "->"
	if it.step < 0 {
		dir = "<-"
	}
	return fmt.Sprintf("%s@%d%s", it.s, it.i, dir)
}

func (it Iterator) Read() byte {
	return it.s[it.i]
}

func (it Iterator) Eq(it2 Iterator) bool {
	return it.i == it2.i
}

func (it Iterator) AllowMultiplePass() {}

func (it Iterator) Next() Iterator {
	return it.AdvanceN(1)
}

func (it Iterator) Prev() Iterator {
	return it.AdvanceN(-1)
}

func (it Iterator) AdvanceN(n int) Iterator {
	return Iterator{
		s:    it.s,
		i:    it.i + n*it.step,
		step: it.step,
	}
}

func (it Iterator) Distance(it2 Iterator) int {
	return (it2.i - it.i) * it.step
}

func (it Iterator) Less(it2 Iterator) bool {
	if it.step > 0 {
		return it.i < it2.i
	}
	return it.i > it2.i
}

// StringBuilderInserter is an OutputIter that wraps a strings.Builder.
type StringBuilderInserter[T any] struct {
	strings.Builder
	Delimiter string
}

func (si *StringBuilderInserter[T]) Write(x T) {
	if si.Builder.Len() > 0 && si.Delimiter != "" {
		si.Builder.WriteString(si.Delimiter)
	}
	switch v := any(x).(type) {
	case byte:
		si.Builder.WriteByte(v)
	case rune:
		si.Builder.WriteRune(v)
	case []byte:
		si.Builder.Write(v)
	case string:
		si.Builder.WriteString(v)
	default:
		si.Builder.WriteString(fmt.Sprint(x))
	}
}

// MakeString creates a string from the range specified by [first, last). The value
// type should be byte or rune.
func MakeString[T byte | rune, It iter.ForwardReader[T, It]](first, last It) string {
	var s strings.Builder
	for ; !first.Eq(last); first = first.Next() {
		switch v := any(first.Read()).(type) {
		case byte:
			s.WriteByte(v)
		case rune:
			s.WriteRune(v)
		}
	}
	return s.String()
}
