package iter

// readonly algorithms.

func AllOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _eq(FindIfNot(first, last, pred), last)
}

func AnyOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _ne(FindIf(first, last, pred), last)
}

func NoneOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _eq(FindIf(first, last, pred), last)
}

func ForEach(first, last ForwardReader, f IterFunction) IterFunction {
	for ; _ne(first, last); first = NextReader(first) {
		f(first)
	}
	return f
}

func ForEachN(first ForwardReader, n int, f IterFunction) IterFunction {
	for ; n > 0; n, first = n-1, NextReader(first) {
		f(first)
	}
	return f
}

func Count(first, last ForwardReader, v Any) int {
	return CountIf(first, last, _eq1(v))
}

func CountIf(first, last ForwardReader, pred UnaryPredicate) int {
	var ret int
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			ret++
		}
	}
	return ret
}

func Mismatch(first1, last1, first2, last2 ForwardReader) (ForwardReader, ForwardReader) {
	return MismatchIf(first1, last1, first2, last2, _eq)
}

func MismatchIf(first1, last1, first2, last2 ForwardReader, pred BinaryPredicate) (ForwardReader, ForwardReader) {
	for _ne(first1, last1) && _ne(first2, last2) && pred(first1.Read(), first2.Read()) {
		first1, first2 = NextReader(first1), NextReader(first2)
	}
	return first1, first2
}

func Find(first, last ForwardReader, v Any) ForwardReader {
	return FindIf(first, last, _eq1(v))
}

func FindIf(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			return first
		}
	}
	return last
}

func FindIfNot(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	for ; _ne(first, last); first = NextReader(first) {
		if !pred(first.Read()) {
			return first
		}
	}
	return last
}

func FindEnd(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindEndIf(first, last, sFirst, sLast, _eq)
}

func FindEndIf(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
	if _eq(sFirst, sLast) {
		return last
	}
	result := last
	for {
		if newResult := SearchIf(first, last, sFirst, sLast, pred); _eq(newResult, last) {
			break
		} else {
			result = newResult
			first = NextReader(result)
		}
	}
	return result
}

func FindFirstOf(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindFirstOfIf(first, last, sFirst, sLast, _eq)
}

func FindFirstOfIf(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
	for ; _ne(first, last); first = NextReader(first) {
		for it := sFirst; _ne(it, sLast); it = NextReader(it) {
			if pred(first.Read(), it.Read()) {
				return first
			}
		}
	}
	return last
}

func AdjacentFind(first, last ForwardReader) ForwardReader {
	return AdjacentFindIf(first, last, _eq)
}

func AdjacentFindIf(first, last ForwardReader, pred BinaryPredicate) ForwardReader {
	if _eq(first, last) {
		return last
	}
	for next := NextReader(first); _ne(next, last); first, next = NextReader(first), NextReader(next) {
		if pred(first.Read(), next.Read()) {
			return first
		}
	}
	return last
}

func Search(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return SearchIf(first, last, sFirst, sLast, _eq)
}

func SearchIf(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
	for {
		it := first
		for sIt := sFirst; ; sIt, it = NextReader(sIt), NextReader(it) {
			if _eq(sIt, sLast) {
				return first
			}
			if _eq(it, last) {
				return last
			}
			if !pred(it.Read(), sIt.Read()) {
				break
			}
		}
		first = NextReader(first)
	}
}

func SearchN(first, last ForwardReader, count int, v Any) ForwardReader {
	return SearchNIf(first, last, count, v, _eq)
}

func SearchNIf(first, last ForwardReader, count int, v Any, pred BinaryPredicate) ForwardReader {
	if count <= 0 {
		return first
	}
	for ; _ne(first, last); first = NextReader(first) {
		if !pred(first, v) {
			continue
		}
		candidate := first
		var curCount int
		for {
			curCount++
			if curCount >= count {
				return candidate
			}
			if first = NextReader(first); _eq(first, last) {
				return last
			}
			if !pred(first.Read(), v) {
				break
			}
		}
	}
	return last
}
