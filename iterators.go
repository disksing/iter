package iter

import (
	"fmt"
	"io"
)

type chanReader[T any, C Chan[T]] struct {
	ch    C
	cur   T
	read1 bool
	eof   bool
}

func (cr *chanReader[T]) recv() {
	v, ok := <-cr.ch
	cr.cur, cr.read1, cr.eof = v, true, !ok
}

func (cr *chanReader[T]) Read() T {
	if !cr.read1 {
		cr.recv()
	}
	return cr.cur
}

func (cr *chanReader[T]) Next() *chanReader[T] {
	if !cr.read1 {
		cr.recv()
	}
	if !cr.eof {
		cr.recv()
	}
	return cr
}

func (cr *chanReader[T]) Eq(x *chanReader[T]) bool {
	if !cr.read1 {
		cr.recv()
	}
	return cr.eof && x == nil
}

type chanWriter[T any, C Chan[T]] struct {
	ch C
}

func (cr *chanWriter[T]) Write(x T) {
	cr.ch <- x
}

// // ChanReader returns an InputIter that reads from a channel.
// func ChanReader[T any, C Chan[T]](c C) *chanReader[T, C] {
// 	return &chanReader[T, C]{
// 		ch: c,
// 	}
// }

// // ChanWriter returns an OutIter that writes to a channel.
// func ChanWriter[T any, C Chan[T]](c C) *chanWriter[T, C] {
// 	return &chanWriter[T, C]{
// 		ch: c,
// 	}
// }

type ioWriter struct {
	w         io.Writer
	written   bool
	delimiter []byte
}

func (w *ioWriter) Write(x any) {
	if w.written && len(w.delimiter) > 0 {
		_, err := w.w.Write(w.delimiter)
		if err != nil {
			panic(err)
		}
	} else {
		w.written = true
	}

	_, err := fmt.Fprint(w.w, x)
	if err != nil {
		panic(err)
	}
}

// IOWriter returns an OutputIter that writes values to an io.Writer.
// It panics if meet any error.
func IOWriter(w io.Writer, delimiter string) *ioWriter {
	return &ioWriter{w: w, delimiter: []byte(delimiter)}
}
