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

// FindEndBy searches for the last occurrence of the sequence [sFirst, sLast) in
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

// Copy copies the elements in the range, defined by [first, last), to another
// range beginning at dFirst. It returns an iterator in the destination range,
// pointing past the last element copied.
func Copy(first, last ForwardReader, dFirst ForwardWriter) ForwardWriter {
	return CopyIf(first, last, dFirst, _true1)
}

// CopyIf copies the elements in the range, defined by [first, last), and
// predicate pred returns true, to another range beginning at dFirst. It returns
// an iterator in the destination range, pointing past the last element copied.
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

// CopyN copies exactly count values from the range beginning at first to the
// range beginning at result. It returns an iterator in the destination range,
// pointing past the last element copied.
func CopyN(first ForwardReader, count int, result ForwardWriter) ForwardWriter {
	for ; count > 0; count-- {
		result.Write(first.Read())
		first, result = NextReader(first), NextWriter(result)
	}
	return result
}

// CopyBackward copies the elements from the range, defined by [first, last), to
// another range ending at dLast. The elements are copied in reverse order (the
// last element is copied first), but their relative order is preserved. It
// returns an iterator to the last element copied.
func CopyBackward(first, last BidiReader, dLast BidiWriter) BidiWriter {
	for _ne(first, last) {
		last, dLast = PrevBidiReader(last), PrevBidiWriter(dLast)
		dLast.Write(last.Read())
	}
	return dLast
}

// Fill assigns the given value to the elements in the range [first, last).
func Fill(first, last ForwardWriter, v Any) {
	for ; _ne(first, last); first = NextWriter(first) {
		first.Write(v)
	}
}

// FillN assigns the given value to the first count elements in the range
// beginning at first if count > 0. Does nothing otherwise.
func FillN(first ForwardWriter, count int, v Any) {
	for ; count > 0; count-- {
		first.Write(v)
		first = NextWriter(first)
	}
}

// Transform applies the given function to the range [first, last) and stores
// the result in another range, beginning at dFirst.
func Transform(first, last ForwardReader, dFirst ForwardWriter, op UnaryOperation) ForwardWriter {
	for ; _ne(first, last); dFirst, first = NextWriter(dFirst), NextReader(first) {
		dFirst.Write(op(first.Read()))
	}
	return dFirst
}

// TransformBinary applies the given function to the two ranges [first, last),
// [first2, first2+last-first) and stores the result in another range, beginning
// at dFirst.
func TransformBinary(first1, last1, first2 ForwardReader, dFirst ForwardWriter, op BinaryOperation) ForwardWriter {
	for ; _ne(first1, last1); dFirst, first1, first2 = NextWriter(dFirst), NextReader(first1), NextReader(first2) {
		dFirst.Write(op(first1.Read(), first2.Read()))
	}
	return dFirst
}

// Generate assigns each element in range [first, last) a value generated by the
// given function object g.
func Generate(first, last ForwardWriter, g Generator) {
	for ; _ne(first, last); first = NextWriter(first) {
		first.Write(g())
	}
}

// GenerateN assigns values, generated by given function object g, to the first
// count elements in the range beginning at first, if count>0. Does nothing
// otherwise.
func GenerateN(first ForwardWriter, count int, g Generator) ForwardWriter {
	for ; count > 0; count-- {
		first.Write(g())
		first = NextWriter(first)
	}
	return first
}

// Remove removes all elements equal to v from the range [first, last) and
// returns a past-the-end iterator for the new end of the range.
func Remove(first, last ForwardReadWriter, v Any) ForwardReadWriter {
	return RemoveIf(first, last, _eq1(v))
}

// RemoveIf removes all elements which predicate function returns true from the
// range [first, last) and returns a past-the-end iterator for the new end of
// the range.
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

// RemoveCopy copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements equal to v. Source and destination
// ranges cannot overlap.
func RemoveCopy(first, last ForwardReader, dFirst ForwardWriter, v Any) ForwardWriter {
	return RemoveCopyIf(first, last, dFirst, _eq1(v))
}

// RemoveCopyIf copies elements from the range [first, last), to another range
// beginning at dFirst, omitting the elements which predicate function returns
// true. Source and destination ranges cannot overlap.
func RemoveCopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate) ForwardWriter {
	for ; _ne(first, last); first = NextReader(first) {
		if !pred(first.Read()) {
			dFirst.Write(first.Read())
			dFirst = NextWriter(dFirst)
		}
	}
	return dFirst
}

// Replace replaces all elements equal to old with new in the range [first,
// last).
func Replace(first, last ForwardReadWriter, old, new Any) {
	ReplaceIf(first, last, _eq1(old), new)
}

// ReplaceIf replaces all elements satisfy pred with new in the range [first,
// last).
func ReplaceIf(first, last ForwardReadWriter, pred UnaryPredicate, v Any) {
	for ; _ne(first, last); first = NextReadWriter(first) {
		if pred(first.Read()) {
			first.Write(v)
		}
	}
}

