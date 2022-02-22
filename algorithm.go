package iter

import (
	"container/heap"
	"math/rand"
	"sort"
)

// AllOf checks if unary predicate pred returns true for all elements in the
// range [first, last).
func AllOf(first, last InputIter, pred UnaryPredicate) bool {
	return _eq(FindIfNot(first, last, pred), last)
}

// AnyOf checks if unary predicate pred returns true for at least one element in
// the range [first, last).
func AnyOf(first, last InputIter, pred UnaryPredicate) bool {
	return _ne(FindIf(first, last, pred), last)
}

// NoneOf checks if unary predicate pred returns true for no elements in the
// range [first, last).
func NoneOf(first, last InputIter, pred UnaryPredicate) bool {
	return _eq(FindIf(first, last, pred), last)
}

// ForEach applies the given function f to the result of dereferencing every
// iterator in the range [first, last), in order.
func ForEach(first, last InputIter, f IteratorFunction) IteratorFunction {
	for ; _ne(first, last); first = NextInputIter(first) {
		f(first.Read())
	}
	return f
}

// ForEachN applies the given function f to the result of dereferencing every
// iterator in the range [first, first + n), in order.
func ForEachN(first InputIter, n int, f IteratorFunction) IteratorFunction {
	for ; n > 0; n, first = n-1, NextInputIter(first) {
		f(first.Read())
	}
	return f
}

// Count counts the elements that are equal to value.
func Count(first, last InputIter, v any) int {
	return CountIf(first, last, _eq1(v))
}

// CountIf counts elements for which predicate pred returns true.
func CountIf(first, last InputIter, pred UnaryPredicate) int {
	var ret int
	for ; _ne(first, last); first = NextInputIter(first) {
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
func Mismatch(first1, last1, first2, last2 InputIter) (InputIter, InputIter) {
	return MismatchBy(first1, last1, first2, last2, _eq)
}

// MismatchBy returns the first mismatching pair of elements from two ranges:
// one defined by [first1, last1) and another defined by [first2,last2).
//
// If last2 is nil, it denotes first2 + (last1 - first1). Elements are compared
// using the given comparer eq.
func MismatchBy(first1, last1, first2, last2 InputIter, eq EqComparer) (InputIter, InputIter) {
	for _ne(first1, last1) && (last2 == nil || _ne(first2, last2)) && eq(first1.Read(), first2.Read()) {
		first1, first2 = NextInputIter(first1), NextInputIter(first2)
	}
	return first1, first2
}

// Find returns the first element in the range [first, last) that is equal to
// value.
func Find(first, last InputIter, v any) InputIter {
	return FindIf(first, last, _eq1(v))
}

// FindIf returns the first element in the range [first, last) which predicate
// pred returns true.
func FindIf(first, last InputIter, pred UnaryPredicate) InputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
		if pred(first.Read()) {
			return first
		}
	}
	return last
}

// FindIfNot returns the first element in the range [first, last) which
// predicate pred returns false.
func FindIfNot(first, last InputIter, pred UnaryPredicate) InputIter {
	return FindIf(first, last, _not1(pred))
}

// FindEnd searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last).
//
// If [sFirst, sLast) is empty or such sequence is found, last is returned.
func FindEnd(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindEndBy(first, last, sFirst, sLast, _eq)
}

// FindEndBy searches for the last occurrence of the sequence [sFirst, sLast) in
// the range [first, last).
//
// If [sFirst, sLast) is empty or such sequence is found, last is returned.
// Elements are compared using the given binary comparer eq.
func FindEndBy(first, last, sFirst, sLast ForwardReader, eq EqComparer) ForwardReader {
	if _eq(sFirst, sLast) {
		return last
	}
	result := last
	for {
		if newResult := SearchBy(first, last, sFirst, sLast, eq); _eq(newResult, last) {
			break
		} else {
			result = newResult
			first = NextForwardReader(result)
		}
	}
	return result
}

// FindFirstOf searches the range [first, last) for any of the elements in the
// range [sFirst, sLast).
func FindFirstOf(first, last, sFirst, sLast ForwardReader) ForwardReader {
	return FindFirstOfBy(first, last, sFirst, sLast, _eq)
}

// FindFirstOfBy searches the range [first, last) for any of the elements in the
// range [sFirst, sLast).
//
// Elements are compared using the given binary comparer eq.
func FindFirstOfBy(first, last, sFirst, sLast ForwardReader, eq EqComparer) ForwardReader {
	return FindIf(first, last, func(x any) bool {
		return AnyOf(sFirst, sLast, func(s any) bool {
			return eq(x, s)
		})
	}).(ForwardReader)
}

// AdjacentFind searches the range [first, last) for two consecutive identical
// elements.
func AdjacentFind(first, last ForwardReader) ForwardReader {
	return AdjacentFindBy(first, last, _eq)
}

// AdjacentFindBy searches the range [first, last) for two consecutive identical
// elements.
//
// Elements are compared using the given binary comparer eq.
func AdjacentFindBy(first, last ForwardReader, eq EqComparer) ForwardReader {
	if _eq(first, last) {
		return last
	}
	for next := NextForwardReader(first); _ne(next, last); first, next = NextForwardReader(first), NextForwardReader(next) {
		if eq(first.Read(), next.Read()) {
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
// [sFirst, sLast) in the range [first, last).
//
// Elements are compared using the given binary comparer eq.
func SearchBy(first, last, sFirst, sLast ForwardReader, eq EqComparer) ForwardReader {
	for {
		it := first
		for sIt := sFirst; ; sIt, it = NextForwardReader(sIt), NextForwardReader(it) {
			if _eq(sIt, sLast) {
				return first
			}
			if _eq(it, last) {
				return last
			}
			if !eq(it.Read(), sIt.Read()) {
				break
			}
		}
		first = NextForwardReader(first)
	}
}

// SearchN searches the range [first, last) for the first sequence of count
// identical elements, each equal to the given value.
func SearchN(first, last ForwardReader, count int, v any) ForwardReader {
	return SearchNBy(first, last, count, v, _eq)
}

// SearchNBy searches the range [first, last) for the first sequence of count
// identical elements.
//
// Elements are compared using the given binary comparer eq.
func SearchNBy(first, last ForwardReader, count int, v any, eq EqComparer) ForwardReader {
	if count <= 0 {
		return first
	}
	for ; _ne(first, last); first = NextForwardReader(first) {
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
			if first = NextForwardReader(first); _eq(first, last) {
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
func Copy(first, last InputIter, dFirst OutputIter) OutputIter {
	return CopyIf(first, last, dFirst, _true1)
}

// CopyIf copies the elements in the range, defined by [first, last), and
// predicate pred returns true, to another range beginning at dFirst.
//
// It returns an iterator in the destination range, pointing past the last
// element copied.
func CopyIf(first, last InputIter, dFirst OutputIter, pred UnaryPredicate) OutputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
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
func CopyN(first InputIter, count int, dFirst OutputIter) OutputIter {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, first.Read())
		first = NextInputIter(first)
	}
	return dFirst
}

// CopyBackward copies the elements from the range, defined by [first, last), to
// another range ending at dLast.
//
// The elements are copied in reverse order (the last element is copied first),
// but their relative order is preserved. It returns an iterator to the last
// element copied.
func CopyBackward(first, last BidiReader, dLast BidiWriter) BidiWriter {
	for _ne(first, last) {
		last, dLast = PrevBidiReader(last), PrevBidiWriter(dLast)
		dLast.Write(last.Read())
	}
	return dLast
}

// Fill assigns the given value to the elements in the range [first, last).
func Fill(first, last ForwardWriter, v any) {
	for ; _ne(first, last); first = NextForwardWriter(first) {
		first.Write(v)
	}
}

// FillN assigns the given value to the first count elements in the range
// beginning at dFirst.
//
// If count <= 0, it does nothing.
func FillN(dFirst OutputIter, count int, v any) {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, v)
	}
}

// Transform applies the given function to the range [first, last) and stores
// the result in another range, beginning at dFirst.
func Transform(first, last InputIter, dFirst OutputIter, op UnaryOperation) OutputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
		dFirst = _writeNext(dFirst, op(first.Read()))
	}
	return dFirst
}

