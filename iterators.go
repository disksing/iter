package iter

type (
	Iter   interface{}
	Reader interface {
		Read() Any
	}
	Writer interface {
		Write(Any)
	}
	ReadWriter interface {
		Reader
		Writer
	}
)

type (
	ForwardIter interface {
		Next() ForwardIter
	}
	ForwardReader interface {
		ForwardIter
		Reader
	}
	ForwardWriter interface {
		ForwardIter
		Writer
	}
	ForwardReadWriter interface {
		ForwardIter
		ReadWriter
	}
)

func NextReader(r ForwardReader) ForwardReader {
	return r.Next().(ForwardReader)
}

func NextWriter(w ForwardWriter) ForwardWriter {
	return w.Next().(ForwardWriter)
}

func NextReadWriter(rw ForwardReadWriter) ForwardReadWriter {
	return rw.Next().(ForwardReadWriter)
}

type (
	BackwardIter interface {
		Prev() BackwardIter
	}
	BackwardReader interface {
		BackwardIter
		Reader
	}
	BackwardWriter interface {
		BackwardIter
		Writer
	}
	BackwardReadWriter interface {
		BackwardIter
		ReadWriter
	}
)

func PrevReader(r BackwardReader) BackwardReader {
	return r.Prev().(BackwardReader)
}

func PrevWriter(w BackwardWriter) BackwardWriter {
	return w.Prev().(BackwardWriter)
}

func PrevReadWriter(rw BackwardReadWriter) BackwardReadWriter {
	return rw.Prev().(BackwardReadWriter)
}

type (
	BidiIter interface {
		ForwardIter
		BackwardIter
	}
	BidiReader interface {
		BidiIter
		Reader
	}
	BidiWriter interface {
		BidiIter
		Writer
	}
	BidiReadWriter interface {
		BidiIter
		ReadWriter
	}
)

func NextBidiIter(bi BidiIter) BidiIter {
	return bi.Next().(BidiIter)
}

func PrevBidiIter(bi BidiIter) BidiIter {
	return bi.Prev().(BidiIter)
}

func NextBidiReader(br BidiReader) BidiReader {
	return br.Next().(BidiReader)
}

func PrevBidiReader(br BidiReader) BidiReader {
	return br.Prev().(BidiReader)
}

func NextBidiWriter(br BidiWriter) BidiWriter {
	return br.Next().(BidiWriter)
}

func PrevBidiWriter(br BidiWriter) BidiWriter {
	return br.Prev().(BidiWriter)
}

func NextBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Next().(BidiReadWriter)
}

func PrevBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Prev().(BidiReadWriter)
}

type (
	RandomIter interface {
		BidiIter
		AdvanceN(n int) RandomIter
		Distance(RandomIter) int
	}
	RandomReader interface {
		RandomIter
		Reader
	}
	RandomWriter interface {
		RandomIter
		Writer
	}
	RandomReadWriter interface {
		RandomIter
		ReadWriter
	}
)

func NextRandomIter(bi RandomIter) RandomIter {
	return bi.Next().(RandomIter)
}

func PrevRandomIter(bi RandomIter) RandomIter {
	return bi.Prev().(RandomIter)
}

func NextRandomReader(br RandomReader) RandomReader {
	return br.Next().(RandomReader)
}

func PrevRandomReader(br RandomReader) RandomReader {
	return br.Prev().(RandomReader)
}

func NextRandomWriter(br RandomWriter) RandomWriter {
	return br.Next().(RandomWriter)
}

func PrevRandomWriter(br RandomWriter) RandomWriter {
	return br.Prev().(RandomWriter)
}

func NextRandomReadWriter(br RandomReadWriter) RandomReadWriter {
	return br.Next().(RandomReadWriter)
}

func PrevRandomReadWriter(br RandomReadWriter) RandomReadWriter {
	return br.Prev().(RandomReadWriter)
}

func AdvanceNReader(rr RandomReader, n int) RandomReader {
	return rr.AdvanceN(n).(RandomReader)
}

func AdvanceNWriter(rw RandomWriter, n int) RandomWriter {
	return rw.AdvanceN(n).(RandomWriter)
}

func AdvanceNReadWriter(rw RandomReadWriter, n int) RandomReadWriter {
	return rw.AdvanceN(n).(RandomReadWriter)
}

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
