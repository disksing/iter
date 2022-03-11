package iter

import (
	"container/heap"
	"math/rand"
	"sort"
)

// AllOf checks if unary predicate pred returns true for all elements in the
// range [first, last).
func AllOf[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) bool {
	return __iter_eq(FindIfNot(first, last, pred), last)
}

// AnyOf checks if unary predicate pred returns true for at least one element in
// the range [first, last).
func AnyOf[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) bool {
	return !__iter_eq(FindIf(first, last, pred), last)
}

// NoneOf checks if unary predicate pred returns true for no elements in the
// range [first, last).
func NoneOf[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) bool {
	return __iter_eq(FindIf(first, last, pred), last)
}

// ForEach applies the given function f to the result of dereferencing every
// iterator in the range [first, last), in order.
func ForEach[T any, It InputIter[T, It]](first, last It, f IteratorFunction[T]) IteratorFunction[T] {
	for ; !__iter_eq(first, last); first = first.Next() {
		f(first.Read())
	}
	return f
}

// ForEachN applies the given function f to the result of dereferencing every
// iterator in the range [first, first + n), in order.
func ForEachN[T any, It InputIter[T, It]](first It, n int, f IteratorFunction[T]) IteratorFunction[T] {
	for ; n > 0; n, first = n-1, first.Next() {
		f(first.Read())
	}
	return f
}

// Count counts the elements that are equal to value.
func Count[T comparable, It InputIter[T, It]](first, last It, v T) int {
	return CountIf(first, last, _eq1(v))
}

// CountIf counts elements for which predicate pred returns true.
func CountIf[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) int {
	var ret int
	for ; !__iter_eq(first, last); first = first.Next() {
		if pred(first.Read()) {
			ret++
		}
	}
	return ret
}

// Mismatch returns the first mismatching pair of elements from two ranges: one
// defined by [first1, last1) and another defined by [first2,last2).
//
// If last2 is nil, it denotes first2 + (last1 - first1).
func Mismatch[T comparable, It1 InputIter[T, It1], It2 InputIter[T, It2]](first1, last1 It1, first2 It2, last2 *It2) (It1, It2) {
	return MismatchBy(first1, last1, first2, last2, _eq2[T])
}

// MismatchBy returns the first mismatching pair of elements from two ranges:
// one defined by [first1, last1) and another defined by [first2,last2).
//
// If last2 is nil, it denotes first2 + (last1 - first1). Elements are compared
// using the given comparer eq.
func MismatchBy[T1, T2 any, It1 InputIter[T1, It1], It2 InputIter[T2, It2]](first1, last1 It1, first2 It2, last2 *It2, eq EqComparer[T1, T2]) (It1, It2) {
	for !__iter_eq(first1, last1) && (last2 == nil || !__iter_eq(first2, *last2)) && eq(first1.Read(), first2.Read()) {
		first1, first2 = first1.Next(), first2.Next()
	}
	return first1, first2
}

// Find returns the first element in the range [first, last) that is equal to
// value.
func Find[T comparable, It InputIter[T, It]](first, last It, v T) It {
	return FindIf(first, last, _eq1(v))
}

// FindIf returns the first element in the range [first, last) which predicate
// pred returns true.
func FindIf[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	for ; !__iter_eq(first, last); first = first.Next() {
		if pred(first.Read()) {
			return first
		}
	}
	return last
}

// FindIfNot returns the first element in the range [first, last) which
// predicate pred returns false.
func FindIfNot[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	return FindIf(first, last, _not1(pred))
}

// FindEnd searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last).
//
// If [sFirst, sLast) is empty or such sequence is found, last is returned.
func FindEnd[T comparable, It ForwardReader[T, It]](first, last, sFirst, sLast It) It {
	return FindEndBy(first, last, sFirst, sLast, _eq2[T])
}

// FindEndBy searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last).
//
// If [sFirst, sLast) is empty or such sequence is found, last is returned.
// Elements are compared using the given binary comparer eq.
func FindEndBy[T1, T2 any, It1 ForwardReader[T1, It1], It2 ForwardReader[T2, It2]](first, last It1, sFirst, sLast It2, eq EqComparer[T1, T2]) It1 {
	if __iter_eq(sFirst, sLast) {
		return last
	}
	result := last
	for {
		if newResult := SearchBy(first, last, sFirst, sLast, eq); __iter_eq(newResult, last) {
			break
		} else {
			result = newResult
			first = result.Next()
		}
	}
	return result
}

// FindFirstOf searches the range [first, last) for any of the elements in the
// range [sFirst, sLast).
func FindFirstOf[T comparable, It ForwardReader[T, It]](first, last It, sFirst, sLast It) It {
	return FindFirstOfBy(first, last, sFirst, sLast, _eq2[T])
}

// FindFirstOfBy searches the range [first, last) for any of the elements in the
// range [sFirst, sLast).
//
// Elements are compared using the given binary comparer eq.
func FindFirstOfBy[T1, T2 any, It1 ForwardReader[T1, It1], It2 ForwardReader[T2, It2]](first, last It1, sFirst, sLast It2, eq EqComparer[T1, T2]) It1 {
	return FindIf(first, last, func(x T1) bool {
		return AnyOf(sFirst, sLast, func(s T2) bool {
			return eq(x, s)
		})
	})
}

// AdjacentFind searches the range [first, last) for two consecutive identical
// elements.
func AdjacentFind[T comparable, It ForwardReader[T, It]](first, last It) It {
	return AdjacentFindBy(first, last, _eq2[T])
}

// AdjacentFindBy searches the range [first, last) for two consecutive identical
// elements.
//
// Elements are compared using the given binary comparer eq.
func AdjacentFindBy[T any, It ForwardReader[T, It]](first, last It, eq EqComparer[T, T]) It {
	if __iter_eq(first, last) {
		return last
	}
	for next := first.Next(); !__iter_eq(next, last); first, next = first.Next(), next.Next() {
		if eq(first.Read(), next.Read()) {
			return first
		}
	}
	return last
}

// Search searches for the first occurrence of the sequence of elements
// [sFirst, sLast) in the range [first, last).
func Search[T comparable, It ForwardReader[T, It]](first, last, sFirst, sLast It) It {
	return SearchBy(first, last, sFirst, sLast, _eq2[T])
}

// SearchBy searches for the first occurrence of the sequence of elements
// [sFirst, sLast) in the range [first, last).
//
// Elements are compared using the given binary comparer eq.
func SearchBy[T1, T2 any, It1 ForwardReader[T1, It1], It2 ForwardReader[T2, It2]](first, last It1, sFirst, sLast It2, eq EqComparer[T1, T2]) It1 {
	for {
		it := first
		for sIt := sFirst; ; sIt, it = sIt.Next(), it.Next() {
			if __iter_eq(sIt, sLast) {
				return first
			}
			if __iter_eq(it, last) {
				return last
			}
			if !eq(it.Read(), sIt.Read()) {
				break
			}
		}
		first = first.Next()
	}
}

// SearchN searches the range [first, last) for the first sequence of count
// identical elements, each equal to the given value.
func SearchN[T comparable, It ForwardReader[T, It]](first, last It, count int, v T) It {
	return SearchNBy(first, last, count, v, _eq2[T])
}

// SearchNBy searches the range [first, last) for the first sequence of count
// identical elements.
//
// Elements are compared using the given binary comparer eq.
func SearchNBy[T1, T2 any, It ForwardReader[T1, It]](first, last It, count int, v T2, eq EqComparer[T1, T2]) It {
	if count <= 0 {
		return first
	}
	for ; !__iter_eq(first, last); first = first.Next() {
		if !eq(first.Read(), v) {
			continue
		}
		candidate := first
		var curCount int
		for {
			curCount++
			if curCount >= count {
				return candidate
			}
			if first = first.Next(); __iter_eq(first, last) {
				return last
			}
			if !eq(first.Read(), v) {
				break
			}
		}
	}
	return last
}

// Copy copies the elements in the range, defined by [first, last), to another
// range beginning at dFirst.
//
// It returns an iterator in the destination range, pointing past the last
// element copied.
func Copy[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out) Out {
	return CopyIf(first, last, dFirst, _true1[T])
}

// CopyIf copies the elements in the range, defined by [first, last), and
// predicate pred returns true, to another range beginning at dFirst.
//
// It returns an iterator in the destination range, pointing past the last
// element copied.
func CopyIf[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, pred UnaryPredicate[T]) Out {
	for ; !__iter_eq(first, last); first = first.Next() {
		if v := first.Read(); pred(v) {
			dFirst = _writeNext(dFirst, v)
		}
	}
	return dFirst
}

// CopyN copies exactly count values from the range beginning at first to the
// range beginning at dFirst.
//
// It returns an iterator in the destination range, pointing past the last
// element copied.
func CopyN[T any, In InputIter[T, In], Out OutputIter[T]](first In, count int, dFirst Out) Out {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, first.Read())
		first = first.Next()
	}
	return dFirst
}

// CopyBackward copies the elements from the range, defined by [first, last), to
// another range ending at dLast.
//
// The elements are copied in reverse order (the last element is copied first),
// but their relative order is preserved. It returns an iterator to the last
// element copied.
func CopyBackward[T any, In BidiReader[T, In], Out BidiWriter[T, Out]](first, last In, dLast Out) Out {
	for !__iter_eq(first, last) {
		last, dLast = last.Prev(), dLast.Prev()
		dLast.Write(last.Read())
	}
	return dLast
}