// TransformBinary applies the given function to the two ranges [first, last),
// [first2, first2+last-first) and stores the result in another range, beginning
// at dFirst.
func TransformBinary(first1, last1, first2 ForwardReader, dFirst OutputIter, op BinaryOperation) OutputIter {
	for ; _ne(first1, last1); first1, first2 = NextForwardReader(first1), NextForwardReader(first2) {
		dFirst = _writeNext(dFirst, op(first1.Read(), first2.Read()))
	}
	return dFirst
}

// Generate assigns each element in range [first, last) a value generated by the
// given function object g.
func Generate(first, last ForwardWriter, g Generator) {
	for ; _ne(first, last); first = NextForwardWriter(first) {
		first.Write(g())
	}
}

// GenerateN assigns values, generated by given function object g, to the first
// count elements in the range beginning at dFirst.
//
// If count <= 0, it does nothing.
func GenerateN(dFirst OutputIter, count int, g Generator) OutputIter {
	for ; count > 0; count-- {
		dFirst = _writeNext(dFirst, g())
	}
	return dFirst
}

// Remove removes all elements equal to v from the range [first, last) and
// returns a past-the-end iterator for the new end of the range.
func Remove(first, last ForwardReadWriter, v any) ForwardReadWriter {
	return RemoveIf(first, last, _eq1(v))
}

// RemoveIf removes all elements which predicate function returns true from the
// range [first, last) and returns a past-the-end iterator for the new end of
// the range.
func RemoveIf(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	first = FindIf(first, last, pred).(ForwardReadWriter)
	if _ne(first, last) {
		for i := NextForwardReadWriter(first); _ne(i, last); i = NextForwardReadWriter(i) {
			if !pred(i.Read()) {
				first.Write(i.Read())
				first = NextForwardReadWriter(first)
			}
		}
	}
	return first
}

// RemoveCopy copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements equal to v.
//
// Source and destination ranges cannot overlap.
func RemoveCopy(first, last InputIter, dFirst OutputIter, v any) OutputIter {
	return RemoveCopyIf(first, last, dFirst, _eq1(v))
}

// RemoveCopyIf copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements which predicate function returns
// true.
//
// Source and destination ranges cannot overlap.
func RemoveCopyIf(first, last InputIter, dFirst OutputIter, pred UnaryPredicate) OutputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
		if v := first.Read(); !pred(v) {
			dFirst = _writeNext(dFirst, v)
		}
	}
	return dFirst
}

// Replace replaces all elements equal to old with new in the range [first,
// last).
func Replace(first, last ForwardReadWriter, old, new any) {
	ReplaceIf(first, last, _eq1(old), new)
}

// ReplaceIf replaces all elements satisfy pred with new in the range [first,
// last).
func ReplaceIf(first, last ForwardReadWriter, pred UnaryPredicate, v any) {
	for ; _ne(first, last); first = NextForwardReadWriter(first) {
		if pred(first.Read()) {
			first.Write(v)
		}
	}
}

// ReplaceCopy copies the elements from the range [first, last) to another range
// beginning at dFirst replacing all elements equal to old with new.
//
// The source and destination ranges cannot overlap.
func ReplaceCopy(first, last InputIter, dFirst OutputIter, old, new any) OutputIter {
	return ReplaceCopyIf(first, last, dFirst, _eq1(old), new)
}

// ReplaceCopyIf copies the elements from the range [first, last) to another
// range beginning at dFirst replacing all elements satisfy pred with new.
//
// The source and destination ranges cannot overlap.
func ReplaceCopyIf(first, last InputIter, dFirst OutputIter, pred UnaryPredicate, v any) OutputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
		if v0 := first.Read(); pred(v0) {
			dFirst = _writeNext(dFirst, v)
		} else {
			dFirst = _writeNext(dFirst, v0)
		}
	}
	return dFirst
}

// Swap swaps value of two iterators.
func Swap(a, b ReadWriter) {
	va, vb := a.Read(), b.Read()
	a.Write(vb)
	b.Write(va)
}

// SwapRanges exchanges elements between range [first1, last1) and another range
// starting at first2.
func SwapRanges(first1, last1, first2 ForwardReadWriter) {
	for ; _ne(first1, last1); first1, first2 = NextForwardReadWriter(first1), NextForwardReadWriter(first2) {
		Swap(first1, first2)
	}
}

// Reverse reverses the order of the elements in the range [first, last).
func Reverse(first, last BidiReadWriter) {
	for ; _ne(first, last); first = NextBidiReadWriter(first) {
		last = PrevBidiReadWriter(last)
		if _eq(first, last) {
			return
		}
		Swap(first, last)
	}
}

// ReverseCopy copies the elements from the range [first, last) to another range
// beginning at dFirst in such a way that the elements in the new range are in
// reverse order.
func ReverseCopy(first, last BidiReader, dFirst OutputIter) OutputIter {
	for _ne(first, last) {
		last = PrevBidiReader(last)
		dFirst = _writeNext(dFirst, last.Read())
	}
	return dFirst
}

// Rotate performs a left rotation on a range of elements in such a way, that
// the element nFirst becomes the first element of the new range and nFirst - 1
// becomes the last element.
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
		Swap(write, read)
		write, read = NextForwardReadWriter(write), NextForwardReadWriter(read)
	}
	Rotate(write, nextRead, last)
	return write
}

// RotateCopy copies the elements from the range [first, last), to another range
// beginning at dFirst in such a way, that the element nFirst becomes the first
// element of the new range and nFirst - 1 becomes the last element.
func RotateCopy(first, nFirst, last ForwardReader, dFirst OutputIter) OutputIter {
	return Copy(first, nFirst, Copy(nFirst, last, dFirst))
}

// Shuffle reorders the elements in the given range [first, last) such that each
// possible permutation of those elements has equal probability of appearance.
func Shuffle(first, last RandomReadWriter, r *rand.Rand) {
	r.Shuffle(first.Distance(last), func(i, j int) {
		Swap(AdvanceNReadWriter(first, i), AdvanceNReadWriter(first, j))
	})
}

// Sample selects n elements from the sequence [first; last) such that each
// possible sample has equal probability of appearance, and writes those
// selected elements into the output iterator out.
func Sample(first, last ForwardReader, out OutputIter, n int, r *rand.Rand) OutputIter {
	_, rr := first.(RandomReader)
	rout, rw := out.(RandomWriter)
	if !rr && rw {
		return _reservoirSample(first, last, rout, n, r)
	}
	return _selectionSample(first, last, out, n, r)
}

func _selectionSample(first, last ForwardReader, out OutputIter, n int, r *rand.Rand) OutputIter {
	unsampled := Distance(first, last)
	if n > unsampled {
		n = unsampled
	}
	for ; n != 0; first = NextForwardReader(first) {
		if r.Intn(unsampled) < n {
			out = _writeNext(out, first.Read())
			n--
		}
		unsampled--
	}
	return out
}

func _reservoirSample(first, last ForwardReader, out RandomWriter, n int, r *rand.Rand) RandomWriter {
	var k int
	for ; _ne(first, last) && k < n; first, k = NextForwardReader(first), k+1 {
		AdvanceNWriter(out, k).Write(first.Read())
	}
	if _eq(first, last) {
		return AdvanceNWriter(out, k)
	}
	sz := k
	for ; _ne(first, last); first, k = NextForwardReader(first), k+1 {
		if d := r.Intn(k + 1); d < sz {
			AdvanceNWriter(out, d).Write(first.Read())
		}
	}
	return AdvanceNWriter(out, n)
}

