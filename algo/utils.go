package algo

import (
	. "github.com/disksing/iter/v2"
)

func __iter_eq[It Comparable[It]](x, y It) bool {
	return x.Eq(y)
}

func __write_next[T any, It OutputIter[T]](out It, v T) It {
	out.Write(v)
	if inc, ok := any(out).(ForwardMovable[It]); ok {
		out = inc.Next()
	}
	return out
}

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

func __cmp[T Ordered](x, y T) int {
	if x > y {
		return 1
	}
	if x < y {
		return -1
	}
	return 0
}

func __eq[T comparable](x, y T) bool {
	return x == y
}

func __inc[T Integer](x T) T {
	return x + 1
}

func __add[T1, T2 Numeric](x, y T1) T2 {
	return T2(x + y)
}

func __sub[T1, T2 Numeric](x, y T1) T2 {
	return T2(x - y)
}

func __mul[T1, T2 Numeric](x, y T1) T2 {
	return T2(x * y)
}

func __eq1[T comparable](v T) UnaryPredicate[T] {
	return func(x T) bool {
		return x == v
	}
}

func __eq2[T comparable](v1, v2 T) bool {
	return v1 == v2
}

func __eq_bind1[T1, T2 any](p EqComparer[T1, T2], v T1) UnaryPredicate[T2] {
	return func(x T2) bool {
		return p(v, x)
	}
}

func __not1[T any](p UnaryPredicate[T]) UnaryPredicate[T] {
	return func(x T) bool { return !p(x) }
}

func __less[T Ordered](v1, v2 T) bool {
	return v1 < v2
}

func __true1[T any](T) bool { return true }

func __noop[T any](x T) T { return x }
