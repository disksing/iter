package iter

import "container/list"

type ListIter struct {
	l *list.List
	e *list.Element
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
		l: l,
		e: l.Back(),
	}
}

func ListREnd(l *list.List) *ListIter {
	return &ListIter{
		l: l,
	}
}

func (l *ListIter) Equal(x Any) bool {
	return l.e == x.(*ListIter).e
}

func (l *ListIter) Next() ForwardIter {
	return &ListIter{
		l: l.l,
		e: l.e.Next(),
	}
}

func (l *ListIter) Prev() BackwardIter {
	return &ListIter{
		l: l.l,
		e: l.e.Prev(),
	}
}

func (l *ListIter) Read() Any {
	return l.e.Value
}

func (l *ListIter) Write(x Any) {
	l.e.Value = x
}

func ListBackInserter(l *list.List) ForwardWriter {
	return &listBackInserter{l: l}
}

type listBackInserter struct {
	l *list.List
}

func (li *listBackInserter) Next() ForwardIter {
	return li
}

func (li *listBackInserter) Write(x Any) {
	li.l.PushBack(x)
}

type listInserter struct {
	l *list.List
	e *list.Element
}

func (li *listInserter) Next() ForwardIter {
	return li
}

func (li *listInserter) Write(x Any) {
	li.l.InsertBefore(x, li.e)
}
