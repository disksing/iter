package lists

import "container/list"

// Iterator is a bidirectional iterator over a list.List.
//
// list.List stores values as any, so Read panics if the current element does
// not contain a T.
type Iterator[T any] struct {
	l        *list.List
	e        *list.Element
	backward bool
}

// Begin returns an iterator to the front element of the list.
func Begin[T any](l *list.List) Iterator[T] {
	return Iterator[T]{
		l: l,
		e: l.Front(),
	}
}

// End returns an iterator to the passed last element of the list.
func End[T any](l *list.List) Iterator[T] {
	return Iterator[T]{
		l: l,
	}
}

// RBegin returns an iterator to the back element of the list.
func RBegin[T any](l *list.List) Iterator[T] {
	return Iterator[T]{
		l:        l,
		e:        l.Back(),
		backward: true,
	}
}

// REnd returns an iterator to the passed first element of the list.
func REnd[T any](l *list.List) Iterator[T] {
	return Iterator[T]{
		l:        l,
		backward: true,
	}
}

func (l Iterator[T]) Eq(x Iterator[T]) bool {
	return l.e == x.e
}

func (l Iterator[T]) AllowMultiplePass() {}

func (l Iterator[T]) Next() Iterator[T] {
	var e *list.Element
	if l.backward {
		e = l.e.Prev()
	} else {
		e = l.e.Next()
	}
	return Iterator[T]{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l Iterator[T]) Prev() Iterator[T] {
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
	return Iterator[T]{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l Iterator[T]) Read() T {
	return l.e.Value.(T)
}

func (l Iterator[T]) Write(x T) {
	l.e.Value = x
}

// ListBackInserter returns an OutputIter to insert elements to the back of the
// list.
func ListBackInserter[T any](l *list.List) BackInserter[T] {
	return BackInserter[T]{l: l}
}

// BackInserter is an output iterator that appends values to a list.List.
type BackInserter[T any] struct {
	l *list.List
}

func (li BackInserter[T]) Write(x T) {
	li.l.PushBack(x)
}

// ListInserter returns an OutputIter to insert elements before a node.
func ListInserter[T any](l *list.List, e *list.Element) Inserter[T] {
	return Inserter[T]{l: l, e: e}
}

// Inserter is an output iterator that inserts values before a list element.
type Inserter[T any] struct {
	l *list.List
	e *list.Element
}

func (li Inserter[T]) Write(x T) {
	li.l.InsertBefore(x, li.e)
}
