# iter

C++ STL 迭代器和算法库的 Go 语言实现。

更少的手写循环，更多富有表达力的代码。

[![GoDoc](https://godoc.org/github.com/disksing/iter?status.svg)](https://godoc.org/github.com/disksing/iter)
[![Build Status](https://travis-ci.com/disksing/iter.svg?branch=master)](https://travis-ci.com/disksing/iter)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)
[![Go Report Card](https://goreportcard.com/badge/github.com/disksing/iter)](https://goreportcard.com/report/github.com/disksing/iter)

## 动机

虽然 Go 不支持泛型，我们值得拥有可复用的通用算法。`iter` 可以在以下方面帮助改善代码：

- 一些简单的循环逻辑不太可能写错或者低效，但使用算法调用将**使得代码更简洁和易于理解**。例如 [AllOf](https://godoc.org/github.com/disksing/iter#AllOf)，[FindIf](https://godoc.org/github.com/disksing/iter#FindIf)，[Accumulate](https://godoc.org/github.com/disksing/iter#Accumulate)。

- 一些算法并不是很复杂，但不太容易写对。使用算法库**让代码“一眼看上去就是对的”**。例如 [Shuffle](https://godoc.org/github.com/disksing/iter#Shuffle)，[Sample](https://godoc.org/github.com/disksing/iter#Sample)，[Partition](https://godoc.org/github.com/disksing/iter#Partition)。

- STL 还包含一些复杂算法，可能需要数小时才能搞对。**手动实现它们并不现实**。例如 [NthElement](https://godoc.org/github.com/disksing/iter#NthElement)，[StablePartition](https://godoc.org/github.com/disksing/iter#StablePartition)，[NextPermutation](https://godoc.org/github.com/disksing/iter#NextPermutation)。

- STL 的实现还有一些**鲜为人知的性能优化**。比如，[MinmaxElement](https://godoc.org/github.com/disksing/iter#MinmaxElement) 被实现为每次取两个元素进行比较，这样做可以大幅减少整体比较次数。

有一些开源项目在做类似的事情，比如 [gostl](https://github.com/liyue201/gostl)，[gods](https://github.com/emirpasic/gods) 和 [go-stp](https://github.com/itrabbit/go-stp)。`iter` 的独特之处在于：

- **非侵入性**。`iter` 避免重复造轮子，尽可能地复用 Go 里已有的容器（slice，string，list.List 等），使用迭代器将它们适配给算法库。

- **完整的算法库（>100）**。它实现了几乎所有 C++17 之前的算法。在[这里](https://godoc.org/github.com/disksing/iter)可以查看完整列表。

## 示例

> 这些示例给一些函数定义了别名来使代码更美观，详情见 [example_test.go](https://github.com/disksing/iter/blob/master/examples_test.go)。

<table>
<thead><tr><th colspan="2">控制台输出 list.List</th></tr></thead>
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

<thead><tr><th colspan="2">反转 string</th></tr></thead>
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

<thead><tr><th colspan="2">去重（来自 <a href="https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable">SliceTricks</a>，略微调整）</th></tr></thead>
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

<thead><tr><th colspan="2">对 channel 中的所有整数求和</th></tr></thead>
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

<thead><tr><th colspan="2">删除字符串中的连续空格</th></tr></thead>
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

<thead><tr><th colspan="2">收集 channel 中最大的 N 个整数</th></tr></thead>
<tbody><tr><td>

```go
// 需要手动维护小顶堆。
```

</td><td>

```go
top := make([]int, 5)
PartialSortCopyBy(ChanReader(ch), ChanEOF, begin(top), end(top),
  func(x, y Any) bool { return x.(int) > y.(int) })
Copy(begin(top), end(top), IOWriter(os.Stdout, ", "))
```

</td></tr></tbody>

<thead><tr><th colspan="2">输出 ["a", "b", "c"] 的所有排列</th></tr></thead>
<tbody><tr><td>

```go
// 通常需要引入递归来完成。
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

## 致谢

- [cppreference.com](https://en.cppreference.com/)
- [LLVM libc++](https://libcxx.llvm.org/)

## 开源许可证

BSD 3-Clause
