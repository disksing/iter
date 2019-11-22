package iter

// Max returns the greater of the given values.
func Max(a, b Any) Any {
	return MaxBy(a, b, _less)
}

// MaxBy returns the greater of the given values. Values are compared using the
// given binary comparison function less.
func MaxBy(a, b Any, less BinaryPredicate) Any {
	if less(a, b) {
		return b
	}
	return a
}

// MaxElement returns the largest element in a range.
func MaxElement(first, last ForwardReader) ForwardReader {
	return MaxElementBy(first, last, _less)
}

// MaxElementBy returns the largest element in a range. Values are compared
// using the given binary comparison function less.
func MaxElementBy(first, last ForwardReader, less BinaryPredicate) ForwardReader {
	if _eq(first, last) {
		return last
	}
	max := first
	for first = NextReader(first); _ne(first, last); first = NextReader(first) {
		if less(max.Read(), first.Read()) {
			max = first
		}
	}
	return max
}

// Min returns the smaller of the given values.
func Min(a, b Any) Any {
	return MinBy(a, b, _less)
}

// MinBy returns the smaller of the given values. Values are compared using the
// given binary comparison function less.
func MinBy(a, b Any, less BinaryPredicate) Any {
	if less(a, b) {
		return a
	}
	return b
}

// MinElement returns the smallest element in a range.
func MinElement(first, last ForwardReader) ForwardReader {
	return MinElementBy(first, last, _less)
}

// MinElementBy returns the smallest element in a range. Values are compared
// using the given binary comparison function less.
func MinElementBy(first, last ForwardReader, less BinaryPredicate) ForwardReader {
	if _eq(first, last) {
		return last
	}
	min := first
	for first = NextReader(first); _ne(first, last); first = NextReader(first) {
		if less(first.Read(), min.Read()) {
			min = first
		}
	}
	return min
}

// Minmax returns the smaller and larger of two elements.
func Minmax(a, b Any) (Any, Any) {
	return MinmaxBy(a, b, _less)
}

// MinmaxBy returns the smaller and larger of two elements. Values are compared
// using the given binary comparison function less.
func MinmaxBy(a, b Any, less BinaryPredicate) (Any, Any) {
	if less(b, a) {
		return b, a
	}
	return a, b
}

// MinmaxElement returns the smallest and the largest elements in a range.
func MinmaxElement(first, last ForwardReader) (ForwardReader, ForwardReader) {
	return MinmaxElementBy(first, last, _less)
}

// MinmaxElementBy returns the smallest and the largest elements in a range.
// Values are compared using the given binary comparison function less.
func MinmaxElementBy(first, last ForwardReader, less BinaryPredicate) (ForwardReader, ForwardReader) {
	if _eq(first, last) {
		return first, first
	}
	min, max := first, first
	for first = NextReader(first); _ne(first, last); first = NextReader(first) {
		i := first
		first = NextReader(first)
		if _eq(first, last) {
			if less(i.Read(), min.Read()) {
				min = i
			} else if less(max.Read(), i.Read()) {
				max = i
			}
			break
		} else {
			if less(first.Read(), i.Read()) {
				if less(first.Read(), min.Read()) {
					min = first
				}
				if less(max.Read(), i.Read()) {
					max = i
				}
			} else {
				if less(i.Read(), min.Read()) {
					min = i
				}
				if less(max.Read(), first.Read()) {
					max = first
				}
			}
		}
	}
	return min, max
}

// Clamp clamps a value between a pair of boundary values.
func Clamp(v, lo, hi Any) Any {
	return ClampBy(v, lo, hi, _less)
}

// ClampBy clamps a value between a pair of boundary values. Values are compared
// using the given binary comparison function less.
func ClampBy(v, lo, hi Any, less BinaryPredicate) Any {
	if less(v, lo) {
		return lo
	}
	if less(hi, v) {
		return hi
	}
	return v
}
