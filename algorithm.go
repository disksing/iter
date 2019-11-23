package iter

import "math/rand"

// AllOf checks if unary predicate pred returns true for all elements in the
// range [first, last).
func AllOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _eq(FindIfNot(first, last, pred), last)
}

// AnyOf checks if unary predicate pred returns true for at least one element in
// the range [first, last).
func AnyOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _ne(FindIf(first, last, pred), last)
}

// NoneOf checks if unary predicate pred returns true for no elements in the
// range [first, last).
func NoneOf(first, last ForwardReader, pred UnaryPredicate) bool {
	return _eq(FindIf(first, last, pred), last)
}

// ForEach applies the given function f to the result of dereferencing every
// iterator in the range [first, last), in order.
func ForEach(first, last ForwardReader, f IterFunction) IterFunction {
	for ; _ne(first, last); first = NextReader(first) {
		f(first)
	}
	return f
}

// ForEachN applies the given function f to the result of dereferencing every
// iterator in the range [first, first + n), in order.
func ForEachN(first ForwardReader, n int, f IterFunction) IterFunction {
	for ; n > 0; n, first = n-1, NextReader(first) {
		f(first)
	}
	return f
}

// Count counts the elements that are equal to value.
func Count(first, last ForwardReader, v Any) int {
	return CountIf(first, last, _eq1(v))
}

// CountIf counts elements for which predicate pred returns true.
func CountIf(first, last ForwardReader, pred UnaryPredicate) int {
	var ret int
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			ret++
		}
	}
	return ret
}

// Mismatch returns the first mismatching pair of elements from two ranges: one
// defined by [first1, last1) and another defined by [first2,last2). If last2 is
// nil, it denotes first2 + (last1 - first1).
func Mismatch(first1, last1, first2, last2 ForwardReader) (ForwardReader, ForwardReader) {
	return MismatchBy(first1, last1, first2, last2, _eq)
}

// MismatchBy returns the first mismatching pair of elements from two ranges:
// one defined by [first1, last1) and another defined by [first2,last2). If
// last2 is nil, it denotes first2 + (last1 - first1). Elements are compared
// using the given binary predicate pred.
func MismatchBy(first1, last1, first2, last2 ForwardReader, pred BinaryPredicate) (ForwardReader, ForwardReader) {
	for _ne(first1, last1) && (last2 == nil || _ne(first2, last2)) && pred(first1.Read(), first2.Read()) {
		first1, first2 = NextReader(first1), NextReader(first2)
	}
	return first1, first2
}

// Find returns the first element in the range [first, last) that is equal to
// value.
func Find(first, last ForwardReader, v Any) ForwardReader {
	return FindIf(first, last, _eq1(v))
}

// FindIf returns the first element in the range [first, last) which predicate
// pred returns true.
func FindIf(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			return first
		}
	}
	return last
}

// FindIfNot returns the first element in the range [first, last) which
// predicate pred returns false.
func FindIfNot(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	return FindIf(first, last, _not1(pred))
}

// FindEnd searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last). If [sFirst, sLast) is empty or such sequence is
// found, last is returned.
func FindEnd(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindEndBy(first, last, sFirst, sLast, _eq)
}

// FindEnd searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last). If [sFirst, sLast) is empty or such sequence is
// found, last is returned. Elements are compared using the given binary
// predicate pred.
func FindEndBy(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
	if _eq(sFirst, sLast) {
		return last
	}
	result := last
	for {
		if newResult := SearchBy(first, last, sFirst, sLast, pred); _eq(newResult, last) {
			break
		} else {
			result = newResult
			first = NextReader(result)
		}
	}
	return result
}

// FindFirstOf searches the range [first, last) for any of the elements in the
// range [sFirst, sLast). Elements are compared using the given binary predicate
// pred.
func FindFirstOf(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindFirstOfBy(first, last, sFirst, sLast, _eq)
}

// FindFirstOfBy searches the range [first, last) for any of the elements in the
// range [sFirst, sLast).
func FindFirstOfBy(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
	return FindIf(first, last, func(x Any) bool {
		return AnyOf(sFirst, sLast, func(s Any) bool {
			return pred(x, s)
		})
	})
}

// AdjacentFind searches the range [first, last) for two consecutive identical
// elements.
func AdjacentFind(first, last ForwardReader) ForwardReader {
	return AdjacentFindBy(first, last, _eq)
}

