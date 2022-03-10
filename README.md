# iter

**Note: I'm currently working on making it work with Go generics. For previous version, please check `v1` branch.**

Go implementation of C++ STL iterators and algorithms.

Less hand-written loops, more expressive code.

README translations: [简体中文](README_ZH.md)

[![GoDoc](https://godoc.org/github.com/disksing/iter?status.svg)](https://godoc.org/github.com/disksing/iter)
[![Build Status](https://travis-ci.com/disksing/iter.svg?branch=master)](https://travis-ci.com/disksing/iter)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/disksing/iter)](https://goreportcard.com/report/github.com/disksing/iter)

## Motivation

Although Go doesn't have generics, we deserve to have reuseable general algorithms. `iter` helps improving Go code in several ways:

- Some simple loops are unlikely to be wrong or inefficient, but calling algorithm instead will **make the code more concise and easier to comprehend**. Such as [AllOf](https://godoc.org/github.com/disksing/iter#AllOf), [FindIf](https://godoc.org/github.com/disksing/iter#FindIf), [Accumulate](https://godoc.org/github.com/disksing/iter#Accumulate).

- Some algorithms are not complicated, but it is not easy to write them correctly. **Reusing code makes them easier to reason for correctness**. Such as [Shuffle](https://godoc.org/github.com/disksing/iter#Shuffle), [Sample](https://godoc.org/github.com/disksing/iter#Sample), [Partition](https://godoc.org/github.com/disksing/iter#Partition).

- STL also includes some complicated algorithms that may take hours to make it correct. **Implementing it manually is impractical**. Such as [NthElement](https://godoc.org/github.com/disksing/iter#NthElement), [StablePartition](https://godoc.org/github.com/disksing/iter#StablePartition), [NextPermutation](https://godoc.org/github.com/disksing/iter#NextPermutation).

- The implementation in the library contains some **imperceptible performance optimizations**. For instance, [MinmaxElement](https://godoc.org/github.com/disksing/iter#MinmaxElement) is done by taking two elements at a time. In this way, the overall number of comparisons is significantly reduced.

There are alternative libraries have similar goals, such as [gostl](https://github.com/liyue201/gostl), [gods](https://github.com/emirpasic/gods) and [go-stp](https://github.com/itrabbit/go-stp). What makes `iter` unique is:

- **Non-intrusive**. Instead of introducing new containers, `iter` tends to reuse existed containers in Go (slice, string, list.List, etc.) and use iterators to adapt them to algorithms.

- **Full algorithms (>100)**. It includes almost all algorithms come before C++17. Check the [Full List](https://github.com/disksing/iter/wiki/Algorithms).

## Examples

> The examples are run with some function alias to make it simple. See [example_test.go](https://github.com/disksing/iter/blob/master/examples_test.go) for the detail.

<table>
<thead><tr><th colspan="2">Print a list.List</th></tr></thead>
<tbody><td>

```go
l := list.New()
for i := 1; i <= 5; i++ {
  l.PushBack(i)
}
for e := l.Front(); e != nil; e = e.Next() {
  fmt.Print(e.Value)
  if e.Next() != nil {
    fmt.Print("->")
  }
}
// Output:
// 1->2->3->4->5
```

</td><td>

```go
l := list.New()
GenerateN(ListBackInserter(l), 5, IotaGenerator(1))
Copy(lBegin(l), lEnd(l), IOWriter(os.Stdout, "->"))
// Output:
// 1->2->3->4->5
```

</td></tr></tbody>

<thead><tr><th colspan="2">Reverse a string</th></tr></thead>
<tbody><tr><td>

```go
s := "!dlrow olleH"
var sb strings.Builder
for i := len(s) - 1; i >= 0; i-- {
  sb.WriteByte(s[i])
}
fmt.Println(sb.String())

b := []byte(s)
for i := len(s)/2 - 1; i >= 0; i-- {
  j := len(s) - 1 - i
  b[i], b[j] = b[j], b[i]
}
fmt.Println(string(b))
// Output:
// Hello world!
// Hello world!
```

</td><td>

```go
s := "!dlrow olleH"
fmt.Println(MakeString(StringRBegin(s), StringREnd(s)))

b := []byte(s)
Reverse(begin(b), end(b))
fmt.Println(string(b))
// Output:
// Hello world!
// Hello world!
```

</td></tr></tbody>

<thead><tr><th colspan="2">In-place deduplicate (from <a href="https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable">SliceTricks</a>, with minor change)</th></tr></thead>
<tbody><tr><td>

```go
in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
sort.Ints(in)
j := 0
for i := 1; i < len(in); i++ {
  if in[j] == in[i] {
    continue
  }
  j++
  in[j] = in[i]
}
in = in[:j+1]
fmt.Println(in)
// Output:
// [1 2 3 4]
```

</td><td>

```go
in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
Sort(begin(in), end(in))
Erase(&in, Unique(begin(in), end(in)))
fmt.Println(in)
// Output:
// [1 2 3 4]
```

</td></tr></tbody>

<thead><tr><th colspan="2">Sum all integers received from a channel</th></tr></thead>
<tbody><tr><td>

```go
ch := make(chan int)
go func() {
  for _, x := range rand.Perm(100) {
    ch <- x + 1
  }
  close(ch)
}()
var sum int
for x := range ch {
  sum += x
}
fmt.Println(sum)
// Output:
// 5050
```

</td><td>

```go
ch := make(chan int)
go func() {
  CopyN(IotaReader(1), 100, ChanWriter(ch))
  close(ch)
}()
fmt.Println(Accumulate(ChanReader(ch), ChanEOF, 0))
// Output:
// 5050
```

</td></tr></tbody>

<thead><tr><th colspan="2">Remove consecutive spaces in a string</th></tr></thead>
<tbody><tr><td>

```go
str := "  a  quick   brown  fox  "
var sb strings.Builder
var prevIsSpace bool
for i := 0; i < len(str); i++ {
  if str[i] != ' ' || !prevIsSpace {
    sb.WriteByte(str[i])
  }
  prevIsSpace = str[i] == ' '
}
fmt.Println(sb.String())
// Output:
// a quick brown fox
```

</td><td>

```go
str := "  a  quick   brown  fox  "
var sb StringBuilderInserter
UniqueCopyIf(sBegin(str), sEnd(str), &sb,
  func(x, y Any) bool { return x.(byte) == ' ' && y.(byte) == ' ' })
fmt.Println(sb.String())
// Output:
// a quick brown fox
```

</td></tr></tbody>

<thead><tr><th colspan="2">Collect N maximum elements from a channel</th></tr></thead>
<tbody><tr><td>

```go
// Need to manually mantain a min-heap.
```

</td><td>

```go
top := make([]int, 5)
PartialSortCopyBy(ChanReader(ch), ChanEOF, begin(top), end(top),
  func(x, y Any) bool { return x.(int) > y.(int) })
Copy(begin(top), end(top), IOWriter(os.Stdout, ", "))
```

</td></tr></tbody>

<thead><tr><th colspan="2">Print all permutations of ["a", "b", "c"]</th></tr></thead>
<tbody><tr><td>

```go
// Usually requires some sort of recursion
```

</td><td>

```go
s := []string{"a", "b", "c"}
for ok := true; ok; ok = NextPermutation(begin(s), end(s)) {
  fmt.Println(s)
}
// Output:
// [a b c]
// [a c b]
// [b a c]
// [b c a]
// [c a b]
// [c b a]
```

</td></tr></tbody>
</table>

## Thanks

- [cppreference.com](https://en.cppreference.com/)
- [LLVM libc++](https://libcxx.llvm.org/)

## License

BSD 3-Clause
