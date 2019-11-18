package iter

import "math/rand"

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

func Fill(first, last ForwardWriter, v interface{}) {
	for ; _ne(first, last); first = NextWriter(first) {
		first.Write(v)
	}
}

func FillN(first ForwardWriter, count int, v interface{}) {
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

func Remove(first, last ForwardReadWriter, v interface{}) ForwardReadWriter {
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

func RemoveCopy(first, last ForwardReader, dFirst ForwardWriter, v interface{}) ForwardWriter {
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

func Replace(first, last ForwardReadWriter, old, new interface{}) {
	ReplaceIf(first, last, _eq1(old), new)
}

func ReplaceIf(first, last ForwardReadWriter, pred UnaryPredicate, v interface{}) {
	for ; _ne(first, last); first = NextReadWriter(first) {
		if pred(first.Read()) {
			first.Write(v)
		}
	}
}

func ReplaceCopy(first, last ForwardReader, dFirst ForwardWriter, old, new interface{}) ForwardWriter {
	return ReplaceCopyIf(first, last, dFirst, _eq1(old), new)
}

func ReplaceCopyIf(first, last ForwardReader, dFirst ForwardWriter, pred UnaryPredicate, v interface{}) ForwardWriter {
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