// AdjacentFindBy searches the range [first, last) for two consecutive identical
// elements. Elements are compared using the given binary predicate pred.
func AdjacentFindBy(first, last ForwardReader, pred BinaryPredicate) ForwardReader {
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

// Search searches for the first occurrence of the sequence of elements
// [sFirst, sLast) in the range [first, last).
func Search(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return SearchBy(first, last, sFirst, sLast, _eq)
}

// SearchBy searches for the first occurrence of the sequence of elements
// [sFirst, sLast) in the range [first, last). Elements are compared using the
// given binary predicate pred.
func SearchBy(first, last, sFirst, sLast ForwardReader, pred BinaryPredicate) ForwardReader {
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

// SearchN searches the range [first, last) for the first sequence of count
// identical elements, each equal to the given value.
func SearchN(first, last ForwardReader, count int, v Any) ForwardReader {
	return SearchNBy(first, last, count, v, _eq)
}

// SearchNBy searches the range [first, last) for the first sequence of count
// identical elements. Elements are compared using the given binary predicate
// pred.
func SearchNBy(first, last ForwardReader, count int, v Any, pred BinaryPredicate) ForwardReader {
	if count <= 0 {
		return first
	}
	for ; _ne(first, last); first = NextReader(first) {
		if !pred(first.Read(), v) {
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

func Copy(first, last ForwardReader, dFirst ForwardWriter) ForwardWriter {
	return CopyIf(first, last, dFirst, _true1)
}

func CopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate) ForwardWriter {
	for _ne(first, last) {
		if pred(first.Read()) {
			dFirst.Write(first.Read())
			dFirst = NextWriter(dFirst)
		}
		first = NextReader(first)
	}
	return dFirst
}

func CopyN(first ForwardReader, count int, result ForwardWriter) ForwardWriter {
	for ; count > 0; count-- {
		result.Write(first.Read())
		first, result = NextReader(first), NextWriter(result)
	}
	return result
}

func CopyBackward(first, last BackwardReader, dLast BackwardWriter) BackwardWriter {
	for _ne(first, last) {
		last, dLast = PrevReader(last), PrevWriter(dLast)
		dLast.Write(last.Read())
	}
	return dLast
}

func Fill(first, last ForwardWriter, v Any) {
	for ; _ne(first, last); first = NextWriter(first) {
		first.Write(v)
	}
}

func FillN(first ForwardWriter, count int, v Any) {
	for ; count > 0; count-- {
		first.Write(v)
		first = NextWriter(first)
	}
}

func Transform(first, last ForwardReader, dFirst ForwardWriter, op UnaryOperation) ForwardWriter {
	for ; _ne(first, last); dFirst, first = NextWriter(dFirst), NextReader(first) {
		dFirst.Write(op(first.Read()))
	}
	return dFirst
}

func Transform2(first1, last1, first2 ForwardReader, dFirst ForwardWriter, op BinaryOperation) ForwardWriter {
	for ; _ne(first1, last1); dFirst, first1, first2 = NextWriter(dFirst), NextReader(first1), NextReader(first2) {
		dFirst.Write(op(first1.Read(), first2.Read()))
	}
	return dFirst
}

func Generate(first, last ForwardWriter, g Generator) {
	for ; _ne(first, last); first = NextWriter(first) {
		first.Write(g())
	}
}

func GenerateN(first ForwardWriter, count int, g Generator) ForwardWriter {
	for ; count > 0; count-- {
		first.Write(g())
		first = NextWriter(first)
	}
	return first
}

func Remove(first, last ForwardReadWriter, v Any) ForwardReadWriter {
	return RemoveIf(first, last, _eq1(v))
}

func RemoveIf(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	first = FindIf(first, last, pred).(ForwardReadWriter)
	if _ne(first, last) {
		for i := NextReadWriter(first); _ne(i, last); i = NextReadWriter(i) {
			if !pred(i.Read()) {
				first.Write(i.Read())
				first = NextReadWriter(first)
			}
		}
	}
	return first
}

func RemoveCopy(first, last ForwardReader, dFirst ForwardWriter, v Any) ForwardWriter {
	return RemoveCopyIf(first, last, dFirst, _eq1(v))
}

func RemoveCopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate) ForwardWriter {
	for ; _ne(first, last); first = NextReader(first) {
		if !pred(first.Read()) {
			dFirst.Write(first.Read())
			dFirst = NextWriter(dFirst)
		}
	}
	return dFirst
}

func Replace(first, last ForwardReadWriter, old, new Any) {
	ReplaceIf(first, last, _eq1(old), new)
}

func ReplaceIf(first, last ForwardReadWriter, pred UnaryPredicate, v Any) {
	for ; _ne(first, last); first = NextReadWriter(first) {
		if pred(first.Read()) {
			first.Write(v)
		}
	}
}

func ReplaceCopy(first, last ForwardReader, dFirst ForwardWriter, old, new Any) ForwardWriter {
	return ReplaceCopyIf(first, last, dFirst, _eq1(old), new)
}

func ReplaceCopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate, v Any) ForwardWriter {
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			dFirst.Write(first.Read())
			dFirst = NextWriter(dFirst)
		}
	}
	return dFirst
}

