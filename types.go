package iter

// Function types.
type (
	UnaryPredicate   func(Any) bool
	EqComparer       func(Any, Any) bool
	LessComparer     func(Any, Any) bool
	ThreeWayComparer func(Any, Any) int
	IteratorFunction func(Iter)
	UnaryOperation   func(Any) Any
	BinaryOperation  func(Any, Any) Any
	Generator        func() Any
)

// Returns a Predicate that returns true if the value equals v.
func _eq1(v Any) UnaryPredicate {
	return func(x Any) bool {
		return _eq(v, x)
	}
}

// Returns a Predicate that returns false if the value equals v.
func _ne1(v Any) UnaryPredicate {
	return func(x Any) bool {
		return _ne(v, x)
	}
}

func _not1(p UnaryPredicate) UnaryPredicate {
	return func(x Any) bool { return !p(x) }
}

func _true1(Any) bool { return true }

func _notrans(x Any) Any { return x }
