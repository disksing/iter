package iter

import (
	"container/list"
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

// listIter is an iterator works with list.List.
type listIter struct {
	l        *list.List
	e        *list.Element
	backward bool
}

// ListBegin returns an iterator to the front element of the list.
func ListBegin(l *list.List) BidiReadWriter {
	return listIter{
		l: l,
		e: l.Front(),
	}
}

// ListEnd returns an iterator to the passed last element of the list.
func ListEnd(l *list.List) BidiReadWriter {
	return listIter{
		l: l,
	}
}

// ListRBegin returns an iterator to the back element of the list.
func ListRBegin(l *list.List) BidiReadWriter {
	return listIter{
		l:        l,
		e:        l.Back(),
		backward: true,
	}
}

// ListREnd returns an iterator to the passed first element of the list.
func ListREnd(l *list.List) BidiReadWriter {
	return listIter{
		l:        l,
		backward: true,
	}
}

func (l listIter) Eq(x Iter) bool {
	return l.e == x.(listIter).e
}

func (l listIter) AllowMultiplePass() {}

func (l listIter) Next() Incrementable {
	var e *list.Element
	if l.backward {
		e = l.e.Prev()
	} else {
		e = l.e.Next()
	}
	return listIter{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l listIter) Prev() BidiIter {
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
	return listIter{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l listIter) Read() Any {
	return l.e.Value
}

func (l listIter) Write(x Any) {
	l.e.Value = x
}

// ListBackInserter returns an OutputIter to insert elements to the back of the
// list.
func ListBackInserter(l *list.List) OutputIter {
	return listBackInserter{l: l}
}

type listBackInserter struct {
	l *list.List
}

func (li listBackInserter) Write(x Any) {
	li.l.PushBack(x)
}

// ListInserter returns an OutputIter to insert elements before a node.
func ListInserter(l *list.List, e *list.Element) OutputIter {
	return listInserter{l: l, e: e}
}

type listInserter struct {
	l *list.List
	e *list.Element
}

func (li listInserter) Write(x Any) {
	li.l.InsertBefore(x, li.e)
}

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
	Delimiter string
}

func (si *StringBuilderInserter) Write(x Any) {
	if si.Builder.Len() > 0 && si.Delimiter != "" {
		si.Builder.WriteString(si.Delimiter)
	}
	switch v := x.(type) {
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

type eof int

func (e eof) Eq(x Any) bool {
	if _, ok := x.(eof); ok {
		return true
	}
	return _eq(x, e)
}

func (e eof) Next() Incrementable { return e }

func (e eof) Read() Any { return nil }

// ChanEOF is a sentinel iterator to terminate chan reader.
var ChanEOF InputIter = eof(0)

type chanReader struct {
	ch    reflect.Value
	cur   interface{}
	read1 bool
	eof   bool
}

func (cr *chanReader) recv() {
	v, ok := cr.ch.Recv()
	cr.cur, cr.read1, cr.eof = v.Interface(), true, !ok
}

func (cr *chanReader) Read() Any {
	if !cr.read1 {
		cr.recv()
	}
	return cr.cur
}

func (cr *chanReader) Next() Incrementable {
	if !cr.read1 {
		cr.recv()
	}
	if !cr.eof {
		cr.recv()
	}
	return cr
}

func (cr *chanReader) Eq(x Any) bool {
	if !cr.read1 {
		cr.recv()
	}
	return cr.eof && x == ChanEOF
}

type chanWriter struct {
	ch reflect.Value
}

func (cr *chanWriter) Write(x Any) {
	cr.ch.Send(reflect.ValueOf(x))
}

// ChanReader returns an InputIter that reads from a channel.
func ChanReader(c interface{}) InputIter {
	return &chanReader{
		ch: reflect.ValueOf(c),
	}
}

// ChanWriter returns an OutIter that writes to a channel.
func ChanWriter(c interface{}) OutputIter {
	return &chanWriter{
		ch: reflect.ValueOf(c),
	}
}
