package iter

import "container/list"

type ListIter struct {
	l        *list.List
	e        *list.Element
	backward bool
}

func ListBegin(l *list.List) *ListIter {
	return &ListIter{
		l: l,
		e: l.Front(),
	}
}

func ListEnd(l *list.List) *ListIter {
	return &ListIter{
		l: l,
	}
}

func ListRBegin(l *list.List) *ListIter {
	return &ListIter{
		l:        l,
		e:        l.Back(),
		backward: true,
	}
}

func ListREnd(l *list.List) *ListIter {
	return &ListIter{
		l:        l,
		backward: true,
	}
}

func (l *ListIter) Eq(x Iter) bool {
	return l.e == x.(*ListIter).e
}

func (l *ListIter) AllowMultiplePass() {}

func (l *ListIter) Next() Incrementable {
	var e *list.Element
	if l.backward {
		e = l.e.Prev()
	} else {
		e = l.e.Next()
	}
	return &ListIter{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l *ListIter) Prev() BidiIter {
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
	return &ListIter{
		l:        l.l,
		e:        e,
		backward: l.backward,
	}
}

func (l *ListIter) Read() Any {
	return l.e.Value
}

func (l *ListIter) Write(x Any) {
	l.e.Value = x
}

func ListBackInserter(l *list.List) OutputIter {
	return &listBackInserter{l: l}
}

type listBackInserter struct {
	l *list.List
}

func (li *listBackInserter) Write(x Any) {
	li.l.PushBack(x)
}

type listInserter struct {
	l *list.List
	e *list.Element
}

func (li *listInserter) Write(x Any) {
	li.l.InsertBefore(x, li.e)
}
