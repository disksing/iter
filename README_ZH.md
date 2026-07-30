# iter

C++ 标准算法的 Go 泛型实现。

`iter` 提供类型安全的迭代器约束、100 多个可复用算法，以及 slice、
string、channel 和 `container/list` 的便捷适配器。

> **开发状态：** 泛型 v2 API 正在整理，尚未发布。此前可用的非泛型版本
> 位于 [`v1` 分支](https://github.com/disksing/iter/tree/v1)。

[English](README.md)

[![CI](https://github.com/disksing/iter/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/disksing/iter/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)

## 环境要求

Go 1.26 或更高版本。

## 包结构

| 包 | 用途 |
|---|---|
| `iter` | 迭代器约束，以及生成值、channel、`io.Writer` 等适配器 |
| `iter/algo` | 面向迭代器区间的完整泛型算法 |
| `iter/slices` | 常用 slice 算法，可直接依赖类型推断调用 |
| `iter/lists` | `container/list` 的类型化迭代器适配器 |
| `iter/strs` | 字节级 string 迭代器和字符串输出工具 |

处理 slice 时优先使用 `slices` 包；需要自定义容器、跨容器复制或完整算法
集合时使用 `algo` 包。

## 示例

### 排序、去重和查找

```go
package main

import (
	"fmt"

	iterslices "github.com/disksing/iter/v2/slices"
)

func main() {
	values := []int{3, 2, 1, 4, 3, 2, 1}
	iterslices.Sort(values)
	values = iterslices.Unique(values)

	fmt.Println(values)
	fmt.Println(iterslices.BinarySearch(values, 3))
	fmt.Println(iterslices.LowerBound(values, 3))
}
```

输出：

```text
[1 2 3 4]
true
2
```

### 在不同容器之间运行算法

```go
src := list.New()
algo.GenerateN(iterlists.ListBackInserter[int](src), 5, iter.IotaGenerator(1))

var dst []int
algo.Copy(
	iterlists.Begin[int](src),
	iterlists.End[int](src),
	iterslices.Appender(&dst),
)
fmt.Println(dst) // [1 2 3 4 5]
```

### 从 channel 读取

```go
ch := make(chan int)
go func() {
	algo.CopyN[int](iter.IotaReader(1), 100, iter.ChanWriter(ch))
	close(ch)
}()

sum := algo.Accumulate(iter.ChanReader(ch), nil, 0)
fmt.Println(sum) // 5050
```

更多可执行示例见 [`examples_test.go`](examples_test.go)。

## 迭代器模型

算法使用左闭右开区间 `[first, last)`。泛型接口表达不同的迭代能力：

- `InputIter`：单次向前读取；
- `ForwardReader`：允许多次遍历；
- `BidiReadWriter`：可双向移动并修改值；
- `RandomReadWriter`：额外支持常数时间的距离和按偏移移动；
- `OutputIter`：向目标写值。

具体类型均已导出，调用者可以在自己的字段和函数签名中使用：

- `slices.Iterator[T]`
- `lists.Iterator[T]`
- `strs.Iterator`
- `iter.IotaIterator[T]`、`iter.RepeatIterator[T]`
- `iter.ChannelReader[T]`、`iter.ChannelWriter[T]`
- `iter.OutputWriter[T]`

默认算法在类型约束允许时使用运算符；以 `By` 结尾的变体接收自定义 predicate、
相等比较器或顺序比较器。

## 与标准库 iter 的关系

标准库 `iter.Seq` 很适合只读和惰性序列，但不能直接表达可写、双向或随机访问
的位置。`NthElement`、heap、原地 partition、permutation 等正是本项目保留
cursor 模型的原因。与标准 `iter.Seq` 的互操作会作为后续工作单独设计，不与
当前 API 加固混在一起。

## 开发验证

```sh
go mod verify
go vet ./...
go test -race ./...
go test ./algo -run=^$ -fuzz=^FuzzSort$ -fuzztime=10s
go test ./algo -run=^$ -fuzz=^FuzzNthElement$ -fuzztime=10s
```

## 许可证

BSD 3-Clause，详见 [LICENSE](LICENSE)。
