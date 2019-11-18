package iter

type Iter interface{}

type Readable interface {
	Read() interface{}
}

type Writable interface {
	Write(interface{})
}

type ReadWriter interface {
	Readable
	Writable
}

type ForwardIter interface {
	Next() ForwardIter
}

type BackwardIter interface {
	Prev() BackwardIter
}

type ForwardReader interface {
	Readable
	ForwardIter
}

func NextReader(r ForwardReader) ForwardReader {
	return r.Next().(ForwardReader)
}

type BackwardReader interface {
	Readable
	BackwardIter
}

func PrevReader(r BackwardReader) BackwardReader {
	return r.Prev().(BackwardReader)
}

type ForwardWriter interface {
	Writable
	ForwardIter
}

func NextWriter(w ForwardWriter) ForwardWriter {
	return w.Next().(ForwardWriter)
}

type BackwardWriter interface {
	Writable
	BackwardIter
}

func PrevWriter(w BackwardWriter) BackwardWriter {
	return w.Prev().(BackwardWriter)
}

type ForwardReadWriter interface {
	ReadWriter
	ForwardIter
}

func NextReadWriter(rw ForwardReadWriter) ForwardReadWriter {
	return rw.Next().(ForwardReadWriter)
}

type BidiIter interface {
	ForwardIter
	BackwardIter
}

func NextBidiIter(bi BidiIter) BidiIter {
	return bi.Next().(BidiIter)
}

func PrevBidiIter(bi BidiIter) BidiIter {
	return bi.Prev().(BidiIter)
}

type BidiReader interface {
	Readable
	BidiIter
}

func NextBidiReader(br BidiReader) BidiReader {
	return br.Next().(BidiReader)
}

func PrevBidiReader(br BidiReader) BidiReader {
	return br.Prev().(BidiReader)
}

type BidiWriter interface {
	Writable
	BidiIter
}

func NextBidiWriter(br BidiWriter) BidiWriter {
	return br.Next().(BidiWriter)
}

func PrevBidiWriter(br BidiWriter) BidiWriter {
	return br.Prev().(BidiWriter)
}

type BidiReadWriter interface {
	ReadWriter
	BidiIter
}

func NextBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Next().(BidiReadWriter)
}

func PrevBidiReadWriter(br BidiReadWriter) BidiReadWriter {
	return br.Prev().(BidiReadWriter)
}

type RandomIter interface {
	ForwardIter
	BackwardIter
	AdvanceN(n int) RandomIter
	Distance(RandomIter) int
}

type RandomReader interface {
	Readable
	RandomIter
}

func AdvanceNReader(rr RandomReader, n int) RandomReader {
	return rr.AdvanceN(n).(RandomReader)
}

type RandomWriter interface {
	Writable
	RandomIter
}

func AdvanceNWriter(rw RandomWriter, n int) RandomWriter {
	return rw.AdvanceN(n).(RandomWriter)
}

type RandomReadWriter interface {
	ReadWriter
	RandomIter
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

func Next(it Iter) Iter {
	return AdvanceN(it, 1)
}

func Prev(it Iter) Iter {
	return AdvanceN(it, -1)
}