// Unique eliminates all but the first element from every consecutive group of
// equivalent elements from the range [first, last) and returns a past-the-end
// iterator for the new logical end of the range.
func Unique(first, last ForwardReadWriter) ForwardReadWriter {
	return UniqueIf(first, last, _eq)
}

// UniqueIf eliminates all but the first element from every consecutive group of
// equivalent elements from the range [first, last) and returns a past-the-end
// iterator for the new logical end of the range.
//
// Elements are compared using the given binary comparer eq.
func UniqueIf(first, last ForwardReadWriter, eq EqComparer) ForwardReadWriter {
	if _eq(first, last) {
		return last
	}
	result := first
	for {
		first = NextForwardReadWriter(first)
		if _eq(first, last) {
			return NextForwardReadWriter(result)
		}
		if !eq(result.Read(), first.Read()) {
			if result = NextForwardReadWriter(result); _ne(result, first) {
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
func UniqueCopy(first, last InputIter, dFirst OutputIter) OutputIter {
	return UniqueCopyIf(first, last, dFirst, _eq)
}

// UniqueCopyIf copies the elements from the range [first, last), to another
// range beginning at dFirst in such a way that there are no consecutive equal
// elements.
//
// Only the first element of each group of equal elements is copied. Elements
// are compared using the given binary comparer eq.
func UniqueCopyIf(first, last InputIter, dFirst OutputIter, eq EqComparer) OutputIter {
	if _ne(first, last) {
		v := first.Read()
		dFirst = _writeNext(dFirst, v)
		for first = NextInputIter(first); _ne(first, last); first = NextInputIter(first) {
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
func IsPartitioned(first, last InputIter, pred UnaryPredicate) bool {
	return NoneOf(FindIfNot(first, last, pred), last, pred)
}

// Partition reorders the elements in the range [first, last) in such a way that
// all elements for which the predicate pred returns true precede the elements
// for which predicate pred returns false.
//
// Relative order of the elements is not preserved.
func Partition(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	first = FindIfNot(first, last, pred).(ForwardReadWriter)
	if _eq(first, last) {
		return first
	}
	for i := NextForwardReadWriter(first); _ne(i, last); i = NextForwardReadWriter(i) {
		if pred(i.Read()) {
			Swap(first, i)
			first = NextForwardReadWriter(first)
		}
	}
	return first
}

// PartitionCopy copies the elements from the range [first, last) to two
// different ranges depending on the value returned by the predicate pred. The
// elements that satisfy the predicate pred are copied to the range beginning at
// outTrue. The rest of the elements are copied to the range beginning at
// outFalse.
func PartitionCopy(first, last InputIter, outTrue, outFalse OutputIter, pred UnaryPredicate) (OutputIter, OutputIter) {
	for ; _ne(first, last); first = NextInputIter(first) {
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
func StablePartition(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	for {
		if _eq(first, last) {
			return first
		}
		if !pred(first.Read()) {
			break
		}
		first = NextForwardReadWriter(first)
	}
	if bfirst, ok := first.(BidiReadWriter); ok {
		if blast, ok := last.(BidiReadWriter); ok {
			for {
				blast = PrevBidiReadWriter(blast)
				if _eq(first, blast) {
					return first
				}
				if pred(blast.Read()) {
					break
				}
			}
			return _stablePartitionBidi(bfirst, blast, pred, Distance(first, blast)+1)
		}
	}
	return _stablePartitionForward(first, last, pred, Distance(first, last))
}

func _stablePartitionBidi(first, last BidiReadWriter, pred UnaryPredicate, l int) BidiReadWriter {
	if l == 2 {
		Swap(first, last)
		return last
	}
	if l == 3 {
		m := NextBidiReadWriter(first)
		if pred(m.Read()) {
			Swap(first, m)
			Swap(m, last)
			return last
		}
		Swap(m, last)
		Swap(first, m)
		return m
	}
	m, l2 := first, l/2
	m = AdvanceN(m, l2).(BidiReadWriter)
	// F???????????????T
	// f       m       l
	m1, lh := m, l2
	for m1 = PrevBidiReadWriter(m1); !pred(m1.Read()); m1 = PrevBidiReadWriter(m1) {
		if _eq(m1, first) {
			break
		}
		lh--
	}
	firstFalse := first
	if _ne(m1, first) {
		// F????TFF????????T
		// f    m1 m       l
		firstFalse = _stablePartitionBidi(first, m1, pred, lh)
	}
	// TTFFFFFF????????T
	// f ff m1 m       l
	m1, lh = m, l-l2
	for pred(m1.Read()) {
		m1 = NextBidiReadWriter(m1)
		if _eq(m1, last) {
			break
		}
		lh--
	}
	secondFalse := NextBidiReadWriter(last)
	if _ne(m1, last) {
		// TTFFFFFFTTTF?????T
		// f ff m1 m  m1    l
		secondFalse = _stablePartitionBidi(m1, last, pred, lh)
	}
	// TTFFFFFFTTTTTTFFFF
	// f ff m1 m  m1 sf l
	return Rotate(firstFalse, m, secondFalse).(BidiReadWriter)
}

func _stablePartitionForward(first, last ForwardReadWriter, pred UnaryPredicate, l int) ForwardReadWriter {
	if l == 1 {
		return first
	}
	if l == 2 {
		m := NextForwardReadWriter(first)
		if pred(m.Read()) {
			Swap(first, m)
			return m
		}
		return first
	}
	l2 := l / 2
	m := AdvanceN(first, l2).(ForwardReadWriter)
	// F?????????????????
	// f       m         l
	firstFalse := _stablePartitionForward(first, m, pred, l2)
	// TTTFFFFF??????????
	// f  ff   m         l
	m1, lh := m, l-l2
	for pred(m1.Read()) {
		m1 = NextForwardReadWriter(m1)
		if _eq(m1, last) {
			break
		}
		lh--
	}
	secondFalse := last
	if _ne(m1, last) {
		// TTTFFFFFTTTF??????
		// f  ff   m  m1     l
		secondFalse = _stablePartitionForward(m1, last, pred, lh)
	}
	// TTTFFFFFTTTTTFFFFF
	// f  ff   m    sf   l
	return Rotate(firstFalse, m, secondFalse)
}

// PartitionPoint examines the partitioned (as if by Partition) range [first,
// last) and locates the end of the first partition, that is, the first element
// that does not satisfy pred or last if all elements satisfy pred.
func PartitionPoint(first, last ForwardReader, pred UnaryPredicate) ForwardReader {
	l := Distance(first, last)
	for l != 0 {
		l2 := l / 2
		m := AdvanceN(first, l2).(ForwardReader)
		if pred(m.Read()) {
			first = NextForwardReader(m)
			l -= l2 + 1
		} else {
			l = l2
		}
	}
	return first
}

// IsSorted checks if the elements in range [first, last) are sorted in
// non-descending order.
func IsSorted(first, last ForwardReader) bool {
	return _eq(IsSortedUntil(first, last), last)
}

// IsSortedBy checks if the elements in range [first, last) are sorted in
// non-descending order.
//
// Elements are compared using the given binary comparer less.
func IsSortedBy(first, last ForwardReader, less LessComparer) bool {
	return _eq(IsSortedUntilBy(first, last, less), last)
}

// IsSortedUntil examines the range [first, last) and finds the largest range
// beginning at first in which the elements are sorted in ascending order.
func IsSortedUntil(first, last ForwardReader) ForwardReader {
	return IsSortedUntilBy(first, last, _less)
}

// IsSortedUntilBy examines the range [first, last) and finds the largest range
// beginning at first in which the elements are sorted in ascending order.
//
// Elements are compared using the given binary comparer less.
func IsSortedUntilBy(first, last ForwardReader, less LessComparer) ForwardReader {
	if _ne(first, last) {
		for next := NextForwardReader(first); _ne(next, last); next = NextForwardReader(next) {
			if less(next.Read(), first.Read()) {
				return next
			}
			first = next
		}
	}
	return last
}

// Adapt RandomIter to sort.Interface.
type sortHelper struct {
	first RandomReadWriter
	n     int
	less  LessComparer
}

func (s *sortHelper) Len() int {
	return s.n
}

func (s *sortHelper) Less(i, j int) bool {
	return s.less(
		AdvanceNReader(s.first, i).Read(),
		AdvanceNReader(s.first, j).Read(),
	)
}

func (s *sortHelper) Swap(i, j int) {
	it1, it2 := AdvanceNReadWriter(s.first, i), AdvanceNReadWriter(s.first, j)
	v1, v2 := it1.Read(), it2.Read()
	it1.Write(v2)
	it2.Write(v1)
}

// Adapt RandomIter to sort.Interface.
type heapHelper struct {
	*sortHelper
}

func (h *heapHelper) Less(i, j int) bool {
	return h.less(
		AdvanceNReader(h.first, j).Read(),
		AdvanceNReader(h.first, i).Read(),
	)
}

func (h *heapHelper) Push(x interface{}) {
	h.n++
}

func (h *heapHelper) Pop() interface{} {
	h.n--
	return nil
}

// Sort sorts the elements in the range [first, last) in ascending order. The
// order of equal elements is not guaranteed to be preserved.
func Sort(first, last RandomReadWriter) {
	SortBy(first, last, _less)
}

// SortBy sorts the elements in the range [first, last) in ascending order. The
// order of equal elements is not guaranteed to be preserved.
//
// Elements are compared using the given binary comparer less.
func SortBy(first, last RandomReadWriter, less LessComparer) {
	sort.Sort(&sortHelper{
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
func PartialSort(first, middle, last RandomReadWriter) {
	PartialSortBy(first, middle, last, _less)
}

// PartialSortBy rearranges elements such that the range [first, middle)
// contains the sorted (middle-first) smallest elements in the range [first,
// last).
//
// The order of equal elements is not guaranteed to be preserved. The order of
// the remaining elements in the range [middle, last) is unspecified. Elements
// are compared using the given binary comparer less.
func PartialSortBy(first, middle, last RandomReadWriter, less LessComparer) {
	MakeHeapBy(first, middle, less)
	for i := middle; _ne(i, last); i = NextRandomReadWriter(i) {
		if less(i.Read(), first.Read()) {
			Swap(first, i)
			heap.Fix(&heapHelper{
				&sortHelper{
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
func PartialSortCopy(first, last InputIter, dFirst, dLast RandomReadWriter) {
	PartialSortCopyBy(first, last, dFirst, dLast, _less)
}

// PartialSortCopyBy sorts some of the elements in the range [first, last) in
// ascending order, storing the result in the range [dFirst, dLast).
//
// At most dLast - dFirst of the elements are placed sorted to the range
// [dFirst, dFirst + n). n is the number of elements to sort (n = min(last -
// first, dLast - dFirst)). The order of equal elements is not guaranteed to be
// preserved. Elements are compared using the given binary comparer less.
func PartialSortCopyBy(first, last InputIter, dFirst, dLast RandomReadWriter, less LessComparer) {
	if _eq(dFirst, dLast) {
		return
	}
	r, len := dFirst, dFirst.Distance(dLast)
	for ; _ne(first, last) && _ne(r, dLast); first, r = NextInputIter(first), NextRandomReadWriter(r) {
		r.Write(first.Read())
	}
	MakeHeapBy(dFirst, dLast, less)
	for ; _ne(first, last); first = NextInputIter(first) {
		if less(first.Read(), dFirst.Read()) {
			dFirst.Write(first.Read())
			heap.Fix(&heapHelper{
				&sortHelper{
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
func StableSort(first, last RandomReadWriter) {
	StableSortBy(first, last, _less)
}

// StableSortBy sorts the elements in the range [first, last) in ascending
// order.
//
// The order of equivalent elements is guaranteed to be preserved. Elements are
// compared using the given binary comparer less.
func StableSortBy(first, last RandomReadWriter, less LessComparer) {
	sort.Stable(&sortHelper{
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
func NthElement(first, nth, last RandomReadWriter) {
	NthElementBy(first, nth, last, _less)
}

// NthElementBy is a partial sorting algorithm that rearranges elements in
// [first, last) such that:
// a. The element pointed at by nth is changed to whatever element would occur
// in that position if [first, last) were sorted.
// b. All of the elements before this new nth element are less than or equal to
// the elements after the new nth element.
//
// Elements are compared using the given binary comparer less.
func NthElementBy(first, nth, last RandomReadWriter, less LessComparer) {
Restart:
	for {
		if _eq(nth, last) {
			return
		}
		len := first.Distance(last)
		if len <= 7 {
			SortBy(first, last, less)
			return
		}

		m := AdvanceNReadWriter(first, len/2)
		last1 := PrevRandomReadWriter(last)

		// sort {first, m, last1}
		var maybeSorted bool
		if !less(m.Read(), first.Read()) {
			// first<=m
			if !less(last1.Read(), m.Read()) {
				// first<=m<=last1
				maybeSorted = true
			} else {
				// first<=m,m>last1
				Swap(m, last1)
				// first<=last1,m<last1
				if less(m.Read(), first.Read()) {
					// m<first<=last1
					Swap(first, m)
					// first<m<=last1
				}
				// first<=m<last1
			}
		} else if less(last1.Read(), m.Read()) {
			// first>m>last1
			Swap(first, last1)
			// first<m<last1
		} else {
			// first>m,m<=last1
			Swap(first, m)
			// first<m,first<=last1
			if less(last1.Read(), m.Read()) {
				// first<=last1<m
				Swap(m, last1)
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
				if j = PrevRandomReadWriter(j); _eq(i, j) {
					// 0+++++++++++++
					// f             l
					// i
					// j
					if i, j = NextRandomReadWriter(i), last1; !less(first.Read(), j.Read()) {
						// 0++++++++++++0
						// fi           jl
						for {
							if _eq(i, j) {
								// 00000000000000
								// f            jl
								//              i
								return
							}
							if less(first.Read(), i.Read()) {
								// 00000>+++++++0
								// f    i       jl
								Swap(i, j)
								maybeSorted = false
								i = NextRandomReadWriter(i)
								break
							}
							i = NextRandomReadWriter(i)
						}
					}
					// 000000+++++++>
					// f     i      jl
					if _eq(i, j) {
						// 0000000000000>
						// f            jl
						//              i
						return
					}
					for {
						for !less(first.Read(), i.Read()) {
							i = NextRandomReadWriter(i)
						}
						for j = PrevRandomReadWriter(j); less(first.Read(), j.Read()); j = PrevRandomReadWriter(j) {
						}
						// 000000>+++++0++
						// f     i     j  l
						if !_less(i, j) {
							break
						}
						Swap(i, j)
						maybeSorted = false
						i = NextRandomReadWriter(i)
					}
					// 000000000+++++++
					// f       ji      l
					if _less(nth, i) {
						return
					}
					first = i
					continue Restart
				}
				if less(j.Read(), m.Read()) {
					// 0???-+++++++++
					// f             l
					// i   j
					Swap(i, j)
					maybeSorted = false
					break
				}
			}
		}

		// i.Read() < m.Read()
		i = NextRandomReadWriter(i)
		// -??????0????????+         -??????0????0+++
		// fi     m        jl  [OR]  fi     m    j   l
		if _less(i, j) {
			for {
				for less(i.Read(), m.Read()) {
					i = NextRandomReadWriter(i)
				}
				for j = PrevRandomReadWriter(j); !less(j.Read(), m.Read()); j = PrevRandomReadWriter(j) {
				}
				// ----+??0?????<+++        -------0?????<+++
				// f      m         l [OR]  f      m         l
				//     i        j                  i     j
				if !less(i, j) {
					// -------0--+++++++       -------0+++++++++
					// f      m         l [OR] f      m         l
					//          ji                   ji
					break
				}
				Swap(i, j)
				maybeSorted = false
				if _eq(m, i) {
					m = j
				}
				// -----??0?????++++       --------?????0+++
				// f   i  m     j   l [OR] f            m   l
				//                                i     j
				i = NextRandomReadWriter(i)
			}
		}

		// -------+++0+++
		// f      i  m   l
		if _ne(i, m) && less(m.Read(), i.Read()) {
			Swap(i, m)
			maybeSorted = false
		}
		// -------0++++++
		// f      i      l
		if _eq(nth, i) {
			return
		}
		if _less(nth, i) {
			if maybeSorted && IsSortedBy(first, i, less) {
				return
			}
			last = i
		} else {
			if maybeSorted && IsSortedBy(i, last, less) {
				return
			}
			first = NextRandomReadWriter(i)
		}
	}
}

// LowerBound returns an iterator pointing to the first element in the range
// [first, last) that is not less than (i.e. greater or equal to) value, or last
// if no such element is found.
func LowerBound(first, last ForwardReader, v any) ForwardReader {
	return LowerBoundBy(first, last, v, _less)
}

// LowerBoundBy returns an iterator pointing to the first element in the range
// [first, last) that is not less than (i.e. greater or equal to) value, or last
// if no such element is found.
//
// Elements are compared using the given binary comparer less.
func LowerBoundBy(first, last ForwardReader, v any, less LessComparer) ForwardReader {
	for len := Distance(first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN(first, l2).(ForwardReader)
		if less(m.Read(), v) {
			first = NextForwardReader(m)
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
func UpperBound(first, last ForwardReader, v any) ForwardReader {
	return UpperBoundBy(first, last, v, _less)
}

// UpperBoundBy returns an iterator pointing to the first element in the range
// [first, last) that is greater than value, or last if no such element is
// found.
//
// Elements are compared using the given binary comparer less.
func UpperBoundBy(first, last ForwardReader, v any, less LessComparer) ForwardReader {
	for len := Distance(first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN(first, l2).(ForwardReader)
		if less(v, m.Read()) {
			len = l2
		} else {
			first = NextForwardReader(m)
			len -= l2 + 1
		}
	}
	return first
}

// BinarySearch checks if an element equivalent to value appears within the
// range [first, last).
func BinarySearch(first, last ForwardReader, v any) bool {
	return BinarySearchBy(first, last, v, _less)
}

// BinarySearchBy checks if an element equivalent to value appears within the
// range [first, last).
//
// Elements are compared using the given binary comparer less.
func BinarySearchBy(first, last ForwardReader, v any, less LessComparer) bool {
	first = LowerBoundBy(first, last, v, less)
	return _ne(first, last) && !(less(v, first.Read()))
}

// EqualRange returns a range containing all elements equivalent to value in the
// range [first, last).
func EqualRange(first, last ForwardReader, v any) (ForwardReader, ForwardReader) {
	return EqualRangeBy(first, last, v, _less)
}

// EqualRangeBy returns a range containing all elements equivalent to value in
// the range [first, last).
//
// Elements are compared using the given binary comparer less.
func EqualRangeBy(first, last ForwardReader, v any, less LessComparer) (ForwardReader, ForwardReader) {
	for len := Distance(first, last); len != 0; {
		l2 := len / 2
		m := AdvanceN(first, l2).(ForwardReader)
		if less(m.Read(), v) {
			first = NextForwardReader(m)
			len -= l2 + 1
		} else if less(v, m.Read()) {
			last = m
			len = l2
		} else {
			return LowerBoundBy(first, m, v, less), UpperBoundBy(NextForwardReader(m), last, v, less)
		}
	}
	return first, first
}

// Merge merges two sorted ranges [first1, last1) and [first2, last2) into one
// sorted range beginning at dFirst.
func Merge(first1, last1, first2, last2 InputIter, dFirst OutputIter) OutputIter {
	return MergeBy(first1, last1, first2, last2, dFirst, _less)
}

// MergeBy merges two sorted ranges [first1, last1) and [first2, last2) into one
// sorted range beginning at dFirst.
//
// Elements are compared using the given binary comparer less.
func MergeBy(first1, last1, first2, last2 InputIter, dFirst OutputIter, less LessComparer) OutputIter {
	for _ne(first1, last1) {
		if _eq(first2, last2) {
			return Copy(first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v2, v1) {
			dFirst = _writeNext(dFirst, v2)
			first2 = NextInputIter(first2)
		} else {
			dFirst = _writeNext(dFirst, v1)
			first1 = NextInputIter(first1)
		}
	}
	return Copy(first2, last2, dFirst)
}

// InplaceMerge Merges two consecutive sorted ranges [first, middle) and
// [middle, last) into one sorted range [first, last). For equivalent elements
// in the original two ranges, the elements from the first range (preserving
// their original order) precede the elements from the second range (preserving
// their original order).
func InplaceMerge(first, middle, last BidiReadWriter) {
	InplaceMergeBy(first, middle, last, _less)
}

// InplaceMergeBy Merges two consecutive sorted ranges [first, middle) and
// [middle, last) into one sorted range [first, last). For equivalent elements
// in the original two ranges, the elements from the first range (preserving
// their original order) precede the elements from the second range (preserving
// their original order).
//
// Elements are compared using the given binary comparer less.
func InplaceMergeBy(first, middle, last BidiReadWriter, less LessComparer) {
	len1, len2 := Distance(first, middle), Distance(middle, last)
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
			first = NextBidiReadWriter(first)
			len1--
		}
		var len11, len21 int
		var m1, m2 BidiReadWriter
		if len1 < len2 {
			len21 = len2 / 2
			m2 = AdvanceN(middle, len21).(BidiReadWriter)
			m1 = UpperBoundBy(first, middle, m2.Read(), less).(BidiReadWriter)
			len11 = Distance(first, m1)
		} else {
			if len1 == 1 {
				Swap(first, middle)
				return
			}
			len11 = len1 / 2
			m1 = AdvanceN(first, len11).(BidiReadWriter)
			m2 = LowerBoundBy(middle, last, m1.Read(), less).(BidiReadWriter)
			len21 = Distance(middle, m2)
		}
		len12, len22 := len1-len11, len2-len21
		middle = Rotate(m1, middle, m2).(BidiReadWriter)
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
func Includes(first1, last1, first2, last2 InputIter) bool {
	return IncludesBy(first1, last1, first2, last2, _less)
}

// IncludesBy returns true if the sorted range [first2, last2) is a subsequence
// of the sorted range [first1, last1). (A subsequence need not be contiguous.)
//
// Elements are compared using the given binary comparer less.
func IncludesBy(first1, last1, first2, last2 InputIter, less LessComparer) bool {
	for ; _ne(first2, last2); first1 = NextInputIter(first1) {
		if _eq(first1, last1) || less(first2.Read(), first1.Read()) {
			return false
		}
		if !less(first1.Read(), first2.Read()) {
			first2 = NextInputIter(first2)
		}
	}
	return true
}

// SetDifference copies the elements from the sorted range [first1, last1) which
// are not found in the sorted range [first2, last2) to the range beginning at
// dFirst.
func SetDifference(first1, last1, first2, last2 InputIter, dFirst OutputIter) OutputIter {
	return SetDifferenceBy(first1, last1, first2, last2, dFirst, _less)
}

// SetDifferenceBy copies the elements from the sorted range [first1, last1)
// which are not found in the sorted range [first2, last2) to the range
// beginning at dFirst.
//
// Elements are compared using the given binary comparer less.
func SetDifferenceBy(first1, last1, first2, last2 InputIter, dFirst OutputIter, less LessComparer) OutputIter {
	for _ne(first1, last1) {
		if _eq(first2, last2) {
			return Copy(first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			dFirst = _writeNext(dFirst, v1)
			first1 = NextInputIter(first1)
		} else {
			if !less(v2, v1) {
				first1 = NextInputIter(first1)
			}
			first2 = NextInputIter(first2)
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
func SetIntersection(first1, last1, first2, last2 InputIter, dFirst OutputIter) OutputIter {
	return SetIntersectionBy(first1, last1, first2, last2, dFirst, _less)
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
func SetIntersectionBy(first1, last1, first2, last2 InputIter, dFirst OutputIter, less LessComparer) OutputIter {
	for _ne(first1, last1) && _ne(first2, last2) {
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			first1 = NextInputIter(first1)
		} else {
			if !less(v2, v1) {
				dFirst = _writeNext(dFirst, v1)
				first1 = NextInputIter(first1)
			}
			first2 = NextInputIter(first2)
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
func SetSymmetricDifference(first1, last1, first2, last2 InputIter, dFirst OutputIter) OutputIter {
	return SetSymmetricDifferenceBy(first1, last1, first2, last2, dFirst, _less)
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
func SetSymmetricDifferenceBy(first1, last1, first2, last2 InputIter, dFirst OutputIter, less LessComparer) OutputIter {
	for _ne(first1, last1) {
		if _eq(first2, last2) {
			return Copy(first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v1, v2) {
			dFirst = _writeNext(dFirst, v1)
			first1 = NextInputIter(first1)
		} else {
			if less(v2, v1) {
				dFirst = _writeNext(dFirst, v2)
			} else {
				first1 = NextInputIter(first1)
			}
			first2 = NextInputIter(first2)
		}
	}
	return Copy(first2, last2, dFirst)
}

// SetUnion constructs a sorted union beginning at dFirst consisting of the set
// of elements present in one or both sorted ranges [first1, last1) and [first2,
// last2).
//
// If some element is found m times in [first1, last1) and n times in [first2,
// last2), then all m elements will be copied from [first1, last1) to dFirst,
// preserving order, and then exactly Max(n-m, 0) elements will be copied from
// [first2, last2) to dFirst, also preserving order.
func SetUnion(first1, last1, first2, last2 InputIter, dFirst OutputIter) OutputIter {
	return SetUnionBy(first1, last1, first2, last2, dFirst, _less)
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
func SetUnionBy(first1, last1, first2, last2 InputIter, dFirst OutputIter, less LessComparer) OutputIter {
	for _ne(first1, last1) {
		if _eq(first2, last2) {
			return Copy(first1, last1, dFirst)
		}
		if v1, v2 := first1.Read(), first2.Read(); less(v2, v1) {
			dFirst = _writeNext(dFirst, v2)
			first2 = NextInputIter(first2)
		} else {
			dFirst = _writeNext(dFirst, v1)
			if !less(v1, v2) {
				first2 = NextInputIter(first2)
			}
			first1 = NextInputIter(first1)
		}
	}
	return Copy(first2, last2, dFirst)
}

// IsHeap checks if the elements in range [first, last) are a max heap.
func IsHeap(first, last RandomReader) bool {
	return IsHeapBy(first, last, _less)
}

// IsHeapBy checks if the elements in range [first, last) are a max heap.
//
// Elements are compared using the given binary comparer less.
func IsHeapBy(first, last RandomReader, less LessComparer) bool {
	return _eq(IsHeapUntilBy(first, last, less), last)
}

// IsHeapUntil examines the range [first, last) and finds the largest range
// beginning at first which is a max heap.
func IsHeapUntil(first, last RandomReader) RandomReader {
	return IsHeapUntilBy(first, last, _less)
}

// IsHeapUntilBy examines the range [first, last) and finds the largest range
// beginning at first which is a max heap.
//
// Elements are compared using the given binary comparer less.
func IsHeapUntilBy(first, last RandomReader, less LessComparer) RandomReader {
	len, p, c, pp := first.Distance(last), 0, 1, first
	for c < len {
		cp := AdvanceNReader(first, c)
		if less(pp.Read(), cp.Read()) {
			return cp
		}
		c++
		if c == len {
			return last
		}
		cp = NextRandomReader(cp)
		if less(pp.Read(), cp.Read()) {
			return cp
		}
		p++
		pp = NextRandomReader(pp)
		c = 2*p + 1
	}
	return last
}

// MakeHeap constructs a max heap in the range [first, last).
func MakeHeap(first, last RandomReadWriter) {
	MakeHeapBy(first, last, _less)
}

// MakeHeapBy constructs a max heap in the range [first, last).
//
// Elements are compared using the given binary comparer less.
func MakeHeapBy(first, last RandomReadWriter, less LessComparer) {
	heap.Init(&heapHelper{
		&sortHelper{
			first: first,
			n:     first.Distance(last),
			less:  less,
		}})
}

// PushHeap inserts the element at the position last-1 into the max heap defined
// by the range [first, last-1).
func PushHeap(first, last RandomReadWriter) {
	PushHeapBy(first, last, _less)
}

// PushHeapBy inserts the element at the position last-1 into the max heap
// defined by the range [first, last-1).
//
// Elements are compared using the given binary comparer less.
func PushHeapBy(first, last RandomReadWriter, less LessComparer) {
	heap.Push(&heapHelper{
		&sortHelper{
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
func PopHeap(first, last RandomReadWriter) {
	PopHeapBy(first, last, _less)
}

// PopHeapBy swaps the value in the position first and the value in the position
// last-1 and makes the subrange [first, last-1) into a heap. This has the
// effect of removing the first element from the heap defined by the range
// [first, last).
//
// Elements are compared using the given binary comparer less.
func PopHeapBy(first, last RandomReadWriter, less LessComparer) {
	heap.Pop(&heapHelper{
		&sortHelper{
			first: first,
			n:     first.Distance(last),
			less:  less,
		}})
}

// SortHeap converts the max heap [first, last) into a sorted range in ascending
// order. The resulting range no longer has the heap property.
func SortHeap(first, last RandomReadWriter) {
	SortHeapBy(first, last, _less)
}

// SortHeapBy converts the max heap [first, last) into a sorted range in ascending
// order. The resulting range no longer has the heap property.
//
// Elements are compared using the given binary comparer less.
func SortHeapBy(first, last RandomReadWriter, less LessComparer) {
	for ; _ne(first, last); last = PrevRandomReadWriter(last) {
		PopHeapBy(first, last, less)
	}
}

// Max returns the greater of the given values.
func Max(a, b any) any {
	return MaxBy(a, b, _less)
}

// MaxBy returns the greater of the given values.
//
// Values are compared using the given binary comparer less.
func MaxBy(a, b any, less LessComparer) any {
	if less(a, b) {
		return b
	}
	return a
}

// MaxElement returns the largest element in a range.
func MaxElement(first, last ForwardReader) ForwardReader {
	return MaxElementBy(first, last, _less)
}

// MaxElementBy returns the largest element in a range.
//
// Values are compared using the given binary comparer less.
func MaxElementBy(first, last ForwardReader, less LessComparer) ForwardReader {
	if _eq(first, last) {
		return last
	}
	max := first
	for first = NextForwardReader(first); _ne(first, last); first = NextForwardReader(first) {
		if less(max.Read(), first.Read()) {
			max = first
		}
	}
	return max
}

// Min returns the smaller of the given values.
func Min(a, b any) any {
	return MinBy(a, b, _less)
}

// MinBy returns the smaller of the given values.
//
// Values are compared using the given binary comparer less.
func MinBy(a, b any, less LessComparer) any {
	if less(a, b) {
		return a
	}
	return b
}

// MinElement returns the smallest element in a range.
func MinElement(first, last ForwardReader) ForwardReader {
	return MinElementBy(first, last, _less)
}

// MinElementBy returns the smallest element in a range.
//
// Values are compared using the given binary comparer less.
func MinElementBy(first, last ForwardReader, less LessComparer) ForwardReader {
	if _eq(first, last) {
		return last
	}
	min := first
	for first = NextForwardReader(first); _ne(first, last); first = NextForwardReader(first) {
		if less(first.Read(), min.Read()) {
			min = first
		}
	}
	return min
}

// Minmax returns the smaller and larger of two elements.
func Minmax(a, b any) (any, any) {
	return MinmaxBy(a, b, _less)
}

// MinmaxBy returns the smaller and larger of two elements.
//
// Values are compared using the given binary comparer less.
func MinmaxBy(a, b any, less LessComparer) (any, any) {
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
//
// Values are compared using the given binary comparer less.
func MinmaxElementBy(first, last ForwardReader, less LessComparer) (ForwardReader, ForwardReader) {
	if _eq(first, last) {
		return first, first
	}
	min, max := first, first
	for first = NextForwardReader(first); _ne(first, last); first = NextForwardReader(first) {
		i := first
		first = NextForwardReader(first)
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
func Clamp(v, lo, hi any) any {
	return ClampBy(v, lo, hi, _less)
}

// ClampBy clamps a value between a pair of boundary values.
//
// Values are compared using the given binary comparer less.
func ClampBy(v, lo, hi any, less LessComparer) any {
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
func Equal(first1, last1, first2, last2 InputIter) bool {
	return EqualBy(first1, last1, first2, last2, _eq)
}

// EqualBy returns true if the range [first1, last1) is equal to the range
// [first2, last2), and false otherwise.
//
// If last2 is nil, it denotes first2 + (last1 - first1). Elements are compared
// using the given binary comparer eq.
func EqualBy(first1, last1, first2, last2 InputIter, eq EqComparer) bool {
	for ; _ne(first1, last1); first1, first2 = NextInputIter(first1), NextInputIter(first2) {
		if (last2 != nil && _eq(first2, last2)) || !eq(first1.Read(), first2.Read()) {
			return false
		}
	}
	return last2 == nil || _eq(first2, last2)
}

// LexicographicalCompare checks if the first range [first1, last1) is
// lexicographically less than the second range [first2, last2).
func LexicographicalCompare(first1, last1, first2, last2 InputIter) bool {
	return LexicographicalCompareBy(first1, last1, first2, last2, _less)
}

// LexicographicalCompareBy checks if the first range [first1, last1) is
// lexicographically less than the second range [first2, last2).
//
// Elements are compared using the given binary comparer less.
func LexicographicalCompareBy(first1, last1, first2, last2 InputIter, less LessComparer) bool {
	for ; _ne(first2, last2); first1, first2 = NextInputIter(first1), NextInputIter(first2) {
		if _eq(first1, last1) || less(first1.Read(), first2.Read()) {
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
func LexicographicalCompareThreeWay(first1, last1, first2, last2 InputIter) int {
	return LexicographicalCompareThreeWayBy(first1, last1, first2, last2, _cmp)
}

// LexicographicalCompareThreeWayBy lexicographically compares two ranges [first1,
// last1) and [first2, last2) using three-way comparison. The result will be 0
// if [first1, last1) == [first2, last2), -1 if [first1, last1) < [first2,
// last2), 1 if [first1, last1) > [first2, last2).
//
// Elements are compared using the given binary predicate cmp.
func LexicographicalCompareThreeWayBy(first1, last1, first2, last2 InputIter, cmp ThreeWayComparer) int {
	for ; _ne(first2, last2); first1, first2 = NextInputIter(first1), NextInputIter(first2) {
		if _eq(first1, last1) {
			return -1
		}
		if x := cmp(first1.Read(), first2.Read()); x != 0 {
			return x
		}
	}
	if _eq(first1, last1) {
		return 0
	}
	return 1
}

// IsPermutation returns true if there exists a permutation of the elements in
// the range [first1, last1) that makes that range equal to the range
// [first2,last2), where last2 denotes first2 + (last1 - first1) if it was not
// given.
func IsPermutation(first1, last1, first2, last2 ForwardReader) bool {
	return IsPermutationBy(first1, last1, first2, last2, _eq)
}

// IsPermutationBy returns true if there exists a permutation of the elements in
// the range [first1, last1) that makes that range equal to the range
// [first2,last2), where last2 denotes first2 + (last1 - first1) if it was not
// given.
//
// Elements are compared using the given binary comparer eq.
func IsPermutationBy(first1, last1, first2, last2 ForwardReader, eq EqComparer) bool {
	l := Distance(first1, last1)
	if last2 == nil {
		last2 = AdvanceN(first2, l).(ForwardReader)
	} else if Distance(first2, last2) != l {
		return false
	}
	m1, m2 := MismatchBy(first1, last1, first2, last2, eq)
	first1, first2 = m1.(ForwardReader), m2.(ForwardReader)
	if _eq(first1, last1) {
		return true
	}
	for i := first1; _ne(i, last1); i = NextForwardReader(i) {
		pred := _eq1(i.Read())
		if _ne(FindIf(first1, i, pred), i) {
			continue
		}
		c2 := CountIf(first2, last2, pred)
		if c2 == 0 || c2 != 1+CountIf(NextForwardReader(i), last1, pred) {
			return false
		}
	}
	return true
}

// NextPermutation transforms the range [first, last) into the next permutation
// from the set of all permutations that are lexicographically ordered. Returns
// true if such permutation exists, otherwise transforms the range into the
// first permutation (as if by Sort(first, last)) and returns false.
func NextPermutation(first, last BidiReadWriter) bool {
	return NextPermutationBy(first, last, _less)
}

// NextPermutationBy transforms the range [first, last) into the next
// permutation from the set of all permutations that are lexicographically
// ordered with respect to less. Returns true if such permutation exists,
// otherwise transforms the range into the first permutation (as if by
// Sort(first, last)) and returns false.
//
// Elements are compared using the given
// binary comparer less.
func NextPermutationBy(first, last BidiReadWriter, less LessComparer) bool {
	if _eq(first, last) {
		return false
	}
	i := PrevBidiReadWriter(last)
	if _eq(first, i) {
		return false
	}
	for {
		ip1 := i
		i = PrevBidiReadWriter(i)
		if less(i.Read(), ip1.Read()) {
			j := PrevBidiReadWriter(last)
			for ; !less(i.Read(), j.Read()); j = PrevBidiReadWriter(j) {
			}
			Swap(i, j)
			Reverse(ip1, last)
			return true
		}
		if _eq(i, first) {
			Reverse(first, last)
			return false
		}
	}
}

// PrevPermutation transforms the range [first, last) into the previous
// permutation from the set of all permutations that are lexicographically
// ordered. Returns true if such permutation exists, otherwise transforms the
// range into the last permutation (as if by Sort(first, last); Reverse(first,
// last);) and returns false.
func PrevPermutation(first, last BidiReadWriter) bool {
	return PrevPermutationBy(first, last, _less)
}

// PrevPermutationBy transforms the range [first, last) into the previous
// permutation from the set of all permutations that are lexicographically
// ordered with respect to less. Returns true if such permutation exists,
// otherwise transforms the range into the last permutation (as if by
// Sort(first, last); Reverse(first, last);) and returns false.
//
// Elements are compared using the given binary comparer less.
func PrevPermutationBy(first, last BidiReadWriter, less LessComparer) bool {
	if _eq(first, last) {
		return false
	}
	i := PrevBidiReadWriter(last)
	if _eq(first, i) {
		return false
	}
	for {
		ip1 := i
		i = PrevBidiReadWriter(i)
		if less(ip1.Read(), i.Read()) {
			j := PrevBidiReadWriter(last)
			for ; !less(j.Read(), i.Read()); j = PrevBidiReadWriter(j) {
			}
			Swap(i, j)
			Reverse(ip1, last)
			return true
		}
		if _eq(i, first) {
			Reverse(first, last)
			return false
		}
	}
}

// Iota fills the range [first, last) with sequentially increasing values,
// starting with v and repetitively evaluating v++/v.Inc().
func Iota(first, last ForwardWriter, v any) {
	IotaBy(first, last, v, _inc)
}

// IotaBy fills the range [first, last) with sequentially increasing values,
// starting with v and repetitively evaluating inc(v).
func IotaBy(first, last ForwardWriter, v any, inc UnaryOperation) {
	for ; _ne(first, last); first, v = NextForwardWriter(first), inc(v) {
		first.Write(v)
	}
}

// Accumulate computes the sum of the given value v and the elements in the
// range [first, last), using v+=x or v=v.Add(x).
func Accumulate(first, last InputIter, v any) any {
	return AccumulateBy(first, last, v, _add)
}

// AccumulateBy computes the sum of the given value v and the elements in the
// range [first, last), using v=add(v,x).
func AccumulateBy(first, last InputIter, v any, add BinaryOperation) any {
	for ; _ne(first, last); first = NextInputIter(first) {
		v = add(v, first.Read())
	}
	return v
}

// InnerProduct computes inner product (i.e. sum of products) or performs
// ordered map/reduce operation on the range [first1, last1), using v=v+x*y or
// v=v.Add(x.Mul(y)).
func InnerProduct(first1, last1, first2 InputIter, v any) any {
	return InnerProductBy(first1, last1, first2, v, _add, _mul)
}

// InnerProductBy computes inner product (i.e. sum of products) or performs
// ordered map/reduce operation on the range [first1, last1), using
// v=add(v,mul(x,y)).
func InnerProductBy(first1, last1, first2 InputIter, v any, add, mul BinaryOperation) any {
	for ; _ne(first1, last1); first1, first2 = NextInputIter(first1), NextInputIter(first2) {
		v = add(v, mul(first1.Read(), first2.Read()))
	}
	return v
}

// AdjacentDifference computes the differences between the second and the first
// of each adjacent pair of elements of the range [first, last) and writes them
// to the range beginning at dFirst + 1. An unmodified copy of first is
// written to dFirst. Differences are calculated by cur-prev or cur.Sub(prev).
func AdjacentDifference(first, last InputIter, dFirst OutputIter) OutputIter {
	return AdjacentDifferenceBy(first, last, dFirst, _sub)
}

// AdjacentDifferenceBy computes the differences between the second and the
// first of each adjacent pair of elements of the range [first, last) and writes
// them to the range beginning at dFirst + 1. An unmodified copy of first is
// written to dFirst. Differences are calculated by sub(cur,prev).
func AdjacentDifferenceBy(first, last InputIter, dFirst OutputIter, sub BinaryOperation) OutputIter {
	if _eq(first, last) {
		return dFirst
	}
	prev := first.Read()
	dFirst = _writeNext(dFirst, prev)
	for first = NextInputIter(first); _ne(first, last); first = NextInputIter(first) {
		cur := first.Read()
		dFirst = _writeNext(dFirst, sub(cur, prev))
		prev = cur
	}
	return dFirst
}

// PartialSum computes the partial sums of the elements in the subranges of the
// range [first, last) and writes them to the range beginning at dFirst. Sums
// are calculated by sum=sum+cur or sum=sum.Add(cur).
func PartialSum(first, last InputIter, dFirst OutputIter) OutputIter {
	return PartialSumBy(first, last, dFirst, _add)
}

// PartialSumBy computes the partial sums of the elements in the subranges of
// the range [first, last) and writes them to the range beginning at dFirst.
// Sums are calculated by sum=add(sum,cur).
func PartialSumBy(first, last InputIter, dFirst OutputIter, add BinaryOperation) OutputIter {
	if _eq(first, last) {
		return dFirst
	}
	sum := first.Read()
	dFirst = _writeNext(dFirst, sum)
	for first = NextInputIter(first); _ne(first, last); first = NextInputIter(first) {
		sum = add(sum, first.Read())
		dFirst = _writeNext(dFirst, sum)
	}
	return dFirst
}

// ExclusiveScan computes an exclusive prefix sum operation using v=v+cur or
// v=v.Add(cur) for the range [first, last), using v as the initial value, and
// writes the results to the range beginning at dFirst. "exclusive" means that
// the i-th input element is not included in the i-th sum.
func ExclusiveScan(first, last InputIter, dFirst OutputIter, v any) OutputIter {
	return ExclusiveScanBy(first, last, dFirst, v, _add)
}

// ExclusiveScanBy computes an exclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value, and writes the
// results to the range beginning at dFirst. "exclusive" means that the i-th
// input element is not included in the i-th sum.
func ExclusiveScanBy(first, last InputIter, dFirst OutputIter, v any, add BinaryOperation) OutputIter {
	return TransformExclusiveScanBy(first, last, dFirst, v, add, _noop)
}

// InclusiveScan computes an inclusive prefix sum operation using v=v+cur or
// v=v.Add(cur) for the range [first, last), using v as the initial value (if
// provided), and writes the results to the range beginning at dFirst.
// "inclusive" means that the i-th input element is included in the i-th sum.
func InclusiveScan(first, last InputIter, dFirst OutputIter, v any) OutputIter {
	return InclusiveScanBy(first, last, dFirst, v, _add)
}

// InclusiveScanBy computes an inclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value (if provided), and
// writes the results to the range beginning at dFirst. "inclusive" means that
// the i-th input element is included in the i-th sum.
func InclusiveScanBy(first, last InputIter, dFirst OutputIter, v any, add BinaryOperation) OutputIter {
	return TransformInclusiveScanBy(first, last, dFirst, v, add, _noop)
}

// TransformExclusiveScan transforms each element in the range [first, last)
// with op, then computes an exclusive prefix sum operation using v=v+cur or
// v=v.Add(cur) for the range [first, last), using v as the initial value, and
// writes the results to the range beginning at dFirst. "exclusive" means that
// the i-th input element is not included in the i-th sum.
func TransformExclusiveScan(first, last InputIter, dFirst OutputIter, v any, op UnaryOperation) OutputIter {
	return TransformExclusiveScanBy(first, last, dFirst, v, _add, op)
}

// TransformExclusiveScanBy transforms each element in the range [first, last)
// with op, then computes an exclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value, and writes the
// results to the range beginning at dFirst. "exclusive" means that the i-th
// input element is not included in the i-th sum.
func TransformExclusiveScanBy(first, last InputIter, dFirst OutputIter, v any, add BinaryOperation, op UnaryOperation) OutputIter {
	if _eq(first, last) {
		return dFirst
	}
	saved := v
	for {
		v = add(v, op(first.Read()))
		dFirst = _writeNext(dFirst, saved)
		saved = v
		first = NextInputIter(first)
		if _eq(first, last) {
			break
		}
	}
	return dFirst
}

// TransformInclusiveScan transforms each element in the range [first, last)
// with op, then computes an inclusive prefix sum operation using v=v+cur or
// v=v.Add(cur) for the range [first, last), using v as the initial value (if
// provided), and writes the results to the range beginning at dFirst.
// "inclusive" means that the i-th input element is included in the i-th sum.
func TransformInclusiveScan(first, last InputIter, dFirst OutputIter, v any, op UnaryOperation) OutputIter {
	return TransformInclusiveScanBy(first, last, dFirst, v, _add, op)
}

// TransformInclusiveScanBy transforms each element in the range [first, last)
// with op, then computes an inclusive prefix sum operation using v=add(v,cur)
// for the range [first, last), using v as the initial value (if provided), and
// writes the results to the range beginning at dFirst. "inclusive" means that
// the i-th input element is included in the i-th sum.
func TransformInclusiveScanBy(first, last InputIter, dFirst OutputIter, v any, add BinaryOperation, op UnaryOperation) OutputIter {
	for ; _ne(first, last); first = NextInputIter(first) {
		v = add(v, op(first.Read()))
		dFirst = _writeNext(dFirst, v)
	}
	return dFirst
}
