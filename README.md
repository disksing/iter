# iter

Generic Go implementations of the C++ standard algorithms.

`iter` provides type-safe iterator contracts, more than 100 reusable
algorithms, and convenient adapters for slices, strings, channels, and
`container/list`.

> **Development status:** the generic v2 API is being prepared and has not been
> released. Use the [`v1` branch](https://github.com/disksing/iter/tree/v1) for
> the previous non-generic version.

[简体中文](README_ZH.md)

[![CI](https://github.com/disksing/iter/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/disksing/iter/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/disksing/iter/branch/master/graph/badge.svg)](https://codecov.io/gh/disksing/iter)

## Requirements

Go 1.26 or later.

## Packages

| Package | Purpose |
|---|---|
| `iter` | Iterator contracts and adapters for generated values, channels, and `io.Writer` |
| `iter/algo` | Generic algorithms over iterator ranges |
| `iter/slices` | Direct, type-inferred APIs for common slice algorithms |
| `iter/lists` | Typed iterator adapters for `container/list` |
| `iter/strs` | Byte iterators and string output helpers |

The `slices` package is the simplest entry point for slice data. The `algo`
package exposes the full iterator-based API for custom containers and mixed
input/output sources.

## Examples

### Sort, deduplicate, and search a slice

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

Output:

```text
[1 2 3 4]
true
2
```

### Run an algorithm across different containers

```go
package main

import (
	"bytes"
	"container/list"
	"fmt"

	"github.com/disksing/iter/v2"
	"github.com/disksing/iter/v2/algo"
	iterlists "github.com/disksing/iter/v2/lists"
)

func main() {
	src := list.New()
	algo.GenerateN(iterlists.ListBackInserter[int](src), 5, iter.IotaGenerator(1))

	var out bytes.Buffer
	algo.Copy(
		iterlists.Begin[int](src),
		iterlists.End[int](src),
		iter.IOWriter[int](&out, ", "),
	)
	fmt.Println(out.String())
}
```

Output:

```text
1, 2, 3, 4, 5
```

To collect values in a slice instead, use an appender:

```go
var dst []int
algo.Copy(
	iterlists.Begin[int](src),
	iterlists.End[int](src),
	iterslices.Appender(&dst),
)
```

### Read from a channel

```go
ch := make(chan int)
go func() {
	algo.CopyN[int](iter.IotaReader(1), 100, iter.ChanWriter(ch))
	close(ch)
}()

sum := algo.Accumulate(iter.ChanReader(ch), nil, 0)
fmt.Println(sum) // 5050
```

More executable examples are in [`examples_test.go`](examples_test.go).

## Iterator model

Algorithms accept half-open ranges `[first, last)`. Iterator capabilities are
expressed by generic interfaces:

- `InputIter` reads and moves forward once.
- `ForwardReader` supports multiple passes.
- `BidiReadWriter` moves both directions and can modify values.
- `RandomReadWriter` adds constant-time distance and indexed movement.
- `OutputIter` writes values to a destination.

Concrete iterator types are exported so callers can use them in their own
fields and signatures:

- `slices.Iterator[T]`
- `lists.Iterator[T]`
- `strs.Iterator`
- `iter.IotaIterator[T]`, `iter.RepeatIterator[T]`
- `iter.ChannelReader[T]`, `iter.ChannelWriter[T]`
- `iter.OutputWriter[T]`

Default algorithms use operators when their type constraint permits it.
Variants ending in `By` accept a predicate, equality function, or ordering
function for custom types.

## Standard-library iterators

Go's standard `iter.Seq` is a good fit for read-only and lazy sequences, but it
does not directly model writable, bidirectional, or random-access positions.
This project keeps cursor-style algorithms for operations such as
`NthElement`, heap manipulation, in-place partitioning, and permutations.
Interoperability with standard `iter.Seq` is planned separately from the
current API hardening work.

## Development

```sh
go mod verify
go vet ./...
go test -race ./...
go test ./algo -run=^$ -fuzz=^FuzzSort$ -fuzztime=10s
go test ./algo -run=^$ -fuzz=^FuzzNthElement$ -fuzztime=10s
```

## License

BSD 3-Clause. See [LICENSE](LICENSE).
