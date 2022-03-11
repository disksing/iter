package iter

import (
	"math/rand"
	"strings"
)

type iotaReader[T Numeric] struct {
	x T
}

func (r iotaReader[T]) Read() T {
	return r.x
}

func (r iotaReader[T]) Next() iotaReader[T] {
	return iotaReader[T]{x: r.x + 1}
}

func (r iotaReader[T]) Eq(iotaReader[T]) bool {
	return false
}

// IotaReader creates an InputIter that returns [x, x+1, x+2...).
func IotaReader[T Numeric, It iotaReader[T]](x T) It {
	return It{x: x}
}

// IotaGenerator creates a Generator that returns [x, x+1, x+2...).
func IotaGenerator[T Numeric](x T) func() T {
	r := IotaReader(x)
	return func() T {
		v := r.Read()
		r = r.Next()
		return v
	}
}

type repeatReader[T any] struct {
	x T
}

func (r repeatReader[T]) Read() T { return r.x }

func (r repeatReader[T]) Next() repeatReader[T] { return r }

func (r repeatReader[T]) Eq(repeatReader[T]) bool { return false }

// RepeatReader creates an InputIter that returns [x, x, x...).
func RepeatReader[T any](x T) repeatReader[T] {
	return repeatReader[T]{x: x}
}

// RepeatGenerator creates an Generator that returns [x, x, x...).
func RepeatGenerator[T any](x T) func() T {
	return func() T { return x }
}

// RandomGenerator creates a generator that returns random item of a slice.
func RandomGenerator[T any](s []T, r *rand.Rand) func() T {
	return func() T { return s[r.Intn(len(s))] }
}

// MakeString creates a string by range spesified by [first, last). The value
// type should be byte or rune.
func MakeString[T byte | rune, It ForwardReader[T, It]](first, last It) string {
	var s strings.Builder
	for ; !first.Eq(last); first = first.Next() {
		switch v := any(first.Read()).(type) {
		case byte:
			s.WriteByte(v)
		case rune:
			s.WriteRune(v)
		}
	}
	return s.String()
}