// Fill assigns the given value to the elements in the range [first, last).
func Fill[T any, It ForwardWriter[T, It]](first, last It, v T) {
	for ; !__iter_eq(first, last); first = first.Next() {
		first.Write(v)
	}
}

// FillN assigns the given value to the first count elements in the range
// beginning at dFirst.
//
// If count <= 0, it does nothing.
func FillN[T any, Out OutputIter[T]](dFirst Out, count int, v T) {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, v)
	}
}

// Transform applies the given function to the range [first, last) and stores
// the result in another range, beginning at dFirst.
func Transform[T1, T2 any, In InputIter[T1, In], Out OutputIter[T2]](first, last In, dFirst Out, op UnaryOperation[T1, T2]) Out {
	for ; !__iter_eq(first, last); first = first.Next() {
		dFirst = _writeNext(dFirst, op(first.Read()))
	}
	return dFirst
}

// TransformBinary applies the given function to the two ranges [first, last),
// [first2, first2+last-first) and stores the result in another range, beginning
// at dFirst.
func TransformBinary[T1, T2, T3 any, In1 ForwardReader[T1, In1], In2 ForwardReader[T2, In2], Out OutputIter[T3]](first1, last1 In1, first2 In2, dFirst Out, op BinaryOperation[T1, T2, T3]) Out {
	for ; !__iter_eq(first1, last1); first1, first2 = first1.Next(), first2.Next() {
		dFirst = _writeNext(dFirst, op(first1.Read(), first2.Read()))
	}
	return dFirst
}

// Generate assigns each element in range [first, last) a value generated by the
// given function object g.
func Generate[T any, It ForwardWriter[T, It]](first, last It, g Generator[T]) {
	for ; !__iter_eq(first, last); first = first.Next() {
		first.Write(g())
	}
}

// GenerateN assigns values, generated by given function object g, to the first
// count elements in the range beginning at dFirst.
//
// If count <= 0, it does nothing.
func GenerateN[T any, It OutputIter[T]](dFirst It, count int, g Generator[T]) It {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, g())
	}
	return dFirst
}

// Remove removes all elements equal to v from the range [first, last) and
// returns a past-the-end iterator for the new end of the range.
func Remove[T comparable, It ForwardReadWriter[T, It]](first, last It, v T) It {
	return RemoveIf(first, last, _eq1(v))
}

// RemoveIf removes all elements which predicate function returns true from the
// range [first, last) and returns a past-the-end iterator for the new end of
// the range.
func RemoveIf[T any, It ForwardReadWriter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	first = FindIf(first, last, pred)
	if !__iter_eq(first, last) {
		for i := first.Next(); !__iter_eq(i, last); i = i.Next() {
			if !pred(i.Read()) {
				first.Write(i.Read())
				first = first.Next()
			}
		}
	}
	return first
}

// RemoveCopy copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements equal to v.
//
// Source and destination ranges cannot overlap.
func RemoveCopy[T comparable, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, v T) Out {
	return RemoveCopyIf(first, last, dFirst, _eq1(v))
}

// RemoveCopyIf copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements which predicate function returns
// true.
//
// Source and destination ranges cannot overlap.
func RemoveCopyIf[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, pred UnaryPredicate[T]) Out {
	for ; !__iter_eq(first, last); first = first.Next() {
		if v := first.Read(); !pred(v) {
			dFirst = _writeNext(dFirst, v)
		}
	}
	return dFirst
}

// Replace replaces all elements equal to old with new in the range [first,
// last).
func Replace[T comparable, It ForwardReadWriter[T, It]](first, last It, old, new T) {
	ReplaceIf(first, last, _eq1(old), new)
}

// ReplaceIf replaces all elements satisfy pred with new in the range [first,
// last).
func ReplaceIf[T any, It ForwardReadWriter[T, It]](first, last It, pred UnaryPredicate[T], v T) {
	for ; !__iter_eq(first, last); first = first.Next() {
		if pred(first.Read()) {
			first.Write(v)
		}
	}
}

// ReplaceCopy copies the elements from the range [first, last) to another range
// beginning at dFirst replacing all elements equal to old with new.
//
// The source and destination ranges cannot overlap.
func ReplaceCopy[T comparable, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, old, new T) Out {
	return ReplaceCopyIf(first, last, dFirst, _eq1(old), new)
}

// ReplaceCopyIf copies the elements from the range [first, last) to another
// range beginning at dFirst replacing all elements satisfy pred with new.
//
// The source and destination ranges cannot overlap.
func ReplaceCopyIf[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, pred UnaryPredicate[T], v T) Out {
	for ; !__iter_eq(first, last); first = first.Next() {
		if v0 := first.Read(); pred(v0) {
			dFirst = _writeNext(dFirst, v)
		} else {
			dFirst = _writeNext(dFirst, v0)
		}
	}
	return dFirst
}

// Swap swaps value of two iterators.
func Swap[T any](a, b ReadWriter[T]) {
	va, vb := a.Read(), b.Read()
	a.Write(vb)
	b.Write(va)
}

// SwapRanges exchanges elements between range [first1, last1) and another range
// starting at first2.
func SwapRanges[T any, It1 ForwardReadWriter[T, It1], It2 ForwardReadWriter[T, It2]](first1, last1 It1, first2 It2) {
	for ; !__iter_eq(first1, last1); first1, first2 = first1.Next(), first2.Next() {
		Swap[T](first1, first2)
	}
}

// Reverse reverses the order of the elements in the range [first, last).
func Reverse[T any, It BidiReadWriter[T, It]](first, last It) {
	for ; !__iter_eq(first, last); first = first.Next() {
		last = last.Prev()
		if __iter_eq(first, last) {
			return
		}
		Swap[T](first, last)
	}
}

// ReverseCopy copies the elements from the range [first, last) to another range
// beginning at dFirst in such a way that the elements in the new range are in
// reverse order.
func ReverseCopy[T any, In BidiReader[T, In], Out OutputIter[T]](first, last In, dFirst Out) Out {
	for !__iter_eq(first, last) {
		last = last.Prev()
		dFirst = _writeNext(dFirst, last.Read())
	}
	return dFirst
}

// Rotate performs a left rotation on a range of elements in such a way, that
// the element nFirst becomes the first element of the new range and nFirst - 1
// becomes the last element.
func Rotate[T any, It ForwardReadWriter[T, It]](first, nFirst, last It) It {
	if __iter_eq(first, nFirst) {
		return last
	}
	if __iter_eq(nFirst, last) {
		return first
	}
	read, write, nextRead := nFirst, first, first
	for !__iter_eq(read, last) {
		if __iter_eq(write, nextRead) {
			nextRead = read
		}
		Swap[T](write, read)
		write, read = write.Next(), read.Next()
	}
	Rotate[T](write, nextRead, last)
	return write
}

// RotateCopy copies the elements from the range [first, last), to another range
// beginning at dFirst in such a way, that the element nFirst becomes the first
// element of the new range and nFirst - 1 becomes the last element.
func RotateCopy[T any, In ForwardReader[T, In], Out OutputIter[T]](first, nFirst, last In, dFirst Out) Out {
	return Copy[T](first, nFirst, Copy[T](nFirst, last, dFirst))
}

// Shuffle reorders the elements in the given range [first, last) such that each
// possible permutation of those elements has equal probability of appearance.
func Shuffle[T any, It RandomReadWriter[T, It]](first, last It, r *rand.Rand) {
	r.Shuffle(first.Distance(last), func(i, j int) {
		Swap[T](first.AdvanceN(i), first.AdvanceN(j))
	})
}

// Sample selects n elements from the sequence [first; last) such that each
// possible sample has equal probability of appearance, and writes those
// selected elements into the output iterator out.
func Sample[T any, In ForwardReader[T, In], Out OutputIter[T]](first, last In, out Out, n int, r *rand.Rand) Out {
	_, rr := any(first).(RandomReader[T, In])
	rout, rw := any(out).(RandomWriter[T, Out])
	if !rr && rw {
		// reservoir sampling
		var k int
		for ; !__iter_eq(first, last) && k < n; first, k = first.Next(), k+1 {
			rout.AdvanceN(k).Write(first.Read())
		}
		if __iter_eq(first, last) {
			return rout.AdvanceN(k)
		}
		sz := k
		for ; !__iter_eq(first, last); first, k = first.Next(), k+1 {
			if d := r.Intn(k + 1); d < sz {
				rout.AdvanceN(d).Write(first.Read())
			}
		}
		return rout.AdvanceN(n)
	}
	// selection sampling
	unsampled := Distance[T](first, last)
	if n > unsampled {
		n = unsampled
	}
	for ; n != 0; first = first.Next() {
		if r.Intn(unsampled) < n {
			out = _writeNext(out, first.Read())
			n--
		}
		unsampled--
	}
	return out
}

// Unique eliminates all but the first element from every consecutive group of
// equivalent elements from the range [first, last) and returns a past-the-end
// iterator for the new logical end of the range.
func Unique[T comparable, It ForwardReadWriter[T, It]](first, last It) It {
	return UniqueIf(first, last, _eq2[T])
}

