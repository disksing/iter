package iter

import (
	"fmt"
	"io"
	"math/rand"
)

// IotaIterator is an input iterator that yields x, x+1, x+2, and so on.
type IotaIterator[T Numeric] struct {
	x T
}

func (r IotaIterator[T]) Read() T {
	return r.x
}

func (r IotaIterator[T]) Next() IotaIterator[T] {
	return IotaIterator[T]{x: r.x + 1}
}

func (r IotaIterator[T]) Eq(IotaIterator[T]) bool {
	return false
}

// IotaReader creates an InputIter that returns [x, x+1, x+2...).
func IotaReader[T Numeric](x T) IotaIterator[T] {
	return IotaIterator[T]{x: x}
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

// RepeatIterator is an input iterator that yields the same value indefinitely.
type RepeatIterator[T any] struct {
	x T
}

func (r RepeatIterator[T]) Read() T { return r.x }

func (r RepeatIterator[T]) Next() RepeatIterator[T] { return r }

func (r RepeatIterator[T]) Eq(RepeatIterator[T]) bool { return false }

// RepeatReader creates an InputIter that returns [x, x, x...).
func RepeatReader[T any](x T) RepeatIterator[T] {
	return RepeatIterator[T]{x: x}
}

// RepeatGenerator creates a Generator that returns [x, x, x...).
func RepeatGenerator[T any](x T) func() T {
	return func() T { return x }
}

// RandomGenerator creates a generator that returns random item of a slice.
func RandomGenerator[T any](s []T, r *rand.Rand) func() T {
	return func() T { return s[r.Intn(len(s))] }
}

// ChannelReader is a single-pass input iterator over a channel.
// A nil *ChannelReader is its end sentinel.
type ChannelReader[T any] struct {
	ch    chan T
	cur   T
	read1 bool
	eof   bool
}

func (cr *ChannelReader[T]) recv() {
	v, ok := <-cr.ch
	cr.cur, cr.read1, cr.eof = v, true, !ok
}

func (cr *ChannelReader[T]) Read() T {
	if !cr.read1 {
		cr.recv()
	}
	return cr.cur
}

func (cr *ChannelReader[T]) Next() *ChannelReader[T] {
	if !cr.read1 {
		cr.recv()
	}
	if !cr.eof {
		cr.recv()
	}
	return cr
}

func (cr *ChannelReader[T]) Eq(x *ChannelReader[T]) bool {
	if !cr.read1 {
		cr.recv()
	}
	return cr.eof && x == nil
}

// ChannelWriter is an output iterator that sends values to a channel.
type ChannelWriter[T any] struct {
	ch chan T
}

func (cr *ChannelWriter[T]) Write(x T) {
	cr.ch <- x
}

// ChanReader returns an InputIter that reads from a channel.
func ChanReader[T any](c chan T) *ChannelReader[T] {
	return &ChannelReader[T]{
		ch: c,
	}
}

// ChanWriter returns an OutputIter that writes to a channel.
func ChanWriter[T any](c chan T) *ChannelWriter[T] {
	return &ChannelWriter[T]{
		ch: c,
	}
}

// OutputWriter is an output iterator that formats values to an io.Writer.
type OutputWriter[T any] struct {
	w         io.Writer
	written   bool
	delimiter []byte
}

func (w *OutputWriter[T]) Write(x T) {
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
// It panics if the writer returns an error.
func IOWriter[T any](w io.Writer, delimiter string) *OutputWriter[T] {
	return &OutputWriter[T]{w: w, delimiter: []byte(delimiter)}
}
