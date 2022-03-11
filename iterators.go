package iter

import (
	"container/list"
	"fmt"
	"io"
	"strings"
)

type sliceIter[T any] struct {
	s    []T
	i    int
	step int
}

// SliceBegin returns an iterator to the front element of the slice.
func SliceBegin[T any, It sliceIter[T]](s []T) It {
	return It{s, 0, 1}
}

// SliceEnd returns an iterator to the passed last element of the slice.
func SliceEnd[T any, It sliceIter[T]](s []T) It {
	return It{s, len(s), 1}
}

// SliceRBegin returns an iterator to the back element of the slice.
func SliceRBegin[T any, It sliceIter[T]](s []T) It {
	return It{s, len(s) - 1, -1}
}

// SliceREnd returns an iterator to the passed first element of the slice.
func SliceREnd[T any, It sliceIter[T]](s []T) It {
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

// SliceBackInserter returns an OutputIter to append elements to the back of the
// slice.
func SliceBackInserter[T any](s *[]T) sliceBackInserter[T] {
	return sliceBackInserter[T]{
		s: s,
	}
}

func (bi sliceBackInserter[T]) Write(x T) {
	*bi.s = append(*bi.s, x)
}

// listIter is an iterator works with list.List.
type listIter[T any] struct {
	l        *list.List
	e        *list.Element
	backward bool
}

// ListBegin returns an iterator to the front element of the list.
func ListBegin[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l: l,
		e: l.Front(),
	}
}

// ListEnd returns an iterator to the passed last element of the list.
func ListEnd[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l: l,
	}
}

// ListRBegin returns an iterator to the back element of the list.
func ListRBegin[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l:        l,
		e:        l.Back(),
		backward: true,
	}
}

// ListREnd returns an iterator to the passed first element of the list.
func ListREnd[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l:        l,
		backward: true,
	}
}

func (l listIter[T]) Eq(x listIter[T]) bool {
	return l.e == x.e
}

func (l listIter[T]) AllowMultiplePass() {}

func (l listIter[T]) Next() listIter[T] {
	var e *list.Element
	if l.backward {
		e = l.e.Prev()
	} else {
		e = l.e.Next()
	}
	return listIter[T]{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l listIter[T]) Prev() listIter[T] {
	var e *list.Element
	switch {
	case l.e == nil && l.backward:
		e = l.l.Front()
	case l.e == nil && !l.backward:
		e = l.l.Back()
	case l.e != nil && l.backward:
		e = l.e.Next()
	case l.e != nil && !l.backward:
		e = l.e.Prev()
	}
	return listIter[T]{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l listIter[T]) Read() T {
	return l.e.Value.(T)
}

func (l listIter[T]) Write(x T) {
	l.e.Value = x
}

// ListBackInserter returns an OutputIter to insert elements to the back of the
// list.
func ListBackInserter[T any](l *list.List) listBackInserter[T] {
	return listBackInserter[T]{l: l}
}

type listBackInserter[T any] struct {
	l *list.List
}

func (li listBackInserter[T]) Write(x T) {
	li.l.PushBack(x)
}

// ListInserter returns an OutputIter to insert elements before a node.
func ListInserter[T any](l *list.List, e *list.Element) listInserter[T] {
	return listInserter[T]{l: l, e: e}
}

type listInserter[T any] struct {
	l *list.List
	e *list.Element
}

func (li listInserter[T]) Write(x T) {
	li.l.InsertBefore(x, li.e)
}

// stringIter is the iterator to access a string in bytes. To travise a string
// by rune, convert the string to []rune then use SliceIter.
type stringIter struct {
	s    string
	i    int
	step int
}

// StringBegin returns an iterator to the front element of the string.
func StringBegin(s string) stringIter {
	return stringIter{
		s:    s,
		step: 1,
	}
}

// StringEnd returns an iterator to the passed last element of the string.
func StringEnd(s string) stringIter {
	return stringIter{
		s:    s,
		i:    len(s),
		step: 1,
	}
}

// StringRBegin returns an iterator to the back element of the string.
func StringRBegin(s string) stringIter {
	return stringIter{
		s:    s,
		i:    len(s) - 1,
		step: -1,
	}
}

// StringREnd returns an iterator to the passed first element of the string.
func StringREnd(s string) stringIter {
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

type chanReader[T any, C Chan[T]] struct {
	ch    C
	cur   T
	read1 bool
	eof   bool
}

func (cr *chanReader[T]) recv() {
	v, ok := <-cr.ch
	cr.cur, cr.read1, cr.eof = v, true, !ok
}

func (cr *chanReader[T]) Read() T {
	if !cr.read1 {
		cr.recv()
	}
	return cr.cur
}

func (cr *chanReader[T]) Next() *chanReader[T] {
	if !cr.read1 {
		cr.recv()
	}
	if !cr.eof {
		cr.recv()
	}
	return cr
}

func (cr *chanReader[T]) Eq(x *chanReader[T]) bool {
	if !cr.read1 {
		cr.recv()
	}
	return cr.eof && x == nil
}

type chanWriter[T any, C Chan[T]] struct {
	ch C
}

func (cr *chanWriter[T]) Write(x T) {
	cr.ch <- x
}

// // ChanReader returns an InputIter that reads from a channel.
// func ChanReader[T any, C Chan[T]](c C) *chanReader[T, C] {
// 	return &chanReader[T, C]{
// 		ch: c,
// 	}
// }

// // ChanWriter returns an OutIter that writes to a channel.
// func ChanWriter[T any, C Chan[T]](c C) *chanWriter[T, C] {
// 	return &chanWriter[T, C]{
// 		ch: c,
// 	}
// }

type ioWriter struct {
	w         io.Writer
	written   bool
	delimiter []byte
}

func (w *ioWriter) Write(x any) {
	if w.written && len(w.delimiter) > 0 {
		_, err := w.w.Write(w.delimiter)
		if err != nil {
			panic(err)
		}
	} else {
		w.written = true
	}

	_, err := fmt.Fprint(w.w, x)
	if err != nil {
		panic(err)
	}
}

// IOWriter returns an OutputIter that writes values to an io.Writer.
// It panics if meet any error.
func IOWriter(w io.Writer, delimiter string) *ioWriter {
	return &ioWriter{w: w, delimiter: []byte(delimiter)}
}
