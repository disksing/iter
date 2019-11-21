package iter

func Max(a, b interface{}) interface{} {
	return MaxBy(a, b, _less)
}

func MaxBy(a, b interface{}, less BinaryPredicate) interface{} {
	if less(a, b) {
		return b
	}
	return a
}

func MaxElement(first, last ForwardReader) ForwardReader {
	return MaxElementBy(first, last, _less)
}

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

func Min(a, b interface{}) interface{} {
	return MinBy(a, b, _less)
}

func MinBy(a, b interface{}, less BinaryPredicate) interface{} {
	if less(a, b) {
		return a
	}
	return b
}

func MinElement(first, last ForwardReader) ForwardReader {
	return MinElementBy(first, last, _less)
}

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

func Minmax(a, b interface{}) (interface{}, interface{}) {
	return MinmaxBy(a, b, _less)
}

func MinmaxBy(a, b interface{}, less BinaryPredicate) (interface{}, interface{}) {
	if less(b, a) {
		return b, a
	}
	return a, b
}

func MinmaxElement(first, last ForwardReader) (ForwardReader, ForwardReader) {
	return MinmaxElementBy(first, last, _less)
}

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

func Clamp(v, lo, hi interface{}) interface{} {
	return ClampBy(v, lo, hi, _less)
}

func ClampBy(v, lo, hi interface{}, less BinaryPredicate) interface{} {
	if less(v, lo) {
		return lo
	}
	if less(hi, v) {
		return hi
	}
	return v
}
