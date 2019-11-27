package iter

type (
	// Iter marks an iterator.
	Iter interface{}
	// Reader is a readable iterator.
	Reader interface {
		Read() Any
	}
	// Writer is a writable iterator.
	Writer interface {
		Write(Any)
	}
	// ReadWriter is an interface that groups Reader and Writer.
	ReadWriter interface {
		Reader
		Writer
	}
)

type (
	// ForwardIter is an iterator that moves forward.
	ForwardIter interface {
		Next() ForwardIter
	}
	// ForwardReader is an interface that groups ForwardIter and Reader.
	ForwardReader interface {
		ForwardIter
		Reader
	}
	// ForwardWriter is an interface that groups ForwardIter and Writer.
	ForwardWriter interface {
		ForwardIter
		Writer
	}
	// ForwardReadWriter is an interface that groups ForwardIter and
	// ReadWriter.
	ForwardReadWriter interface {
		ForwardIter
		ReadWriter
	}
)

// NextReader moves a ForwardReader to next.
func NextReader(r ForwardReader) ForwardReader {
	return r.Next().(ForwardReader)
}

// NextWriter moves a ForwardWriter to next.
func NextWriter(w ForwardWriter) ForwardWriter {
	return w.Next().(ForwardWriter)
}

// NextReadWriter moves a ReadWriter to next.
func NextReadWriter(rw ForwardReadWriter) ForwardReadWriter {
	return rw.Next().(ForwardReadWriter)
}

type (
	// BackwardIter is an iterator moves backward.
	BackwardIter interface {
		Prev() BackwardIter
	}
	// BackwardReader is an interface that groups BackwardIter and Reader.
	BackwardReader interface {
		BackwardIter
		Reader
	}
	// BackwardWriter is an interface that groups BackwardIter and Writer.
	BackwardWriter interface {
		BackwardIter
		Writer
	}
	// BackwardReadWriter is an interface that groups BackwardIter and
	// ReadWriter.
	BackwardReadWriter interface {
		BackwardIter
		ReadWriter
	}
)

// PrevReader moves a BackwardReader to prev.
func PrevReader(r BackwardReader) BackwardReader {
	return r.Prev().(BackwardReader)
}

// PrevWriter moves a BackwardWriter to prev.
func PrevWriter(w BackwardWriter) BackwardWriter {
	return w.Prev().(BackwardWriter)
}

// PrevReadWriter moves a BackwardReadWriter to prev.
func PrevReadWriter(rw BackwardReadWriter) BackwardReadWriter {
	return rw.Prev().(BackwardReadWriter)
}

type (
	// BidiIter is an iterator that moves both direction.
	BidiIter interface {
		ForwardIter
		BackwardIter
	}
	// BidiReader is an interface that groups BidiIter and Reader.
	BidiReader interface {
		BidiIter
		Reader
	}
	// BidiWriter is an interface that groups BidiIter and Writer.
	BidiWriter interface {
		BidiIter
		Writer
	}
	// BidiReadWriter is an interface that groups BidiIter and ReadWriter.
	BidiReadWriter interface {
		BidiIter
		ReadWriter
	}
)

// NextBidiIter moves a BidiIter to next.
func NextBidiIter(bi BidiIter) BidiIter {
	return bi.Next().(BidiIter)
}

// PrevBidiIter moves a BidiIter to prev.
func PrevBidiIter(bi BidiIter) BidiIter {
	return bi.Prev().(BidiIter)
}

// NextBidiReader moves a BidiReader to next.
func NextBidiReader(br BidiReader) BidiReader {
	return br.Next().(BidiReader)
}

// PrevBidiReader moves a BidiReader to prev.
func PrevBidiReader(br BidiReader) BidiReader {
	return br.Prev().(BidiReader)
}

// NextBidiWriter moves a BidiWriter to next.
func NextBidiWriter(br BidiWriter) BidiWriter {
	return br.Next().(BidiWriter)
}

// PrevBidiWriter moves a BidiWriter to prev.
func PrevBidiWriter(br BidiWriter) BidiWriter {
	return br.Prev().(BidiWriter)
}