// UniqueIf eliminates all but the first element from every consecutive group of
// equivalent elements from the range [first, last) and returns a past-the-end
// iterator for the new logical end of the range.
//
// Elements are compared using the given binary comparer eq.
func UniqueIf[T any, It ForwardReadWriter[T, It]](first, last It, eq EqComparer[T, T]) It {
	if __iter_eq(first, last) {
		return last
	}
	result := first
	for {
		first = first.Next()
		if __iter_eq(first, last) {
			return result.Next()
		}
		if !eq(result.Read(), first.Read()) {
			if result = result.Next(); !__iter_eq(result, first) {
				result.Write(first.Read())
			}
		}
	}
}

// UniqueCopy copies the elements from the range [first, last), to another range
// beginning at dFirst in such a way that there are no consecutive equal
// elements.
//
// Only the first element of each group of equal elements is copied.
func UniqueCopy[T comparable, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out) Out {
	return UniqueCopyIf(first, last, dFirst, _eq2[T])
}

// UniqueCopyIf copies the elements from the range [first, last), to another
// range beginning at dFirst in such a way that there are no consecutive equal
// elements.
//
// Only the first element of each group of equal elements is copied. Elements
// are compared using the given binary comparer eq.
func UniqueCopyIf[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, eq EqComparer[T, T]) Out {
	if !__iter_eq(first, last) {
		v := first.Read()
		dFirst = _writeNext(dFirst, v)
		for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
			if v1 := first.Read(); !eq(v, v1) {
				v = v1
				dFirst = _writeNext(dFirst, v)
			}
		}
	}
	return dFirst
}

// IsPartitioned returns true if all elements in the range [first, last) that
// satisfy the predicate pred appear before all elements that don't. Also returns
// true if [first, last) is empty.
func IsPartitioned[T any, It InputIter[T, It]](first, last It, pred UnaryPredicate[T]) bool {
	return NoneOf(FindIfNot(first, last, pred), last, pred)
}

// Partition reorders the elements in the range [first, last) in such a way that
// all elements for which the predicate pred returns true precede the elements
// for which predicate pred returns false.
//
// Relative order of the elements is not preserved.
func Partition[T any, It ForwardReadWriter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	first = FindIfNot(first, last, pred)
	if __iter_eq(first, last) {
		return first
	}
	for i := first.Next(); !__iter_eq(i, last); i = i.Next() {
		if pred(i.Read()) {
			Swap[T](first, i)
			first = first.Next()
		}
	}
	return first
}

// PartitionCopy copies the elements from the range [first, last) to two
// different ranges depending on the value returned by the predicate pred. The
// elements that satisfy the predicate pred are copied to the range beginning at
// outTrue. The rest of the elements are copied to the range beginning at
// outFalse.
func PartitionCopy[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, outTrue, outFalse Out, pred UnaryPredicate[T]) (Out, Out) {
	for ; !__iter_eq(first, last); first = first.Next() {
		if v := first.Read(); pred(v) {
			outTrue = _writeNext(outTrue, v)
		} else {
			outFalse = _writeNext(outFalse, v)
		}
	}
	return outTrue, outFalse
}

// StablePartition reorders the elements in the range [first, last) in such a
// way that all elements for which the predicate pred returns true precede the
// elements for which predicate pred returns false. Relative order of the
// elements is preserved.
func StablePartition[T any, It ForwardReadWriter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	for {
		if __iter_eq(first, last) {
			return first
		}
		if !pred(first.Read()) {
			break
		}
		first = first.Next()
	}
	return _stablePartitionForward(first, last, pred, Distance[T](first, last))
}

func _stablePartitionForward[T any, It ForwardReadWriter[T, It]](first, last It, pred UnaryPredicate[T], l int) It {
	if l == 1 {
		return first
	}
	if l == 2 {
		m := first.Next()
		if pred(m.Read()) {
			Swap[T](first, m)
			return m
		}
		return first
	}
	l2 := l / 2
	m := AdvanceN[T](first, l2)
	// F?????????????????
	// f       m         l
	firstFalse := _stablePartitionForward(first, m, pred, l2)
	// TTTFFFFF??????????
	// f  ff   m         l
	m1, lh := m, l-l2
	for pred(m1.Read()) {
		m1 = m1.Next()
		if __iter_eq(m1, last) {
			break
		}
		lh--
	}
	secondFalse := last
	if !__iter_eq(m1, last) {
		// TTTFFFFFTTTF??????
		// f  ff   m  m1     l
		secondFalse = _stablePartitionForward(m1, last, pred, lh)
	}
	// TTTFFFFFTTTTTFFFFF
	// f  ff   m    sf   l
	return Rotate[T](firstFalse, m, secondFalse)
}

// StablePartitionBidi reorders the elements in the range [first, last) in such a
// way that all elements for which the predicate pred returns true precede the
// elements for which predicate pred returns false. Relative order of the
// elements is preserved.
func StablePartitionBidi[T any, It BidiReadWriter[T, It]](first, last It, pred UnaryPredicate[T]) It {
	for {
		if __iter_eq(first, last) {
			return first
		}
		if !pred(first.Read()) {
			break
		}
		first = first.Next()
	}
	for {
		last = last.Prev()
		if __iter_eq(first, last) {
			return first
		}
		if pred(last.Read()) {
			break
		}
	}
	return _stablePartitionBidi(first, last, pred, Distance[T](first, last)+1)
}

func _stablePartitionBidi[T any, It BidiReadWriter[T, It]](first, last It, pred UnaryPredicate[T], l int) It {
	if l == 2 {
		Swap[T](first, last)
		return last
	}
	if l == 3 {
		m := first.Next()
		if pred(m.Read()) {
			Swap[T](first, m)
			Swap[T](m, last)
			return last
		}
		Swap[T](m, last)
		Swap[T](first, m)
		return m
	}
	m, l2 := first, l/2
	m = AdvanceN[T](m, l2)
	// F???????????????T
	// f       m       l
	m1, lh := m, l2
	for m1 = m1.Prev(); !pred(m1.Read()); m1 = m1.Prev() {
		if __iter_eq(m1, first) {
			break
		}
		lh--
	}
	firstFalse := first
	if !__iter_eq(m1, first) {
		// F????TFF????????T
		// f    m1 m       l
		firstFalse = _stablePartitionBidi(first, m1, pred, lh)
	}
	// TTFFFFFF????????T
	// f ff m1 m       l
	m1, lh = m, l-l2
	for pred(m1.Read()) {
		m1 = m1.Next()
		if __iter_eq(m1, last) {
			break
		}
		lh--
	}
	secondFalse := last.Next()
	if !__iter_eq(m1, last) {
		// TTFFFFFFTTTF?????T
		// f ff m1 m  m1    l
		secondFalse = _stablePartitionBidi(m1, last, pred, lh)
	}
	// TTFFFFFFTTTTTTFFFF
	// f ff m1 m  m1 sf l
	return Rotate[T](firstFalse, m, secondFalse)
}

// PartitionPoint examines the partitioned (as if by Partition) range [first,
// last) and locates the end of the first partition, that is, the first element
// that does not satisfy pred or last if all elements satisfy pred.
func PartitionPoint[T any, It ForwardReader[T, It]](first, last It, pred UnaryPredicate[T]) It {
	l := Distance[T](first, last)
	for l != 0 {
		l2 := l / 2
		m := AdvanceN[T](first, l2)
		if pred(m.Read()) {
			first = m.Next()
			l -= l2 + 1
		} else {
			l = l2
		}
	}
	return first
}

// IsSorted checks if the elements in range [first, last) are sorted in
// non-descending order.
func IsSorted[T Ordered, It ForwardReader[T, It]](first, last It) bool {
	return __iter_eq(IsSortedUntil[T](first, last), last)
}

// IsSortedBy checks if the elements in range [first, last) are sorted in
// non-descending order.
//
// Elements are compared using the given binary comparer less.
func IsSortedBy[T any, It ForwardReader[T, It]](first, last It, less LessComparer[T]) bool {
	return __iter_eq(IsSortedUntilBy(first, last, less), last)
}

// IsSortedUntil examines the range [first, last) and finds the largest range
// beginning at first in which the elements are sorted in ascending order.
func IsSortedUntil[T Ordered, It ForwardReader[T, It]](first, last It) It {
	return IsSortedUntilBy(first, last, _less[T])
}

// IsSortedUntilBy examines the range [first, last) and finds the largest range
// beginning at first in which the elements are sorted in ascending order.
//
// Elements are compared using the given binary comparer less.
func IsSortedUntilBy[T any, It ForwardReader[T, It]](first, last It, less LessComparer[T]) It {
	if !__iter_eq(first, last) {
		for next := first.Next(); !__iter_eq(next, last); next = next.Next() {
			if less(next.Read(), first.Read()) {
				return next
			}
			first = next
		}
	}
	return last
}

// Adapt RandomIter to sort.Interface.
type sortHelper[T any, It RandomReadWriter[T, It]] struct {
	first RandomReadWriter[T, It]
	n     int
	less  LessComparer[T]
}

func (s *sortHelper[T, It]) Len() int {
	return s.n
}

func (s *sortHelper[T, It]) Less(i, j int) bool {
	return s.less(
		s.first.AdvanceN(i).Read(),
		s.first.AdvanceN(j).Read(),
	)
}

