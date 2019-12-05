package iter

import "container/list"

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
	return &listIter{
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
	return &listIter{
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
