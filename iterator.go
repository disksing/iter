package iter

// Iter represents an iterator, just an alias of any.
type Iter[T any] interface{}

type (
	// Reader is a readable iterator.
	Reader[T any] interface {
		Read() T
	}
	// Writer is a writable iterator.
	Writer[T any] interface {
		Write(T)
	}
	// ReadWriter is an interface that groups Reader and Writer.
	ReadWriter[T any] interface {
		Reader[T]
		Writer[T]
	}
)

// Comparable represents an iterator that can be compared.
type Comparable[It Iter[any]] interface {
	Eq(It) bool
}

// ForwardMovable represents iterators that can move forward.
type ForwardMovable[It Iter[any]] interface {
	Next() It
}

// BackwardMovable represents iterators that can move backward.
type BackwardMovable[It Iter[any]] interface {
	Prev() It
}

// InputIter is a readable and forward movable iterator.
type InputIter[T any, It Iter[T]] interface {
	Reader[T]
	ForwardMovable[It]
	Comparable[It]
}

func __eq[It Comparable[It]](x, y It) bool {
	return x.Eq(y)
}

func __ne[It Comparable[It]](x, y It) bool {
	return !__eq(x, y)
}

// OutputIter is a writable and ForwardMovable iterator.
//
// It may not implement the incremental interface, in which case the increment
// logic is done in Write().
type OutputIter[T any] interface {
	Writer[T]
}

func _writeNext[T any, It OutputIter[T]](out It, v T) It {
	out.Write(v)
	if inc, ok := interface{}(out).(ForwardMovable[It]); ok {
		out = inc.Next()
	}
	return out
}

type (
	// ForwardIter is an iterator that moves forward.
	ForwardIter[T any, It Iter[T]] interface {
		ForwardMovable[It]
		Comparable[It]
		AllowMultiplePass() // a marker indicates it can be multiple passed.
	}
	// ForwardReader is an interface that groups ForwardIter and Reader.
	ForwardReader[T any, It Iter[T]] interface {
		ForwardIter[T, It]
		Reader[T]
	}
	// ForwardWriter is an interface that groups ForwardIter and Writer.
	ForwardWriter[T any, It Iter[T]] interface {
		ForwardIter[T, It]
		Writer[T]
	}
	// ForwardReadWriter is an interface that groups ForwardIter and
	// ReadWriter.
	ForwardReadWriter[T any, It Iter[T]] interface {
		ForwardIter[T, It]
		ReadWriter[T]
	}
)

type (
	// BidiIter is an iterator that moves both forward or backward.
	BidiIter[T any, It Iter[T]] interface {
		ForwardIter[T, It]
		BackwardMovable[It]
	}
	// BidiReader is an interface that groups BidiIter and Reader.
	BidiReader[T any, It Iter[T]] interface {
		BidiIter[T, It]
		Reader[T]
	}
	// BidiWriter is an interface that groups BidiIter and Writer.
	BidiWriter[T any, It Iter[T]] interface {
		BidiIter[T, It]
		Writer[T]
	}
	// BidiReadWriter is an interface that groups BidiIter and ReadWriter.
	BidiReadWriter[T any, It Iter[T]] interface {
		BidiIter[T, It]
		ReadWriter[T]
	}
)

type (
	// RandomIter is a random access iterator.
	RandomIter[T any, It Iter[T]] interface {
		BidiIter[T, It]
		AdvanceN(n int) It
		Distance(It) int
		Less(It) bool
	}
	// RandomReader is an interface that groups RandomIter and Reader.
	RandomReader[T any, It Iter[T]] interface {
		RandomIter[T, It]
		Reader[T]
	}
	// RandomWriter is an interface that groups RandomIter and Writer.
	RandomWriter[T any, It Iter[T]] interface {
		RandomIter[T, It]
		Writer[T]
	}
	// RandomReadWriter is an interface that groups RandomIter and
	// ReadWriter.
	RandomReadWriter[T any, It Iter[T]] interface {
		RandomIter[T, It]
		ReadWriter[T]
	}
)

// Distance returns the distance of two iterators.
func Distance[T any, It Iter[T]](first, last It) int {
	ifirst, ilast := interface{}(first), interface{}(last)
	if f, ok := ifirst.(RandomIter[T, It]); ok {
		if l, ok := ilast.(It); ok {
			return f.Distance(l)
		}
	}
	if f, ok := ifirst.(ForwardIter[T, It]); ok {
		if l, ok := ilast.(It); ok {
			var d int
			for ; !f.Eq(l); f = (interface{})(f.Next()).(ForwardIter[T, It]) {
				d++
			}
			return d
		}
	}
	if f, ok := ifirst.(InputIter[T, It]); ok {
		var d int
		for ; !f.Eq(last); f = (interface{})(f.Next()).(InputIter[T, It]) {
			d++
		}
		return d
	}
	panic("cannot get distance")
}

// AdvanceN moves an iterator by step N.
func AdvanceN[T any, It Iter[T]](it It, n int) It {
	if it2, ok := interface{}(it).(RandomIter[T, It]); ok {
		return it2.AdvanceN(n)
	}
	if it2, ok := interface{}(it).(ForwardIter[T, It]); ok && n >= 0 {
		for ; n > 0; n-- {
			it2 = (interface{})(it2.Next()).(ForwardIter[T, It])
		}
		return it2.(It)
	}
	if it2, ok := interface{}(it).(InputIter[T, It]); ok && n >= 0 {
		for ; n > 0; n-- {
			it2 = (interface{})(it2.Next()).(InputIter[T, It])
		}
		return it2.(It)
	}
	if it2, ok := interface{}(it).(BidiIter[T, It]); ok && n <= 0 {
		for ; n < 0; n++ {
			it2 = (interface{})(it2.Prev()).(BidiIter[T, It])
		}
		return it2.(It)
	}
	panic("cannot advance")
}