func (s *sortHelper[T, It]) Swap(i, j int) {
	it1, it2 := s.first.AdvanceN(i), s.first.AdvanceN(j)
	v1, v2 := it1.Read(), it2.Read()
	it1.Write(v2)
	it2.Write(v1)
}

// Adapt RandomIter to sort.Interface.
type heapHelper[T any, It RandomReadWriter[T, It]] struct {
	*sortHelper[T, It]
}

func (h *heapHelper[T, It]) Less(i, j int) bool {
	return h.less(
		h.first.AdvanceN(j).Read(),
		h.first.AdvanceN(i).Read(),
	)
}

func (h *heapHelper[T, It]) Push(x any) {
	h.n++
}

func (h *heapHelper[T, It]) Pop() any {
	h.n--
	return nil
}

// Sort sorts the elements in the range [first, last) in ascending order. The
// order of equal elements is not guaranteed to be preserved.
func Sort[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	SortBy(first, last, _less[T])
}

// SortBy sorts the elements in the range [first, last) in ascending order. The
// order of equal elements is not guaranteed to be preserved.
//
// Elements are compared using the given binary comparer less.
func SortBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	sort.Sort(&sortHelper[T, It]{
		first: first,
		n:     first.Distance(last),
		less:  less,
	})
}

// PartialSort rearranges elements such that the range [first, middle) contains
// the sorted (middle-first) smallest elements in the range [first, last).
//
// The order of equal elements is not guaranteed to be preserved. The order of
// the remaining elements in the range [middle, last) is unspecified.
func PartialSort[T Ordered, It RandomReadWriter[T, It]](first, middle, last It) {
	PartialSortBy(first, middle, last, _less[T])
}

// PartialSortBy rearranges elements such that the range [first, middle)
// contains the sorted (middle-first) smallest elements in the range [first,
// last).
//
// The order of equal elements is not guaranteed to be preserved. The order of
// the remaining elements in the range [middle, last) is unspecified. Elements
// are compared using the given binary comparer less.
func PartialSortBy[T any, It RandomReadWriter[T, It]](first, middle, last It, less LessComparer[T]) {
	MakeHeapBy(first, middle, less)
	for i := middle; !__iter_eq(i, last); i = i.Next() {
		if less(i.Read(), first.Read()) {
			Swap[T](first, i)
			heap.Fix(&heapHelper[T, It]{
				&sortHelper[T, It]{
					first: first,
					n:     first.Distance(middle),
					less:  less,
				}}, 0)
		}
	}
	SortHeapBy(first, middle, less)
}

// PartialSortCopy sorts some of the elements in the range [first, last) in
// ascending order, storing the result in the range [dFirst, dLast).
//
// At most dLast - dFirst of the elements are placed sorted to the range
// [dFirst, dFirst + n). n is the number of elements to sort (n = min(last -
// first, dLast - dFirst)). The order of equal elements is not guaranteed to be
// preserved.
func PartialSortCopy[T Ordered, In InputIter[T, In], Out RandomReadWriter[T, Out]](first, last In, dFirst, dLast Out) {
	PartialSortCopyBy(first, last, dFirst, dLast, _less[T])
}

// PartialSortCopyBy sorts some of the elements in the range [first, last) in
// ascending order, storing the result in the range [dFirst, dLast).
//
// At most dLast - dFirst of the elements are placed sorted to the range
// [dFirst, dFirst + n). n is the number of elements to sort (n = min(last -
// first, dLast - dFirst)). The order of equal elements is not guaranteed to be
// preserved. Elements are compared using the given binary comparer less.
func PartialSortCopyBy[T any, In InputIter[T, In], Out RandomReadWriter[T, Out]](first, last In, dFirst, dLast Out, less LessComparer[T]) {
	if __iter_eq(dFirst, dLast) {
		return
	}
	r, len := dFirst, dFirst.Distance(dLast)
	for ; !__iter_eq(first, last) && !__iter_eq(r, dLast); first, r = first.Next(), r.Next() {
		r.Write(first.Read())
	}
	MakeHeapBy(dFirst, dLast, less)
	for ; !__iter_eq(first, last); first = first.Next() {
		if less(first.Read(), dFirst.Read()) {
			dFirst.Write(first.Read())
			heap.Fix(&heapHelper[T, Out]{
				&sortHelper[T, Out]{
					first: dFirst,
					n:     len,
					less:  less,
				}}, 0)
		}
	}
	SortHeapBy(dFirst, r, less)
}

// StableSort sorts the elements in the range [first, last) in ascending order.
// The order of equivalent elements is guaranteed to be preserved.
func StableSort[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	StableSortBy(first, last, _less[T])
}

// StableSortBy sorts the elements in the range [first, last) in ascending
// order.
//
// The order of equivalent elements is guaranteed to be preserved. Elements are
// compared using the given binary comparer less.
func StableSortBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	sort.Stable(&sortHelper[T, It]{
		first: first,
		n:     first.Distance(last),
		less:  less,
	})
}

// NthElement is a partial sorting algorithm that rearranges elements in [first,
// last) such that:
// a. The element pointed at by nth is changed to whatever element would occur
// in that position if [first, last) were sorted.
// b. All of the elements before this new nth element are less than or equal to
// the elements after the new nth element.
func NthElement[T Ordered, It RandomReadWriter[T, It]](first, nth, last It) {
	NthElementBy(first, nth, last, _less[T])
}

// NthElementBy is a partial sorting algorithm that rearranges elements in
// [first, last) such that:
// a. The element pointed at by nth is changed to whatever element would occur
// in that position if [first, last) were sorted.
// b. All of the elements before this new nth element are less than or equal to
// the elements after the new nth element.
//
// Elements are compared using the given binary comparer less.
func NthElementBy[T any, It RandomReadWriter[T, It]](first, nth, last It, less LessComparer[T]) {
Restart:
	for {
		if __iter_eq(nth, last) {
			return
		}
		len := first.Distance(last)
		if len <= 7 {
			SortBy(first, last, less)
			return
		}

		m := first.AdvanceN(len / 2)
		last1 := last.Prev()

		// sort {first, m, last1}
		var maybeSorted bool
		if !less(m.Read(), first.Read()) {
			// first<=m
			if !less(last1.Read(), m.Read()) {
				// first<=m<=last1
				maybeSorted = true
			} else {
				// first<=m,m>last1
				Swap[T](m, last1)
				// first<=last1,m<last1
				if less(m.Read(), first.Read()) {
					// m<first<=last1
					Swap[T](first, m)
					// first<m<=last1
				}
				// first<=m<last1
			}
		} else if less(last1.Read(), m.Read()) {
			// first>m>last1
			Swap[T](first, last1)
			// first<m<last1
		} else {
			// first>m,m<=last1
			Swap[T](first, m)
			// first<m,first<=last1
			if less(last1.Read(), m.Read()) {
				// first<=last1<m
				Swap[T](m, last1)
				// first<=m<last1
			}
			// first<m<=last1
		}

		i, j := first, last1
		// -????????0???????+   // 0: pivot, -: <=pivot, <: <pivot, +: >=pivot, >: >pivot
		// f        m        l
		// i                j
		if !less(i.Read(), m.Read()) {
			// 0????????0???????+
			// f        m        l
			// i                j
			for {
				if j = j.Prev(); __iter_eq(i, j) {
					// 0+++++++++++++
					// f             l
					// i
					// j
					if i, j = i.Next(), last1; !less(first.Read(), j.Read()) {
						// 0++++++++++++0
						// fi           jl
						for {
							if __iter_eq(i, j) {
								// 00000000000000
								// f            jl
								//              i
								return
							}
							if less(first.Read(), i.Read()) {
								// 00000>+++++++0
								// f    i       jl
								Swap[T](i, j)
								maybeSorted = false
								i = i.Next()
								break
							}
							i = i.Next()
						}
					}
					// 000000+++++++>
					// f     i      jl
					if __iter_eq(i, j) {
						// 0000000000000>
						// f            jl
						//              i
						return
					}
					for {
						for !less(first.Read(), i.Read()) {
							i = i.Next()
						}
						for j = j.Prev(); less(first.Read(), j.Read()); j = j.Prev() {
						}
						// 000000>+++++0++
						// f     i     j  l
						if !i.Less(j) {
							break
						}
						Swap[T](i, j)
						maybeSorted = false
						i = i.Next()
					}
					// 000000000+++++++
					// f       ji      l
					if nth.Less(i) {
						return
					}
					first = i
					continue Restart
				}
				if less(j.Read(), m.Read()) {
					// 0???-+++++++++
					// f             l
					// i   j
					Swap[T](i, j)
					maybeSorted = false
					break
				}
			}
		}

		// i.Read() < m.Read()
		i = i.Next()
		// -??????0????????+         -??????0????0+++
		// fi     m        jl  [OR]  fi     m    j   l
		if i.Less(j) {
			for {
				for less(i.Read(), m.Read()) {
					i = i.Next()
				}
				for j = j.Prev(); !less(j.Read(), m.Read()); j = j.Prev() {
				}
				// ----+??0?????<+++        -------0?????<+++
				// f      m         l [OR]  f      m         l
				//     i        j                  i     j
				if !i.Less(j) {
					// -------0--+++++++       -------0+++++++++
					// f      m         l [OR] f      m         l
					//          ji                   ji
					break
				}
				Swap[T](i, j)
				maybeSorted = false
				if __iter_eq(m, i) {
					m = j
				}
				// -----??0?????++++       --------?????0+++
				// f   i  m     j   l [OR] f            m   l
				//                                i     j
				i = i.Next()
			}
		}

		// -------+++0+++
		// f      i  m   l
		if !__iter_eq(i, m) && less(m.Read(), i.Read()) {
			Swap[T](i, m)
			maybeSorted = false
		}
		// -------0++++++
		// f      i      l
		if __iter_eq(nth, i) {
			return
		}
		if nth.Less(i) {
			if maybeSorted && IsSortedBy(first, i, less) {
				return
			}
			last = i
		} else {
			if maybeSorted && IsSortedBy(i, last, less) {
				return
			}
			first = i.Next()
		}
	}
}

