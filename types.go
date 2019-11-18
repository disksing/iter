package iter

type Any = interface{}

type Equaler interface {
	Equal(Any) bool
}

// type Copier interface {
// 	Copy() Any
// }

// func _cp(x Any) Any {
// 	if c, ok := x.(Copier); ok {
// 		return c.Copy()
// 	}
// 	return x
// }

type (
	UnaryPredicate  func(Any) bool
	BinaryPredicate func(Any, Any) bool
	IterFunction    func(Iter)
	UnaryOperation  func(Any) Any
	BinaryOperation func(Any, Any) Any
	Generator       func() Any
)

func _eq(x Any, y Any) bool {
	if e, ok := x.(Equaler); ok {
		return e.Equal(y)
	}
	if e, ok := y.(Equaler); ok {
		return e.Equal(x)
	}
	return x == y
}

func _ne(x Any, y Any) bool {
	return !_eq(x, y)
}

func _eq1(v Any) UnaryPredicate {
	return func(x Any) bool {
		return _eq(v, x)
	}
}

func _ne1(v Any) UnaryPredicate {
	return func(x Any) bool {
		return _ne(v, x)
	}
}

func _true1(Any) bool { return true }

type Comparable interface {
	Less(Any) bool
}

func _less(x, y Any) bool {
	if c, ok := x.(Comparable); ok {
		return c.Less(y)
	}
	return reflectCompare(x, y) == -1
}