func IterSwap(a, b ReadWriter) {
	t := a.Read()
	a.Write(b.Read())
	b.Write(t)
}

func SwapRanges(first1, last1, first2 ForwardReadWriter) {
	for ; _ne(first1, last1); first1, first2 = NextReadWriter(first1), NextReadWriter(first2) {
		IterSwap(first1, first2)
	}
}

func Reverse(first, last BidiReadWriter) {
	for ; _ne(first, last); first = NextBidiReadWriter(first) {
		last = PrevBidiReadWriter(last)
		if _eq(first, last) {
			return
		}
		IterSwap(first, last)
	}
}

func ReverseCopy(first, last BackwardReader, dFirst ForwardWriter) ForwardWriter {
	for _ne(first, last) {
		last = PrevReader(last)
		dFirst.Write(last.Read())
		dFirst = NextWriter(dFirst)
	}
	return dFirst
}

func Rotate(first, nFirst, last ForwardReadWriter) ForwardReadWriter {
	if _eq(first, nFirst) {
		return last
	}
	if _eq(nFirst, last) {
		return first
	}
	read, write, nextRead := nFirst, first, first
	for _ne(read, last) {
		if _eq(write, nextRead) {
			nextRead = read
		}
		IterSwap(write, read)
		write, read = NextReadWriter(write), NextReadWriter(read)
	}
	Rotate(write, nextRead, last)
	return write
}

func RotateCopy(first, nFirst, last ForwardReader, dFirst ForwardWriter) ForwardWriter {
	dFirst = Copy(nFirst, last, dFirst)
	return Copy(first, nFirst, dFirst)
}

func Shuffle(first, last RandomReadWriter, r *rand.Rand) {
	for n := first.Distance(last) - 1; n > 0; n-- {
		IterSwap(AdvanceNReadWriter(first, n), AdvanceNReadWriter(first, r.Intn(n+1)))
	}
}

func Sample(first, last ForwardReader, out ForwardWriter, n int, r *rand.Rand) ForwardWriter {
	_, okr := first.(RandomReader)
	randWriter, okw := out.(RandomWriter)
	if okr || !okw {
		return _selectionSample(first, last, out, n, r)
	}
	return _reservoirSample(first, last, randWriter, n, r)
}

func _selectionSample(first, last ForwardReader, out ForwardWriter, n int, r *rand.Rand) ForwardWriter {
	unsampled := Distance(first, last)
	if n > unsampled {
		n = unsampled
	}
	for ; n != 0; first = NextReader(first) {
		unsampled--
		if r.Intn(unsampled) < n {
			out.Write(first.Read())
			out, n = NextWriter(out), n-1
		}
	}
	return out
}

func _reservoirSample(first, last ForwardReader, out RandomWriter, n int, r *rand.Rand) RandomWriter {
	var k int
	for ; _ne(first, last) && k < n; first, k = NextReader(first), k+1 {
		AdvanceNWriter(out, k).Write(first.Read())
	}
	if _eq(first, last) {
		return AdvanceNWriter(out, k)
	}
	sz := k
	for ; _ne(first, last); first, k = NextReader(first), k+1 {
		if r.Intn(k) < sz {
			AdvanceNWriter(out, k).Write(first.Read())
		}
	}
	return AdvanceNWriter(out, n)
}

