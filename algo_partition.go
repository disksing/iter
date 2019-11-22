package iter

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
