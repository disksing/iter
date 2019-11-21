package iter

// Equalable values can evaluate equivalent with anthoer value with the same
// type.
type Equalable interface {
	Equal(Any) bool
}

func _eq(x Any, y Any) bool {
	if e, ok := x.(Equalable); ok {
		return e.Equal(y)
	}
	if e, ok := y.(Equalable); ok {
		return e.Equal(x)
	}
	return x == y
}

func _ne(x Any, y Any) bool {
	return !_eq(x, y)
}