func Unique(first, last ForwardReadWriter) ForwardReadWriter {
	return UniqueIf(first, last, _eq)
}

func UniqueIf(first, last ForwardReadWriter, pred BinaryPredicate) ForwardReadWriter {
	if _eq(first, last) {
		return last
	}
	result := first
	for {
		first = NextReadWriter(first)
		if _eq(first, last) {
			return NextReadWriter(result)
		}
		if !pred(result.Read(), first.Read()) {
			if result = NextReadWriter(result); _ne(result, first) {
				result.Write(first.Read())
			}
		}
	}
}

func UniqueCopy(first, last ForwardReader, result ForwardWriter) ForwardWriter {
	return UniqueCopyIf(first, last, result, _eq)
}

func UniqueCopyIf(first, last ForwardReader, result ForwardWriter, pred BinaryPredicate) ForwardWriter {
	if _ne(first, last) {
		v := first.Read()
		result.Write(v)
		result = NextWriter(result)
		for first = NextReader(first); _ne(first, last); first = NextReader(first) {
			if !pred(v, first.Read()) {
				v = first.Read()
				result.Write(v)
				result = NextWriter(result)
			}
		}
	}
	return result
}

func IsPartitioned(first, last ForwardReader, pred UnaryPredicate) bool {
	return NoneOf(FindIfNot(first, last, pred), last, pred)
}

func Partition(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	first = FindIfNot(first, last, pred).(ForwardReadWriter)
	if _eq(first, last) {
		return first
	}
	for i := NextReadWriter(first); _ne(i, last); i = NextReadWriter(i) {
		if pred(i.Read()) {
			IterSwap(first, i)
			first = NextReadWriter(first)
		}
	}
	return first
}

func ParittionCopy(first, last ForwardReader, outTrue, outFalse ForwardWriter, pred UnaryPredicate) (ForwardWriter, ForwardWriter) {
	for ; _ne(first, last); first = NextReader(first) {
		if pred(first.Read()) {
			outTrue.Write(first.Read())
			outTrue = NextWriter(outTrue)
		} else {
			outFalse.Write(first.Read())
			outFalse = NextWriter(outFalse)
		}
	}
	return outTrue, outFalse
}

func StablePartition(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	panic("TODO: not implemented")
}

func PartitionPoint(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	l := Distance(first, last)
	for l != 0 {
		l2 := l / 2
		m := AdvanceN(first, l2).(ForwardReader)
		if pred(m.Read()) {
			first = NextReader(m)
			l -= l2 + 1
		} else {
			l = l2
		}
	}
	return first
}

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

func Equal(first1, last1, first2, last2 ForwardReader) bool {
	return EqualIf(first1, last1, first2, last2, _eq)
}

func EqualIf(first1, last1, first2, last2 ForwardReader, pred BinaryPredicate) bool {
	for ; _ne(first1, last1); first1, first2 = NextReader(first1), NextReader(first2) {
		if !pred(first1.Read(), first2.Read()) {
			return false
		}
	}
	return true
}

// func IsPermutation(first1, last1, first2, last2 ForwardIter) bool {
// 	return IsPermutationIf(first1, last1, first2, last2, eqv2)
// }

// func IsPermutationIf(first1, last1, first2, last2 ForwardIter, eq func(Iter, Iter) bool) bool {
// 	for ; ne(first1, last1); first1, first2 = first1.Next(), first2.Next() {
// 		if !eq(first1, first2) {
// 			break
// 		}
// 	}
// 	if first1 == last1 {
// 		return true
// 	}

// 	l1 := Distance(first1, last1)
// 	if l1 == 1 {
// 		return false
// 	}
// 	last2 = AdvanceN(last2, l1).(ForwardIter)
// 	for i := first1; i != last1; i = i.Next() {
// 		match := first1
// 		for ; match != i; match = match.Next() {
// 			if eq(match, i) {
// 				break
// 			}
// 		}
// 		if match == i {
// 			var c2 int
// 			for j := first2; j != last2; j = j.Next() {
// 				if eq(i, j) {
// 					c2++
// 				}
// 			}
// 			if c2 == 0 {
// 				return false
// 			}
// 			c1 := 1
// 			for j := Advance(i).(ForwardIter); j != last1; j = j.Next() {
// 				if eq(i, j) {
// 					c1++
// 				}
// 			}
// 			if c1 != c2 {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }
