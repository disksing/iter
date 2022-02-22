package iter

type (
	// UnaryPredicate checks if a value satisfy condition.
	UnaryPredicate func(any) bool
	// EqComparer checks if first value equals to the second value.
	EqComparer func(any, any) bool
	// LessComparer checks if first value is less than the second value.
	LessComparer func(any, any) bool
	// ThreeWayComparer compares 2 values, returns 1 if first>second, 0 if
	// first=second, -1 if first<second.
	ThreeWayComparer func(any, any) int
	// IteratorFunction apply some actions to a value.
	IteratorFunction func(any)
	// UnaryOperation transforms a value to another.
	UnaryOperation func(any) any
	// BinaryOperation transforms 2 values to 1 value.
	BinaryOperation func(any, any) any
	// Generator creates a value on each call.
	Generator func() any
)

func _eq(x, y any) bool {
	type ieq interface{ Eq(any) bool }
	if e, ok := x.(ieq); ok {
		return e.Eq(y)
	}
	return x == y
}

func _ne(x, y any) bool {
	return !_eq(x, y)
}

func _less(x, y any) bool {
	type iless interface{ Less(any) bool }
	if c, ok := x.(iless); ok {
		return c.Less(y)
	}
	return _cmp(x, y) < 0
}

func _cmp(x, y any) int {
	type icmp interface{ Cmp(any) int }
	if t, ok := x.(icmp); ok {
		return t.Cmp(y)
	}
	return reflectCmp(x, y)
}

func _inc(x any) any {
	type iinc interface{ Inc() any }
	if i, ok := x.(iinc); ok {
		return i.Inc()
	}
	return reflectInc(x)
}

func _add(x, y any) any {
	type iadd interface{ Add(any) any }
	if a, ok := x.(iadd); ok {
		return a.Add(y)
	}
	return reflectAdd(x, y)
}

func _sub(x, y any) any {
	type isub interface{ Sub(any) any }
	if s, ok := x.(isub); ok {
		return s.Sub(y)
	}
	return reflectSub(x, y)
}

func _mul(x, y any) any {
	type imul interface{ Mul(any) any }
	if m, ok := x.(imul); ok {
		return m.Mul(y)
	}
	return reflectMul(x, y)
}

// Returns a Predicate that returns true if the value equals v.
func _eq1(v any) UnaryPredicate {
	return func(x any) bool {
		return _eq(v, x)
	}
}

func _not1(p UnaryPredicate) UnaryPredicate {
	return func(x any) bool { return !p(x) }
}

func _true1(any) bool { return true }

func _noop(x any) any { return x }