// ReplaceCopy copies the elements from the range [first, last) to another range
// beginning at dFirst replacing all elements equal to old with new. The source
// and destination ranges cannot overlap.
func ReplaceCopy(first, last ForwardReader, dFirst ForwardWriter, old, new Any) ForwardWriter {
	return ReplaceCopyIf(first, last, dFirst, _eq1(old), new)
}

// ReplaceCopyIf copies the elements from the range [first, last) to another
// range beginning at dFirst replacing all elements satisfy pred with new. The
// source and destination ranges cannot overlap.
func ReplaceCopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate, v Any) ForwardWriter {
	for ; _ne(first, last); first, dFirst = NextReader(first), NextWriter(dFirst) {
		if pred(first.Read()) {
			dFirst.Write(v)
		} else {
			dFirst.Write(first.Read())
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
	for ; _ne(first1, last1); first1, first2 = NextReadWriter(first1), NextReadWriter(first2) {
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
func ReverseCopy(first, last BackwardReader, dFirst ForwardWriter) ForwardWriter {
	for _ne(first, last) {
		last = PrevReader(last)
		dFirst.Write(last.Read())
		dFirst = NextWriter(dFirst)
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
		write, read = NextReadWriter(write), NextReadWriter(read)
	}
	Rotate(write, nextRead, last)
	return write
}

// RotateCopy copies the elements from the range [first, last), to another range
// beginning at dFirst in such a way, that the element nFirst becomes the first
// element of the new range and nFirst - 1 becomes the last element.
func RotateCopy(first, nFirst, last ForwardReader, dFirst ForwardWriter) ForwardWriter {
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
func Sample(first, last ForwardReader, out ForwardWriter, n int, r *rand.Rand) ForwardWriter {
	_, rr := first.(RandomReader)
	rout, rw := out.(RandomWriter)
	if !rr && rw {
		return _reservoirSample(first, last, rout, n, r)
	}
	return _selectionSample(first, last, out, n, r)
}

func _selectionSample(first, last ForwardReader, out ForwardWriter, n int, r *rand.Rand) ForwardWriter {
	unsampled := Distance(first, last)
	if n > unsampled {
		n = unsampled
	}
	for ; n != 0; first = NextReader(first) {
		if r.Intn(unsampled) < n {
			out.Write(first.Read())
			out, n = NextWriter(out), n-1
		}
		unsampled--
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
// iterator for the new logical end of the range. Elements are compared using
// the given binary predicate pred.
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

// UniqueCopy copies the elements from the range [first, last), to another range
// beginning at d_first in such a way that there are no consecutive equal
// elements. Only the first element of each group of equal elements is copied.
func UniqueCopy(first, last ForwardReader, result ForwardWriter) ForwardWriter {
	return UniqueCopyIf(first, last, result, _eq)
}

// UniqueCopyIf copies the elements from the range [first, last), to another
// range beginning at d_first in such a way that there are no consecutive equal
// elements. Only the first element of each group of equal elements is copied.
// Elements are compared using the given binary predicate pred.
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

// IsPartitioned returns true if all elements in the range [first, last) that
// satisfy the predicate pred appear before all elements that don't. Also returns
// true if [first, last) is empty.
func IsPartitioned(first, last ForwardReader, pred UnaryPredicate) bool {
	return NoneOf(FindIfNot(first, last, pred), last, pred)
}

// Partition reorders the elements in the range [first, last) in such a way that
// all elements for which the predicate pred returns true precede the elements
// for which predicate pred returns false. Relative order of the elements is not
// preserved.
func Partition(first, last ForwardReadWriter, pred UnaryPredicate) ForwardReadWriter {
	first = FindIfNot(first, last, pred).(ForwardReadWriter)
	if _eq(first, last) {
		return first
	}
	for i := NextReadWriter(first); _ne(i, last); i = NextReadWriter(i) {
		if pred(i.Read()) {
			Swap(first, i)
			first = NextReadWriter(first)
		}
	}
	return first
}

// PartitionCopy copies the elements from the range [first, last) to two
// different ranges depending on the value returned by the predicate pred. The
// elements that satisfy the predicate pred are copied to the range beginning at
// outTrue. The rest of the elements are copied to the range beginning at
// outFalse.
func PartitionCopy(first, last ForwardReader, outTrue, outFalse ForwardWriter, pred UnaryPredicate) (ForwardWriter, ForwardWriter) {
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
		first = NextReadWriter(first)
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
		m := NextReadWriter(first)
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
		m1 = NextReadWriter(m1)
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

// PartitionPoint examines the partitioned (as if by std::partition) range
// [first, last) and locates the end of the first partition, that is, the first
// element that does not satisfy pred or last if all elements satisfy pred.
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
	return EqualBy(first1, last1, first2, last2, _eq)
}

func EqualBy(first1, last1, first2, last2 ForwardReader, pred BinaryPredicate) bool {
	for ; _ne(first1, last1); first1, first2 = NextReader(first1), NextReader(first2) {
		if !pred(first1.Read(), first2.Read()) {
			return false
		}
	}
	return true
}