// LowerBound returns an iterator pointing to the first element in the range
// [first, last) that is not less than (i.e. greater or equal to) value, or last
// if no such element is found.
func LowerBound[T Ordered, It ForwardReader[T, It]](first, last It, v T) It {
	return LowerBoundBy(first, last, v, _less[T])
}

// LowerBoundBy returns an iterator pointing to the first element in the range
// [first, last) that is not less than (i.e. greater or equal to) value, or last
// if no such element is found.
//
// Elements are compared using the given binary comparer less.
func LowerBoundBy[T any, It ForwardReader[T, It]](first, last It, v T, less LessComparer[T]) It {
	for len := Distance[T](first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN[T](first, l2)
		if less(m.Read(), v) {
			first = m.Next()
			len -= l2 + 1
		} else {
			len = l2
		}
	}
	return first
}

// UpperBound returns an iterator pointing to the first element in the range
// [first, last) that is greater than value, or last if no such element is
// found.
func UpperBound[T Ordered, It ForwardReader[T, It]](first, last It, v T) It {
	return UpperBoundBy(first, last, v, _less[T])
}

// UpperBoundBy returns an iterator pointing to the first element in the range
// [first, last) that is greater than value, or last if no such element is
// found.
//
// Elements are compared using the given binary comparer less.
func UpperBoundBy[T any, It ForwardReader[T, It]](first, last It, v T, less LessComparer[T]) It {
	for len := Distance[T](first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN[T](first, l2)
		if less(v, m.Read()) {
			len = l2
		} else {
			first = m.Next()
			len -= l2 + 1
		}
	}
	return first
}

// BinarySearch checks if an element equivalent to value appears within the
// range [first, last).
func BinarySearch[T Ordered, It ForwardReader[T, It]](first, last It, v T) bool {
	return BinarySearchBy(first, last, v, _less[T])
}

// BinarySearchBy checks if an element equivalent to value appears within the
// range [first, last).
//
// Elements are compared using the given binary comparer less.
func BinarySearchBy[T any, It ForwardReader[T, It]](first, last It, v T, less LessComparer[T]) bool {
	first = LowerBoundBy(first, last, v, less)
	return !__iter_eq(first, last) && !(less(v, first.Read()))
}

// EqualRange returns a range containing all elements equivalent to value in the
// range [first, last).
func EqualRange[T Ordered, It ForwardReader[T, It]](first, last It, v T) (It, It) {
	return EqualRangeBy(first, last, v, _less[T])
}

// EqualRangeBy returns a range containing all elements equivalent to value in
// the range [first, last).
//
// Elements are compared using the given binary comparer less.
func EqualRangeBy[T any, It ForwardReader[T, It]](first, last It, v T, less LessComparer[T]) (It, It) {
	for len := Distance[T](first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN[T](first, l2)
		if less(m.Read(), v) {
			first = m.Next()
			len -= l2 + 1
		} else if less(v, m.Read()) {
			last = m
			len = l2
		} else {
			return LowerBoundBy(first, m, v, less), UpperBoundBy(m.Next(), last, v, less)
		}
	}
	return first, first
}

// Merge merges two sorted ranges [first1, last1) and [first2, last2) into one
// sorted range beginning at dFirst.
func Merge[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out) Out {
	return MergeBy(first1, last1, first2, last2, dFirst, _less[T])
}

// MergeBy merges two sorted ranges [first1, last1) and [first2, last2) into one
// sorted range beginning at dFirst.
//
// Elements are compared using the given binary comparer less.
func MergeBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out, less LessComparer[T]) Out {
	for !__iter_eq(first1, last1) {
		if __iter_eq(first2, last2) {
			return Copy[T](first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v2, v1) {
			dFirst = _writeNext(dFirst, v2)
			first2 = first2.Next()
		} else {
			dFirst = _writeNext(dFirst, v1)
			first1 = first1.Next()
		}
	}
	return Copy[T](first2, last2, dFirst)
}

// InplaceMerge Merges two consecutive sorted ranges [first, middle) and
// [middle, last) into one sorted range [first, last). For equivalent elements
// in the original two ranges, the elements from the first range (preserving
// their original order) precede the elements from the second range (preserving
// their original order).
func InplaceMerge[T Ordered, It BidiReadWriter[T, It]](first, middle, last It) {
	InplaceMergeBy(first, middle, last, _less[T])
}

// InplaceMergeBy Merges two consecutive sorted ranges [first, middle) and
// [middle, last) into one sorted range [first, last). For equivalent elements
// in the original two ranges, the elements from the first range (preserving
// their original order) precede the elements from the second range (preserving
// their original order).
//
// Elements are compared using the given binary comparer less.
func InplaceMergeBy[T any, It BidiReadWriter[T, It]](first, middle, last It, less LessComparer[T]) {
	len1, len2 := Distance[T](first, middle), Distance[T](middle, last)
	for {
		if len2 == 0 {
			return
		}
		for {
			if len1 == 0 {
				return
			}
			if less(middle.Read(), first.Read()) {
				break
			}
			first = first.Next()
			len1--
		}
		var len11, len21 int
		var m1, m2 It
		if len1 < len2 {
			len21 = len2 / 2
			m2 = AdvanceN[T](middle, len21)
			m1 = UpperBoundBy(first, middle, m2.Read(), less)
			len11 = Distance[T](first, m1)
		} else {
			if len1 == 1 {
				Swap[T](first, middle)
				return
			}
			len11 = len1 / 2
			m1 = AdvanceN[T](first, len11)
			m2 = LowerBoundBy(middle, last, m1.Read(), less)
			len21 = Distance[T](middle, m2)
		}
		len12, len22 := len1-len11, len2-len21
		middle = Rotate[T](m1, middle, m2)
		if len11+len21 < len12+len22 {
			InplaceMergeBy(first, m1, middle, less)
			first, middle = middle, m2
			len1, len2 = len12, len22
		} else {
			InplaceMergeBy(middle, m2, last, less)
			middle, last = m1, middle
			len1, len2 = len11, len21
		}
	}
}

// Includes returns true if the sorted range [first2, last2) is a subsequence of
// the sorted range [first1, last1). (A subsequence need not be contiguous.)
func Includes[T Ordered, It1 InputIter[T, It1], It2 InputIter[T, It2]](first1, last1 It1, first2, last2 It2) bool {
	return IncludesBy(first1, last1, first2, last2, _less[T])
}

// IncludesBy returns true if the sorted range [first2, last2) is a subsequence
// of the sorted range [first1, last1). (A subsequence need not be contiguous.)
//
// Elements are compared using the given binary comparer less.
func IncludesBy[T any, It1 InputIter[T, It1], It2 InputIter[T, It2]](first1, last1 It1, first2, last2 It2, less LessComparer[T]) bool {
	for ; !__iter_eq(first2, last2); first1 = first1.Next() {
		if __iter_eq(first1, last1) || less(first2.Read(), first1.Read()) {
			return false
		}
		if !less(first1.Read(), first2.Read()) {
			first2 = first2.Next()
		}
	}
	return true
}

// SetDifference copies the elements from the sorted range [first1, last1) which
// are not found in the sorted range [first2, last2) to the range beginning at
// dFirst.
func SetDifference[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out) Out {
	return SetDifferenceBy(first1, last1, first2, last2, dFirst, _less[T])
}

// SetDifferenceBy copies the elements from the sorted range [first1, last1)
// which are not found in the sorted range [first2, last2) to the range
// beginning at dFirst.
//
// Elements are compared using the given binary comparer less.
func SetDifferenceBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out, less LessComparer[T]) Out {
	for !__iter_eq(first1, last1) {
		if __iter_eq(first2, last2) {
			return Copy[T](first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			dFirst = _writeNext(dFirst, v1)
			first1 = first1.Next()
		} else {
			if !less(v2, v1) {
				first1 = first1.Next()
			}
			first2 = first2.Next()
		}
	}
	return dFirst
}

// SetIntersection constructs a sorted range beginning at dFirst consisting of
// elements that are found in both sorted ranges [first1, last1) and [first2,
// last2). If some element is found m times in [first1, last1) and n times in
// [first2, last2), the first Min(m, n) elements will be copied from the first
// range to the destination range.
//
// The order of equivalent elements is preserved. The resulting range cannot
// overlap with either of the input ranges.
func SetIntersection[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out) Out {
	return SetIntersectionBy(first1, last1, first2, last2, dFirst, _less[T])
}

// SetIntersectionBy constructs a sorted range beginning at dFirst consisting of
// elements that are found in both sorted ranges [first1, last1) and [first2,
// last2). If some element is found m times in [first1, last1) and n times in
// [first2, last2), the first Min(m, n) elements will be copied from the first
// range to the destination range.
//
// The order of equivalent elements is preserved. The resulting range cannot
// overlap with either of the input ranges. Elements are compared using the
// given binary comparer less.
func SetIntersectionBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out, less LessComparer[T]) Out {
	for !__iter_eq(first1, last1) && !__iter_eq(first2, last2) {
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			first1 = first1.Next()
		} else {
			if !less(v2, v1) {
				dFirst = _writeNext(dFirst, v1)
				first1 = first1.Next()
			}
			first2 = first2.Next()
		}
	}
	return dFirst
}

