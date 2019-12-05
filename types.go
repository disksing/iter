package iter

// Any represents any type.
type Any interface{}

type (
	// UnaryPredicate checks if a value satisfy condition.
	UnaryPredicate func(Any) bool
	// EqComparer checks if first value equals to the second value.
	EqComparer func(Any, Any) bool
	// LessComparer checks if first value is less than the second value.
	LessComparer func(Any, Any) bool
	// ThreeWayComparer compares 2 values, returns 1 if first>second, 0 if
	// first=second, -1 if first<second.
	ThreeWayComparer func(Any, Any) int
	// IteratorFunction apply some actions to a value.
	IteratorFunction func(Any)
	// UnaryOperation transforms a value to another.
	UnaryOperation func(Any) Any
	// BinaryOperation trasnfoms 2 values to 1 value.
	BinaryOperation func(Any, Any) Any
	// Generator creates a value on each call.
	Generator func() Any
)

func _eq(x, y Any) bool {
	type ieq interface{ Eq(Any) bool }
	if e, ok := x.(ieq); ok {
		return e.Eq(y)
	}
	return x == y
}

func _ne(x, y Any) bool {
	return !_eq(x, y)
}

func _less(x, y Any) bool {
	type iless interface{ Less(Any) bool }
	if c, ok := x.(iless); ok {
		return c.Less(y)
	}
	return _cmp(x, y) == -1
}

func _cmp(x, y Any) int {
	type icmp interface{ Cmp(Any) int }
	if t, ok := x.(icmp); ok {
		return t.Cmp(y)
	}
	return reflectCompare(x, y)
}

func _inc(x Any) Any {
	type iinc interface{ Inc() Any }
	if i, ok := x.(iinc); ok {
		return i.Inc()
	}
	return reflectInc(x)
}

func _add(x, y Any) Any {
	type iadd interface{ Add(Any) Any }
	if a, ok := x.(iadd); ok {
		return a.Add(y)
	}
	return reflectAdd(x, y)
}

func _sub(x, y Any) Any {
	type isub interface{ Sub(Any) Any }
	if s, ok := x.(isub); ok {
		return s.Sub(y)
	}
	return reflectSub(x, y)
}

func _mul(x, y Any) Any {
	type imul interface{ Mul(Any) Any }
	if m, ok := x.(imul); ok {
		return m.Mul(y)
	}
	return reflectMul(x, y)
}

// Returns a Predicate that returns true if the value equals v.
func _eq1(v Any) UnaryPredicate {
	return func(x Any) bool {
		return _eq(v, x)
	}
}

func _not1(p UnaryPredicate) UnaryPredicate {
	return func(x Any) bool { return !p(x) }
}

func _true1(Any) bool { return true }

func _notrans(x Any) Any { return x }
