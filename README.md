# iter

**Note: The generic v2 API is still under development and has not been released. It requires Go 1.26 or later. For the previous version, please check the [`v1` branch](https://github.com/disksing/iter/tree/v1).**

Go implementation of C++ STL iterators and algorithms.

Less hand-written loops, more expressive code.

README translations: [简体中文](README_ZH.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/disksing/iter/v2.svg)](https://pkg.go.dev/github.com/disksing/iter/v2)
[![CI](https://github.com/disksing/iter/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/disksing/iter/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/disksing/iter/v2)](https://goreportcard.com/report/github.com/disksing/iter/v2)

## Motivation

Go now has generics, and we deserve to have reusable general algorithms. `iter` helps improve Go code in several ways:

- Some simple loops are unlikely to be wrong or inefficient, but calling an algorithm instead will **make the code more concise and easier to comprehend**. Such as [AllOf](https://pkg.go.dev/github.com/disksing/iter/v2/algo#AllOf), [FindIf](https://pkg.go.dev/github.com/disksing/iter/v2/algo#FindIf), [Accumulate](https://pkg.go.dev/github.com/disksing/iter/v2/algo#Accumulate).

- Some algorithms are not complicated, but it is not easy to write them correctly. **Reusing code makes them easier to reason for correctness**. Such as [Shuffle](https://pkg.go.dev/github.com/disksing/iter/v2/algo#Shuffle), [Sample](https://pkg.go.dev/github.com/disksing/iter/v2/algo#Sample), [Partition](https://pkg.go.dev/github.com/disksing/iter/v2/algo#Partition).

- STL also includes some complicated algorithms that may take hours to make it correct. **Implementing it manually is impractical**. Such as [NthElement](https://pkg.go.dev/github.com/disksing/iter/v2/algo#NthElement), [StablePartition](https://pkg.go.dev/github.com/disksing/iter/v2/algo#StablePartition), [NextPermutation](https://pkg.go.dev/github.com/disksing/iter/v2/algo#NextPermutation).

- The implementation in the library contains some **imperceptible performance optimizations**. For instance, [MinmaxElement](https://pkg.go.dev/github.com/disksing/iter/v2/algo#MinmaxElement) is done by taking two elements at a time. In this way, the overall number of comparisons is significantly reduced.

There are alternative libraries have similar goals, such as [gostl](https://github.com/liyue201/gostl), [gods](https://github.com/emirpasic/gods) and [go-stp](https://github.com/itrabbit/go-stp). What makes `iter` unique is:

- **Non-intrusive**. Instead of introducing new containers, `iter` tends to reuse existed containers in Go (slice, string, list.List, etc.) and use iterators to adapt them to algorithms.

- **Full algorithms (>100)**. It includes almost all algorithms come before C++17. Check the [Full List](https://github.com/disksing/iter/wiki/Algorithms).

## Examples

> The examples keep package qualifiers so the generic v2 API is explicit. See [examples_test.go](examples_test.go) for complete, executable examples.

<table>
<thead><tr><th colspan="2">Print a list.List</th></tr></thead>
<tbody><tr><td>

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
algo.GenerateN(lists.ListBackInserter[int](l), 5, iter.IotaGenerator(1))
algo.Copy(lists.Begin[int](l), lists.End[int](l), iter.IOWriter[int](os.Stdout, "->"))
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
fmt.Println(strs.MakeString(strs.RBegin(s), strs.REnd(s)))

b := []byte(s)
iterslices.Reverse(b)
fmt.Println(string(b))
// Output:
// Hello world!
// Hello world!
```

</td></tr></tbody>

<thead><tr><th colspan="2">In-place deduplicate (from <a href="https://go.dev/wiki/SliceTricks#in-place-deduplicate-comparable">SliceTricks</a>, with minor change)</th></tr></thead>
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
iterslices.Sort(in)
in = iterslices.Unique(in)
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
  algo.CopyN[int](iter.IotaReader(1), 100, iter.ChanWriter(ch))
  close(ch)
}()
fmt.Println(algo.Accumulate(iter.ChanReader(ch), nil, 0))
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
var sb strs.StringBuilderInserter[byte]
algo.UniqueCopyIf(strs.Begin(str), strs.End(str), &sb,
  func(x, y byte) bool { return x == ' ' && y == ' ' })
fmt.Println(sb.String())
// Output:
// a quick brown fox
```

</td></tr></tbody>

<thead><tr><th colspan="2">Collect N maximum elements from a channel</th></tr></thead>
<tbody><tr><td>

```go
// Need to manually maintain a min-heap.
```

</td><td>

```go
top := make([]int, 5)
algo.PartialSortCopyBy(iter.ChanReader(ch), nil, iterslices.Begin(top), iterslices.End(top),
  func(x, y int) bool { return x > y })
algo.Copy(iterslices.Begin(top), iterslices.End(top), iter.IOWriter[int](os.Stdout, ", "))
```

</td></tr></tbody>

<thead><tr><th colspan="2">Print all permutations of ["a", "b", "c"]</th></tr></thead>
<tbody><tr><td>

```go
// Usually requires some sort of recursion.
```

</td><td>

```go
s := []string{"a", "b", "c"}
for ok := true; ok; ok = iterslices.NextPermutation(s) {
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
