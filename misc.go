package iter

import (
	"fmt"
	"io"
	"math/rand"
)

type iotaReader[T Numeric] struct {
	x T
}

func (r iotaReader[T]) Read() T {
	return r.x
}

func (r iotaReader[T]) Next() iotaReader[T] {
	return iotaReader[T]{x: r.x + 1}
}

func (r iotaReader[T]) Eq(iotaReader[T]) bool {
	return false
}

// IotaReader creates an InputIter that returns [x, x+1, x+2...).
func IotaReader[T Numeric, It iotaReader[T]](x T) It {
	return It{x: x}
}

// IotaGenerator creates a Generator that returns [x, x+1, x+2...).
func IotaGenerator[T Numeric](x T) func() T {
	r := IotaReader(x)
	return func() T {
		v := r.Read()
		r = r.Next()
		return v
	}
}

type repeatReader[T any] struct {
	x T
}

func (r repeatReader[T]) Read() T { return r.x }

func (r repeatReader[T]) Next() repeatReader[T] { return r }

func (r repeatReader[T]) Eq(repeatReader[T]) bool { return false }

// RepeatReader creates an InputIter that returns [x, x, x...).
func RepeatReader[T any](x T) repeatReader[T] {
	return repeatReader[T]{x: x}
}

// RepeatGenerator creates an Generator that returns [x, x, x...).
func RepeatGenerator[T any](x T) func() T {
	return func() T { return x }
}

// RandomGenerator creates a generator that returns random item of a slice.
func RandomGenerator[T any](s []T, r *rand.Rand) func() T {
	return func() T { return s[r.Intn(len(s))] }
}

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
