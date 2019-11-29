package iter

import (
	"fmt"
	"strings"
)

// StringIter is the iterator to access a string in bytes. To travise a string
// by rune, convert the string to []rune then use SliceIter.
type StringIter struct {
	s        string
	i        int
	backward bool
}

func StringBegin(s string) StringIter {
	return StringIter{
		s: s,
	}
}

func StringEnd(s string) StringIter {
	return StringIter{
		s: s,
		i: len(s),
	}
}

func StringRBegin(s string) StringIter {
	return StringIter{
		s:        s,
		i:        len(s) - 1,
		backward: true,
	}
}

func StringREnd(s string) StringIter {
	return StringIter{
		s:        s,
		i:        -1,
		backward: true,
	}
}

func (it StringIter) String() string {
	dir := "->"
	if it.backward {
		dir = "<-"
	}
	return fmt.Sprintf("%s@%d%s", it.s, it.i, dir)
}

func (it StringIter) Read() Any {
	return it.s[it.i]
}

func (it StringIter) Eq(it2 Any) bool {
	return it.i == it2.(StringIter).i
}

func (it StringIter) Next() ForwardIter {
	return it.AdvanceN(1)
}

func (it StringIter) Prev() BackwardIter {
	return it.AdvanceN(-1)
}

func (it StringIter) AdvanceN(n int) RandomIter {
	if it.backward {
		n = -n
	}
	return StringIter{
		s:        it.s,
		i:        it.i + n,
		backward: it.backward,
	}
}

func (it StringIter) Distance(it2 RandomIter) int {
	d := it2.(StringIter).i - it.i
	if it.backward {
		return -d
	}
	return d
}

// MakeString creates a string by range spesified by [first, last). The value
// type should be byte or rune.
func MakeString(first, last ForwardReader) string {
	var s strings.Builder
	for ; _ne(first, last); first = NextReader(first) {
		switch v := first.Read().(type) {
		case byte:
			s.WriteByte(v)
		case rune:
			s.WriteRune(v)
		}
	}
	return s.String()
}
