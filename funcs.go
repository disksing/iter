package iter

// Any represents any type.
type Any = interface{}

// Function types.
type (
	UnaryPredicate  func(Any) bool
	BinaryPredicate func(Any, Any) bool
	IterFunction    func(Iter)
	UnaryOperation  func(Any) Any
	BinaryOperation func(Any, Any) Any
	Generator       func() Any
)

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
