package lists

import "container/list"

// listIter is an iterator works with list.List.
type listIter[T any] struct {
	l        *list.List
	e        *list.Element
	backward bool
}

// Begin returns an iterator to the front element of the list.
func Begin[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l: l,
		e: l.Front(),
	}
}

// End returns an iterator to the passed last element of the list.
func End[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l: l,
	}
}

// RBegin returns an iterator to the back element of the list.
func RBegin[T any](l *list.List) listIter[T] {
	return listIter[T]{
		l:        l,
		e:        l.Back(),
		backward: true,
	}
}

// REnd returns an iterator to the passed first element of the list.
func REnd[T any](l *list.List) listIter[T] {
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