// SetSymmetricDifference computes symmetric difference of two sorted ranges:
// the elements that are found in either of the ranges, but not in both of them
// are copied to the range beginning at dFirst. The resulting range is also
// sorted.
//
// If some element is found m times in [first1, last1) and n times in [first2,
// last2), it will be copied to dFirst exactly Abs(m-n) times. If m>n, then the
// last m-n of those elements are copied from [first1,last1), otherwise the last
// n-m elements are copied from [first2,last2). The resulting range cannot
// overlap with either of the input ranges.
func SetSymmetricDifference[T Ordered, In1 InputIter[T, In1], Out OutputIter[T]](first1, last1 In1, first2, last2 In1, dFirst Out) Out {
	return SetSymmetricDifferenceBy(first1, last1, first2, last2, dFirst, _less[T])
}

// SetSymmetricDifferenceBy computes symmetric difference of two sorted ranges:
// the elements that are found in either of the ranges, but not in both of them
// are copied to the range beginning at dFirst. The resulting range is also
// sorted.
//
// If some element is found m times in [first1, last1) and n times in [first2,
// last2), it will be copied to dFirst exactly Abs(m-n) times. If m>n, then the
// last m-n of those elements are copied from [first1,last1), otherwise the last
// n-m elements are copied from [first2,last2). The resulting range cannot
// overlap with either of the input ranges. Elements are compared using the
// given binary comparer less.
func SetSymmetricDifferenceBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out, less LessComparer[T]) Out {
	for !__iter_eq(first1, last1) {
		if __iter_eq(first2, last2) {
			return Copy[T](first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			dFirst = _writeNext(dFirst, v1)
			first1 = first1.Next()
		} else {
			if less(v2, v1) {
				dFirst = _writeNext(dFirst, v2)
			} else {
				first1 = first1.Next()
			}
			first2 = first2.Next()
		}
	}
	return Copy[T](first2, last2, dFirst)
}

// SetUnion constructs a sorted union beginning at dFirst consisting of the set
// of elements present in one or both sorted ranges [first1, last1) and [first2,
// last2).
//
// If some element is found m times in [first1, last1) and n times in [first2,
// last2), then all m elements will be copied from [first1, last1) to dFirst,
// preserving order, and then exactly Max(n-m, 0) elements will be copied from
// [first2, last2) to dFirst, also preserving order.
func SetUnion[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out) Out {
	return SetUnionBy(first1, last1, first2, last2, dFirst, _less[T])
}

// SetUnionBy constructs a sorted union beginning at dFirst consisting of the
// set of elements present in one or both sorted ranges [first1, last1) and
// [first2, last2).
//
// If some element is found m times in [first1, last1) and n times in [first2,
// last2), then all m elements will be copied from [first1, last1) to dFirst,
// preserving order, and then exactly Max(n-m, 0) elements will be copied from
// [first2, last2) to dFirst, also preserving order.  Elements are compared
// using the given binary comparer less.
func SetUnionBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2], Out OutputIter[T]](first1, last1 In1, first2, last2 In2, dFirst Out, less LessComparer[T]) Out {
	for !__iter_eq(first1, last1) {
		if __iter_eq(first2, last2) {
			return Copy[T](first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v2, v1) {
			dFirst = _writeNext(dFirst, v2)
			first2 = first2.Next()
		} else {
			dFirst = _writeNext(dFirst, v1)
			if !less(v1, v2) {
				first2 = first2.Next()
			}
			first1 = first1.Next()
		}
	}
	return Copy[T](first2, last2, dFirst)
}

// IsHeap checks if the elements in range [first, last) are a max heap.
func IsHeap[T Ordered, It RandomReader[T, It]](first, last It) bool {
	return IsHeapBy(first, last, _less[T])
}

// IsHeapBy checks if the elements in range [first, last) are a max heap.
//
// Elements are compared using the given binary comparer less.
func IsHeapBy[T any, It RandomReader[T, It]](first, last It, less LessComparer[T]) bool {
	return __iter_eq(IsHeapUntilBy(first, last, less), last)
}

// IsHeapUntil examines the range [first, last) and finds the largest range
// beginning at first which is a max heap.
func IsHeapUntil[T Ordered, It RandomReader[T, It]](first, last It) It {
	return IsHeapUntilBy(first, last, _less[T])
}

// IsHeapUntilBy examines the range [first, last) and finds the largest range
// beginning at first which is a max heap.
//
// Elements are compared using the given binary comparer less.
func IsHeapUntilBy[T any, It RandomReader[T, It]](first, last It, less LessComparer[T]) It {
	len, p, c, pp := first.Distance(last), 0, 1, first
	for c < len {
		cp := first.AdvanceN(c)
		if less(pp.Read(), cp.Read()) {
			return cp
		}
		c++
		if c == len {
			return last
		}
		cp = cp.Next()
		if less(pp.Read(), cp.Read()) {
			return cp
		}
		p++
		pp = pp.Next()
		c = 2*p + 1
	}
	return last
}

// MakeHeap constructs a max heap in the range [first, last).
func MakeHeap[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	MakeHeapBy(first, last, _less[T])
}

// MakeHeapBy constructs a max heap in the range [first, last).
//
// Elements are compared using the given binary comparer less.
func MakeHeapBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	heap.Init(&heapHelper[T, It]{
		&sortHelper[T, It]{
			first: first,
			n:     first.Distance(last),
			less:  less,
		}})
}

// PushHeap inserts the element at the position last-1 into the max heap defined
// by the range [first, last-1).
func PushHeap[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	PushHeapBy(first, last, _less[T])
}

// PushHeapBy inserts the element at the position last-1 into the max heap
// defined by the range [first, last-1).
//
// Elements are compared using the given binary comparer less.
func PushHeapBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	heap.Push(&heapHelper[T, It]{
		&sortHelper[T, It]{
			first: first,
			n:     first.Distance(last) - 1,
			less:  less,
		},
	}, nil)
}

// PopHeap swaps the value in the position first and the value in the position
// last-1 and makes the subrange [first, last-1) into a heap. This has the
// effect of removing the first element from the heap defined by the range
// [first, last).
func PopHeap[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	PopHeapBy(first, last, _less[T])
}

// PopHeapBy swaps the value in the position first and the value in the position
// last-1 and makes the subrange [first, last-1) into a heap. This has the
// effect of removing the first element from the heap defined by the range
// [first, last).
//
// Elements are compared using the given binary comparer less.
func PopHeapBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	heap.Pop(&heapHelper[T, It]{
		&sortHelper[T, It]{
			first: first,
			n:     first.Distance(last),
			less:  less,
		}})
}

// SortHeap converts the max heap [first, last) into a sorted range in ascending
// order. The resulting range no longer has the heap property.
func SortHeap[T Ordered, It RandomReadWriter[T, It]](first, last It) {
	SortHeapBy(first, last, _less[T])
}

// SortHeapBy converts the max heap [first, last) into a sorted range in ascending
// order. The resulting range no longer has the heap property.
//
// Elements are compared using the given binary comparer less.
func SortHeapBy[T any, It RandomReadWriter[T, It]](first, last It, less LessComparer[T]) {
	for ; !__iter_eq(first, last); last = last.Prev() {
		PopHeapBy(first, last, less)
	}
}

// Max returns the greater of the given values.
func Max[T Ordered](a, b T) T {
	return MaxBy(a, b, _less[T])
}

// MaxBy returns the greater of the given values.
//
// Values are compared using the given binary comparer less.
func MaxBy[T any](a, b T, less LessComparer[T]) T {
	if less(a, b) {
		return b
	}
	return a
}

// MaxElement returns the largest element in a range.
func MaxElement[T Ordered, It ForwardReader[T, It]](first, last It) It {
	return MaxElementBy(first, last, _less[T])
}

// MaxElementBy returns the largest element in a range.
//
// Values are compared using the given binary comparer less.
func MaxElementBy[T any, It ForwardReader[T, It]](first, last It, less LessComparer[T]) It {
	if __iter_eq(first, last) {
		return last
	}
	max := first
	for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
		if less(max.Read(), first.Read()) {
			max = first
		}
	}
	return max
}

// Min returns the smaller of the given values.
func Min[T Ordered](a, b T) T {
	return MinBy(a, b, _less[T])
}

// MinBy returns the smaller of the given values.
//
// Values are compared using the given binary comparer less.
func MinBy[T any](a, b T, less LessComparer[T]) T {
	if less(a, b) {
		return a
	}
	return b
}

