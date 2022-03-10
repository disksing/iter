package iter

type (
	// UnaryPredicate checks if a value satisfy condition.
	UnaryPredicate[T any] func(T) bool
	// EqComparer checks if first value equals to the second value.
	EqComparer[T1, T2 any] func(T1, T2) bool
	// LessComparer checks if first value is less than the second value.
	LessComparer[T any] func(T, T) bool
	// ThreeWayComparer compares 2 values, returns 1 if first>second, 0 if
	// first=second, -1 if first<second.
	ThreeWayComparer[T1, T2 any] func(T1, T2) int
	// IteratorFunction apply some actions to a value.
	IteratorFunction[T any] func(T)
	// UnaryOperation transforms a value to another.
	UnaryOperation[T1, T2 any] func(T1) T2
	// BinaryOperation transforms 2 values to 1 value.
	BinaryOperation[T1, T2, T3 any] func(T1, T2) T3
	// Generator creates a value on each call.
	Generator[T any] func() T
)

func _cmp[T Ordered](x, y T) int {
	if x > y {
		return 1
	}
	if x < y {
		return -1
	}
	return 0
}

func _eq[T comparable](x, y T) bool {
	return x == y
}

func _inc[T Integer](x T) T {
	return x + 1
}

func _add[T1, T2 Numeric](x, y T1) T2 {
	return T2(x + y)
}

func _sub[T1, T2 Numeric](x, y T1) T2 {
	return T2(x - y)
}

func _mul[T1, T2 Numeric](x, y T1) T2 {
	return T2(x * y)
}

// Returns a Predicate that returns true if the value equals v.
func _eq1[T comparable](v T) UnaryPredicate[T] {
	return func(x T) bool {
		return x == v
	}
}

func _eq2[T comparable](v1, v2 T) bool {
	return v1 == v2
}

func _eq_bind1[T1, T2 any](p EqComparer[T1, T2], v T1) UnaryPredicate[T2] {
	return func(x T2) bool {
		return p(v, x)
	}
}

func _eq_bind2[T1, T2 any](p EqComparer[T1, T2], v T2) UnaryPredicate[T1] {
	return func(x T1) bool {
		return p(x, v)
	}
}

func _not1[T any](p UnaryPredicate[T]) UnaryPredicate[T] {
	return func(x T) bool { return !p(x) }
}

func _less[T Ordered](v1, v2 T) bool {
	return v1 < v2
}

func _true1[T any](T) bool { return true }

func _noop[T any](x T) T { return x }