// NextBidiReadWriter moves a BidiReadWriter to next.
func NextBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Next().(BidiReadWriter)
}

// PrevBidiReadWriter moves a BidiReadWriter to prev.
func PrevBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Prev().(BidiReadWriter)
}

type (
	// RandomIter is a random access iterator.
	RandomIter interface {
		BidiIter
		AdvanceN(n int) RandomIter
		Distance(RandomIter) int
	}
	// RandomReader is an interface that groups RandomIter and Reader.
	RandomReader interface {
		RandomIter
		Reader
	}
	// RandomWriter is an interface that groups RandomIter and Writer.
	RandomWriter interface {
		RandomIter
		Writer
	}
	// RandomReadWriter is an interface that groups RandomIter and
	// ReadWriter.
	RandomReadWriter interface {
		RandomIter
		ReadWriter
	}
)

// NextRandomIter moves a RandomIter to next.
func NextRandomIter(bi RandomIter) RandomIter {
	return bi.Next().(RandomIter)
}

// PrevRandomIter moves a RandomIter to prev.
func PrevRandomIter(bi RandomIter) RandomIter {
	return bi.Prev().(RandomIter)
}

// NextRandomReader moves a RandomReader to next.
func NextRandomReader(br RandomReader) RandomReader {
	return br.Next().(RandomReader)
}

// PrevRandomReader moves a RandomReader to prev.
func PrevRandomReader(br RandomReader) RandomReader {
	return br.Prev().(RandomReader)
}

// NextRandomWriter moves a RandomWriter to next.
func NextRandomWriter(br RandomWriter) RandomWriter {
	return br.Next().(RandomWriter)
}

// PrevRandomWriter moves a RandomWriter to prev.
func PrevRandomWriter(br RandomWriter) RandomWriter {
	return br.Prev().(RandomWriter)
}

// NextRandomReadWriter moves a RandomReadWriter to next.
func NextRandomReadWriter(br RandomReadWriter) RandomReadWriter {
	return br.Next().(RandomReadWriter)
}

// PrevRandomReadWriter moves a RandomReadWriter to prev.
func PrevRandomReadWriter(br RandomReadWriter) RandomReadWriter {
	return br.Prev().(RandomReadWriter)
}

// AdvanceNReader moves a RandomReader by step N.
func AdvanceNReader(rr RandomReader, n int) RandomReader {
	return rr.AdvanceN(n).(RandomReader)
}

// AdvanceNWriter moves a RandomWriter by step N.
func AdvanceNWriter(rw RandomWriter, n int) RandomWriter {
	return rw.AdvanceN(n).(RandomWriter)
}

// AdvanceNReadWriter moves a RandomReadWriter by step N.
func AdvanceNReadWriter(rw RandomReadWriter, n int) RandomReadWriter {
	return rw.AdvanceN(n).(RandomReadWriter)
}

// Distance returns the distance of two iterators.
func Distance(first, last Iter) int {
	if f, ok := first.(RandomIter); ok {
		if l, ok := last.(RandomIter); ok {
			return f.Distance(l)
		}
	}
	if f, ok := first.(ForwardIter); ok {
		if l, ok := last.(ForwardIter); ok {
			var d int
			for ; _ne(f, l); f = f.Next() {
				d++
			}
			return d
		}
	}
	if f, ok := first.(BackwardIter); ok {
		if l, ok := last.(BackwardIter); ok {
			var d int
			for ; _ne(f, l); l = l.Prev() {
				d++
			}
			return d
		}
	}
	panic("cannot get distance")
}

// AdvanceN moves an iterator by step N.
func AdvanceN(it Iter, n int) Iter {
	if it2, ok := it.(RandomIter); ok {
		return it2.AdvanceN(n)
	}
	if it2, ok := it.(ForwardIter); ok && n >= 0 {
		for ; n > 0; n-- {
			it2 = it2.Next()
		}
		return it2
	}
	if it2, ok := it.(BackwardIter); ok && n <= 0 {
		for ; n < 0; n++ {
			it2 = it2.Prev()
		}
		return it2
	}
	panic("cannot advance")
}