// MinElement returns the smallest element in a range.
func MinElement[T Ordered, It ForwardReader[T, It]](first, last It) It {
	return MinElementBy(first, last, _less[T])
}

// MinElementBy returns the smallest element in a range.
//
// Values are compared using the given binary comparer less.
func MinElementBy[T any, It ForwardReader[T, It]](first, last It, less LessComparer[T]) It {
	if __iter_eq(first, last) {
		return last
	}
	min := first
	for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
		if less(first.Read(), min.Read()) {
			min = first
		}
	}
	return min
}

// Minmax returns the smaller and larger of two elements.
func Minmax[T Ordered](a, b T) (T, T) {
	return MinmaxBy(a, b, _less[T])
}

// MinmaxBy returns the smaller and larger of two elements.
//
// Values are compared using the given binary comparer less.
func MinmaxBy[T any](a, b T, less LessComparer[T]) (T, T) {
	if less(b, a) {
		return b, a
	}
	return a, b
}

// MinmaxElement returns the smallest and the largest elements in a range.
func MinmaxElement[T Ordered, It ForwardReader[T, It]](first, last It) (It, It) {
	return MinmaxElementBy(first, last, _less[T])
}

// MinmaxElementBy returns the smallest and the largest elements in a range.
//
// Values are compared using the given binary comparer less.
func MinmaxElementBy[T any, It ForwardReader[T, It]](first, last It, less LessComparer[T]) (It, It) {
	if __iter_eq(first, last) {
		return first, first
	}
	min, max := first, first
	for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
		i := first
		first = first.Next()
		if __iter_eq(first, last) {
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
func Clamp[T Ordered](v, lo, hi T) T {
	return ClampBy(v, lo, hi, _less[T])
}

// ClampBy clamps a value between a pair of boundary values.
//
// Values are compared using the given binary comparer less.
func ClampBy[T any](v, lo, hi T, less LessComparer[T]) T {
	if less(v, lo) {
		return lo
	}
	if less(hi, v) {
		return hi
	}
	return v
}

// Equal returns true if the range [first1, last1) is equal to the range
// [first2, last2), and false otherwise.
//
// If last2 is nil, it denotes first2 + (last1 - first1).
func Equal[T comparable, In1 InputIter[T, In1], In2 InputIter[T, In2]](first1, last1 In1, first2 In2, last2 *In2) bool {
	return EqualBy(first1, last1, first2, last2, _eq2[T])
}

// EqualBy returns true if the range [first1, last1) is equal to the range
// [first2, last2), and false otherwise.
//
// If last2 is nil, it denotes first2 + (last1 - first1). Elements are compared
// using the given binary comparer eq.
func EqualBy[T1, T2 any, In1 InputIter[T1, In1], In2 InputIter[T2, In2]](first1, last1 In1, first2 In2, last2 *In2, eq EqComparer[T1, T2]) bool {
	for ; !__iter_eq(first1, last1); first1, first2 = first1.Next(), first2.Next() {
		if (last2 != nil && __iter_eq(first2, *last2)) || !eq(first1.Read(), first2.Read()) {
			return false
		}
	}
	return last2 == nil || __iter_eq(first2, *last2)
}

// LexicographicalCompare checks if the first range [first1, last1) is
// lexicographically less than the second range [first2, last2).
func LexicographicalCompare[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2]](first1, last1 In1, first2, last2 In2) bool {
	return LexicographicalCompareBy(first1, last1, first2, last2, _less[T])
}

// LexicographicalCompareBy checks if the first range [first1, last1) is
// lexicographically less than the second range [first2, last2).
//
// Elements are compared using the given binary comparer less.
func LexicographicalCompareBy[T any, In1 InputIter[T, In1], In2 InputIter[T, In2]](first1, last1 In1, first2, last2 In2, less LessComparer[T]) bool {
	for ; !__iter_eq(first2, last2); first1, first2 = first1.Next(), first2.Next() {
		if __iter_eq(first1, last1) || less(first1.Read(), first2.Read()) {
			return true
		}
		if less(first2.Read(), first1.Read()) {
			return false
		}
	}
	return false
}

// LexicographicalCompareThreeWay lexicographically compares two ranges [first1,
// last1) and [first2, last2) using three-way comparison. The result will be 0
// if [first1, last1) == [first2, last2), -1 if [first1, last1) < [first2,
// last2), 1 if [first1, last1) > [first2, last2).
func LexicographicalCompareThreeWay[T Ordered, In1 InputIter[T, In1], In2 InputIter[T, In2]](first1, last1 In1, first2, last2 In2) int {
	return LexicographicalCompareThreeWayBy(first1, last1, first2, last2, _cmp[T])
}

// LexicographicalCompareThreeWayBy lexicographically compares two ranges [first1,
// last1) and [first2, last2) using three-way comparison. The result will be 0
// if [first1, last1) == [first2, last2), -1 if [first1, last1) < [first2,
// last2), 1 if [first1, last1) > [first2, last2).
//
// Elements are compared using the given binary predicate cmp.
func LexicographicalCompareThreeWayBy[T1, T2 any, In1 InputIter[T1, In1], In2 InputIter[T2, In2]](first1, last1 In1, first2, last2 In2, cmp ThreeWayComparer[T1, T2]) int {
	for ; !__iter_eq(first2, last2); first1, first2 = first1.Next(), first2.Next() {
		if __iter_eq(first1, last1) {
			return -1
		}
		if x := cmp(first1.Read(), first2.Read()); x != 0 {
			return x
		}
	}
	if __iter_eq(first1, last1) {
		return 0
	}
	return 1
}

// IsPermutation returns true if there exists a permutation of the elements in
// the range [first1, last1) that makes that range equal to the range
// [first2,last2), where last2 denotes first2 + (last1 - first1) if it was not
// given.
func IsPermutation[T comparable, It1 ForwardReader[T, It1], It2 ForwardReader[T, It2]](first1, last1 It1, first2 It2, last2 *It2) bool {
	return IsPermutationBy(first1, last1, first2, last2, _eq[T])
}

// IsPermutationBy returns true if there exists a permutation of the elements in
// the range [first1, last1) that makes that range equal to the range
// [first2,last2), where last2 denotes first2 + (last1 - first1) if it was not
// given.
//
// Elements are compared using the given binary comparer eq.
func IsPermutationBy[T any, It1 ForwardReader[T, It1], It2 ForwardReader[T, It2]](first1, last1 It1, first2 It2, last2 *It2, eq EqComparer[T, T]) bool {
	l := Distance[T](first1, last1)
	if last2 == nil {
		l2 := AdvanceN[T](first2, l)
		last2 = &l2
	} else if Distance[T](first2, *last2) != l {
		return false
	}
	first1, first2 = MismatchBy(first1, last1, first2, last2, eq)
	if __iter_eq(first1, last1) {
		return true
	}
	for i := first1; !__iter_eq(i, last1); i = i.Next() {
		pred := _eq_bind1(eq, i.Read())
		if !__iter_eq(FindIf(first1, i, pred), i) {
			continue
		}
		c2 := CountIf(first2, *last2, pred)
		if c2 == 0 || c2 != 1+CountIf(i.Next(), last1, pred) {
			return false
		}
	}
	return true
}

// NextPermutation transforms the range [first, last) into the next permutation
// from the set of all permutations that are lexicographically ordered. Returns
// true if such permutation exists, otherwise transforms the range into the
// first permutation (as if by Sort(first, last)) and returns false.
func NextPermutation[T Ordered, It BidiReadWriter[T, It]](first, last It) bool {
	return NextPermutationBy(first, last, _less[T])
}

// NextPermutationBy transforms the range [first, last) into the next
// permutation from the set of all permutations that are lexicographically
// ordered with respect to less. Returns true if such permutation exists,
// otherwise transforms the range into the first permutation (as if by
// Sort(first, last)) and returns false.
//
// Elements are compared using the given
// binary comparer less.
func NextPermutationBy[T any, It BidiReadWriter[T, It]](first, last It, less LessComparer[T]) bool {
	if __iter_eq(first, last) {
		return false
	}
	i := last.Prev()
	if __iter_eq(first, i) {
		return false
	}
	for {
		ip1 := i
		i = i.Prev()
		if less(i.Read(), ip1.Read()) {
			j := last.Prev()
			for ; !less(i.Read(), j.Read()); j = j.Prev() {
			}
			Swap[T](i, j)
			Reverse[T](ip1, last)
			return true
		}
		if __iter_eq(i, first) {
			Reverse[T](first, last)
			return false
		}
	}
}

// PrevPermutation transforms the range [first, last) into the previous
// permutation from the set of all permutations that are lexicographically
// ordered. Returns true if such permutation exists, otherwise transforms the
// range into the last permutation (as if by Sort(first, last); Reverse(first,
// last);) and returns false.
func PrevPermutation[T Ordered, It BidiReadWriter[T, It]](first, last It) bool {
	return PrevPermutationBy(first, last, _less[T])
}

// PrevPermutationBy transforms the range [first, last) into the previous
// permutation from the set of all permutations that are lexicographically
// ordered with respect to less. Returns true if such permutation exists,
// otherwise transforms the range into the last permutation (as if by
// Sort(first, last); Reverse(first, last);) and returns false.
//
// Elements are compared using the given binary comparer less.
func PrevPermutationBy[T any, It BidiReadWriter[T, It]](first, last It, less LessComparer[T]) bool {
	if __iter_eq(first, last) {
		return false
	}
	i := last.Prev()
	if __iter_eq(first, i) {
		return false
	}
	for {
		ip1 := i
		i = i.Prev()
		if less(ip1.Read(), i.Read()) {
			j := last.Prev()
			for ; !less(j.Read(), i.Read()); j = j.Prev() {
			}
			Swap[T](i, j)
			Reverse[T](ip1, last)
			return true
		}
		if __iter_eq(i, first) {
			Reverse[T](first, last)
			return false
		}
	}
}

// Iota fills the range [first, last) with sequentially increasing values,
// starting with v and repetitively evaluating v++.
func Iota[T Integer, It ForwardWriter[T, It]](first, last It, v T) {
	IotaBy(first, last, v, _inc[T])
}

// IotaBy fills the range [first, last) with sequentially increasing values,
// starting with v and repetitively evaluating inc(v).
func IotaBy[T any, It ForwardWriter[T, It]](first, last It, v T, inc UnaryOperation[T, T]) {
	for ; !__iter_eq(first, last); first, v = first.Next(), inc(v) {
		first.Write(v)
	}
}

// Accumulate computes the sum of the given value v and the elements in the
// range [first, last), using v+=x.
func Accumulate[T Numeric, It InputIter[T, It]](first, last It, v T) T {
	return AccumulateBy(first, last, v, _add[T, T])
}

// AccumulateBy computes the sum of the given value v and the elements in the
// range [first, last), using v=add(v,x).
func AccumulateBy[T1, T2 any, It InputIter[T1, It]](first, last It, v T2, add BinaryOperation[T2, T1, T2]) T2 {
	for ; !__iter_eq(first, last); first = first.Next() {
		v = add(v, first.Read())
	}
	return v
}

// InnerProduct computes inner product (i.e. sum of products) or performs
// ordered map/reduce operation on the range [first1, last1), using v=v+x*y.
func InnerProduct[T Numeric, It1 InputIter[T, It1], It2 InputIter[T, It2]](first1, last1 It1, first2 It2, v T) T {
	return InnerProductBy(first1, last1, first2, v, _add[T, T], _mul[T, T])
}

// InnerProductBy computes inner product (i.e. sum of products) or performs
// ordered map/reduce operation on the range [first1, last1), using
// v=add(v,mul(x,y)).
func InnerProductBy[T1, T2, T3, T4 any, It1 InputIter[T1, It1], It2 InputIter[T2, It2]](first1, last1 It1, first2 It2, v T4, add BinaryOperation[T4, T3, T4], mul BinaryOperation[T1, T2, T3]) T4 {
	for ; !__iter_eq(first1, last1); first1, first2 = first1.Next(), first2.Next() {
		v = add(v, mul(first1.Read(), first2.Read()))
	}
	return v
}

// AdjacentDifference computes the differences between the second and the first
// of each adjacent pair of elements of the range [first, last) and writes them
// to the range beginning at dFirst + 1. An unmodified copy of first is
// written to dFirst. Differences are calculated by cur-prev.
func AdjacentDifference[T Numeric, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out) Out {
	return AdjacentDifferenceBy(first, last, dFirst, _sub[T, T])
}

// AdjacentDifferenceBy computes the differences between the second and the
// first of each adjacent pair of elements of the range [first, last) and writes
// them to the range beginning at dFirst + 1. An unmodified copy of first is
// written to dFirst. Differences are calculated by sub(cur,prev).
func AdjacentDifferenceBy[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, sub BinaryOperation[T, T, T]) Out {
	if __iter_eq(first, last) {
		return dFirst
	}
	prev := first.Read()
	dFirst = _writeNext(dFirst, prev)
	for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
		cur := first.Read()
		dFirst = _writeNext(dFirst, sub(cur, prev))
		prev = cur
	}
	return dFirst
}

