// Package lists adapts container/list values to the generic iterator model.
//
// A container/list stores values as any. Reading an element through
// Iterator[T] therefore panics when the stored value is not a T.
package lists
