package strs

import (
	"fmt"
	"strings"

	"github.com/disksing/iter/v2"
)

// stringIter is the iterator to access a string in bytes. To travise a string
// by rune, convert the string to []rune then use SliceIter.
type stringIter struct {
	s    string
	i    int
	step int
}

// Begin returns an iterator to the front element of the string.
func Begin(s string) stringIter {
	return stringIter{
		s:    s,
		step: 1,
	}
}

// End returns an iterator to the passed last element of the string.
func End(s string) stringIter {
	return stringIter{
		s:    s,
		i:    len(s),
		step: 1,
	}
}

// RBegin returns an iterator to the back element of the string.
func RBegin(s string) stringIter {
	return stringIter{
		s:    s,
		i:    len(s) - 1,
		step: -1,
	}
}

// REnd returns an iterator to the passed first element of the string.
func REnd(s string) stringIter {
	return stringIter{
		s:    s,
		i:    -1,
		step: -1,
	}
}

func (it stringIter) String() string {
	dir := "->"
	if it.step < 0 {
		dir = "<-"
	}
	return fmt.Sprintf("%s@%d%s", it.s, it.i, dir)
}

func (it stringIter) Read() byte {
	return it.s[it.i]
}

func (it stringIter) Eq(it2 stringIter) bool {
	return it.i == it2.i
}

func (it stringIter) AllowMultiplePass() {}

func (it stringIter) Next() stringIter {
	return it.AdvanceN(1)
}

func (it stringIter) Prev() stringIter {
	return it.AdvanceN(-1)
}

func (it stringIter) AdvanceN(n int) stringIter {
	return stringIter{
		s:    it.s,
		i:    it.i + n*it.step,
		step: it.step,
	}
}

func (it stringIter) Distance(it2 stringIter) int {
	return (it2.i - it.i) * it.step
}

func (it stringIter) Less(it2 stringIter) bool {
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

// MakeString creates a string by range spesified by [first, last). The value
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