// PartialSum computes the partial sums of the elements in the subranges of the
// range [first, last) and writes them to the range beginning at dFirst. Sums
// are calculated by sum=sum+cur.
func PartialSum[T Numeric, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out) Out {
	return PartialSumBy(first, last, dFirst, _add[T, T])
}

// PartialSumBy computes the partial sums of the elements in the subranges of
// the range [first, last) and writes them to the range beginning at dFirst.
// Sums are calculated by sum=add(sum,cur).
func PartialSumBy[T any, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, add BinaryOperation[T, T, T]) Out {
	if __iter_eq(first, last) {
		return dFirst
	}
	sum := first.Read()
	dFirst = _writeNext(dFirst, sum)
	for first = first.Next(); !__iter_eq(first, last); first = first.Next() {
		sum = add(sum, first.Read())
		dFirst = _writeNext(dFirst, sum)
	}
	return dFirst
}

// ExclusiveScan computes an exclusive prefix sum operation using v=v+cur
// for the range [first, last), using v as the initial value, and
// writes the results to the range beginning at dFirst. "exclusive" means that
// the i-th input element is not included in the i-th sum.
func ExclusiveScan[T Numeric, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, v T) Out {
	return ExclusiveScanBy(first, last, dFirst, v, _add[T, T])
}

// ExclusiveScanBy computes an exclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value, and writes the
// results to the range beginning at dFirst. "exclusive" means that the i-th
// input element is not included in the i-th sum.
func ExclusiveScanBy[T1, T2 any, In InputIter[T1, In], Out OutputIter[T2]](first, last In, dFirst Out, v T2, add BinaryOperation[T2, T1, T2]) Out {
	return TransformExclusiveScanBy(first, last, dFirst, v, add, _noop[T1])
}

// InclusiveScan computes an inclusive prefix sum operation using v=v+cur
// for the range [first, last), using v as the initial value (if
// provided), and writes the results to the range beginning at dFirst.
// "inclusive" means that the i-th input element is included in the i-th sum.
func InclusiveScan[T Numeric, In InputIter[T, In], Out OutputIter[T]](first, last In, dFirst Out, v T) Out {
	return InclusiveScanBy(first, last, dFirst, v, _add[T, T])
}

// InclusiveScanBy computes an inclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value (if provided), and
// writes the results to the range beginning at dFirst. "inclusive" means that
// the i-th input element is included in the i-th sum.
func InclusiveScanBy[T1, T2 any, In InputIter[T1, In], Out OutputIter[T2]](first, last In, dFirst Out, v T2, add BinaryOperation[T2, T1, T2]) Out {
	return TransformInclusiveScanBy(first, last, dFirst, v, add, _noop[T1])
}

// TransformExclusiveScan transforms each element in the range [first, last)
// with op, then computes an exclusive prefix sum operation using v=v+cur
// for the range [first, last), using v as the initial value, and
// writes the results to the range beginning at dFirst. "exclusive" means that
// the i-th input element is not included in the i-th sum.
func TransformExclusiveScan[T1, T2 Numeric, In InputIter[T1, In], Out OutputIter[T2]](first, last In, dFirst Out, v T2, op UnaryOperation[T1, T2]) Out {
	return TransformExclusiveScanBy(first, last, dFirst, v, _add[T2, T2], op)
}

// TransformExclusiveScanBy transforms each element in the range [first, last)
// with op, then computes an exclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value, and writes the
// results to the range beginning at dFirst. "exclusive" means that the i-th
// input element is not included in the i-th sum.
func TransformExclusiveScanBy[T1, T2, T3 any, In InputIter[T1, In], Out OutputIter[T3]](first, last In, dFirst Out, v T3, add BinaryOperation[T3, T2, T3], op UnaryOperation[T1, T2]) Out {
	if __iter_eq(first, last) {
		return dFirst
	}
	saved := v
	for {
		v = add(v, op(first.Read()))
		dFirst = _writeNext(dFirst, saved)
		saved = v
		first = first.Next()
		if __iter_eq(first, last) {
			break
		}
	}
	return dFirst
}

// TransformInclusiveScan transforms each element in the range [first, last)
// with op, then computes an inclusive prefix sum operation using v=v+cur
// for the range [first, last), using v as the initial value (if
// provided), and writes the results to the range beginning at dFirst.
// "inclusive" means that the i-th input element is included in the i-th sum.
func TransformInclusiveScan[T1, T2 Numeric, In InputIter[T1, In], Out OutputIter[T2]](first, last In, dFirst Out, v T2, op UnaryOperation[T1, T2]) Out {
	return TransformInclusiveScanBy(first, last, dFirst, v, _add[T2, T2], op)
}

// TransformInclusiveScanBy transforms each element in the range [first, last)
// with op, then computes an inclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value (if provided), and
// writes the results to the range beginning at dFirst. "inclusive" means that
// the i-th input element is included in the i-th sum.
func TransformInclusiveScanBy[T1, T2, T3 any, In InputIter[T1, In], Out OutputIter[T3]](first, last In, dFirst Out, v T3, add BinaryOperation[T3, T2, T3], op UnaryOperation[T1, T2]) Out {
	for ; !__iter_eq(first, last); first = first.Next() {
		v = add(v, op(first.Read()))
		dFirst = _writeNext(dFirst, v)
	}
	return dFirst
}
