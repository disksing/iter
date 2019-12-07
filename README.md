# iter

GO implementation of C++ STL iterators and algorithms.

Less hand-written loops, more expressive code.

[![GoDoc](https://godoc.org/github.com/disksing/iter?status.svg)](https://godoc.org/github.com/disksing/iter)
[![Build Status](https://travis-ci.com/disksing/iter.svg?branch=master)](https://travis-ci.com/disksing/iter)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/disksing/iter)](https://goreportcard.com/report/github.com/disksing/iter)

### Motivation

Although Go doesn't have generics, we deserve to have reuseable general algorithms. `iter` can help improve Go code in several ways:

- Some simple loops are unlikely to be wrong or inefficient, but calling algorithm instead will **make the code more concise and easier to comprehend**. Such as [AllOf](https://godoc.org/github.com/disksing/iter#AllOf), [FindIf](https://godoc.org/github.com/disksing/iter#FindIf), [Accumulate](https://godoc.org/github.com/disksing/iter#Accumulate).

- Some algorithms are not complicated, but it is not easy to write them correctly. **Reusing code makes them easier to reason for correctness**. Such as [Shuffle](https://godoc.org/github.com/disksing/iter#Shuffle), [Sample](https://godoc.org/github.com/disksing/iter#Sample), [Partition](https://godoc.org/github.com/disksing/iter#Partition).

- STL also includes some complicated algorithms that may take hours to make it correct. **Implementing it manually is impractical**. Such as [NthElement](https://godoc.org/github.com/disksing/iter#NthElement), [StablePartition](https://godoc.org/github.com/disksing/iter#StablePartition), [NextPermutation](https://godoc.org/github.com/disksing/iter#NextPermutation).

- The implementation in the library contains some **imperceptible performance optimizations**. For instance, [MinmaxElement](https://godoc.org/github.com/disksing/iter#MinmaxElement) is done by taking two elements at a time. In this way, the overall number of comparisons is significantly reduced.

There are alternative libraries have similar goals, such as [gostl](https://github.com/liyue201/gostl), [gods](https://github.com/emirpasic/gods) and [go-stp](https://github.com/itrabbit/go-stp). What makes `iter` unique is:

- **None-intrusive**. Instead of introducing new containers, `iter` tends to reuse existed containers in Go (slice, string, list.List, etc.) and use iterators to adapt them to algorithms.

- **Full algorithms (>100)**. It includes almost all algorithms come before C ++ 17. Check the [Full List](https://godoc.org/github.com/disksing/iter).

### Examples

### Thanks

- [cppreference.com](https://en.cppreference.com/)
- [LLVM libc++](https://libcxx.llvm.org/)

### License

BSD 3-Clause
