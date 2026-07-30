package iter_test

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"

	"github.com/disksing/iter/v2"
	"github.com/disksing/iter/v2/algo"
	"github.com/disksing/iter/v2/lists"
	"github.com/disksing/iter/v2/slices"
	"github.com/disksing/iter/v2/strs"
)

// Sort, deduplicate, and search a slice.
func Example_sliceAlgorithms() {
	values := []int{3, 2, 1, 4, 3, 2, 1}
	slices.Sort(values)
	values = slices.Unique(values)

	fmt.Println(values)
	fmt.Println(slices.BinarySearch(values, 3))
	fmt.Println(slices.LowerBound(values, 3))
	// Output:
	// [1 2 3 4]
	// true
	// 2
}

// Print all list items.
func ExampleIOWriter() {
	l := list.New()
	algo.GenerateN(lists.ListBackInserter[int](l), 5, iter.IotaGenerator(1))

	var out bytes.Buffer
	algo.Copy(lists.Begin[int](l), lists.End[int](l), iter.IOWriter[int](&out, "->"))
	fmt.Println(out.String())
	// Output:
	// 1->2->3->4->5
}

// Copy values between different container types.
func ExampleAppender() {
	src := list.New()
	algo.GenerateN(lists.ListBackInserter[int](src), 5, iter.IotaGenerator(1))

	var dst []int
	algo.Copy(lists.Begin[int](src), lists.End[int](src), slices.Appender(&dst))
	fmt.Println(dst)
	// Output:
	// [1 2 3 4 5]
}

// Reverse a string.
func ExampleMakeString() {
	s := "!dlrow olleH"
	fmt.Println(strs.MakeString[byte](strs.RBegin(s), strs.REnd(s)))
	b := []byte(s)
	algo.Reverse[byte](slices.Begin(b), slices.End(b))
	fmt.Println(string(b))
	// Output:
	// Hello world!
	// Hello world!
}

// Deduplicate elements.
func ExampleUnique() {
	in := []int{3, 2, 1, 4, 3, 2, 1, 4, 1}
	algo.Sort[int](slices.Begin(in), slices.End(in))
	in = in[:slices.Begin(in).Distance(algo.Unique[int](slices.Begin(in), slices.End(in)))]
	fmt.Println(in)
	// Output:
	// [1 2 3 4]
}

// Sum all integers received from a channel.
func ExampleChanReader() {
	ch := make(chan int)
	go func() {
		algo.CopyN[int](iter.IotaReader(1), 100, iter.ChanWriter(ch))
		close(ch)
	}()
	fmt.Println(algo.Accumulate(iter.ChanReader(ch), nil, 0))
	// Output:
	// 5050
}

// Remove consecutive spaces in a string.
func ExampleUniqueCopyIf() {
	str := "  a  quick   brown  fox  "
	var sb strs.StringBuilderInserter[byte]
	algo.UniqueCopyIf(strs.Begin(str), strs.End(str), &sb,
		func(x, y byte) bool { return x == ' ' && y == ' ' })
	fmt.Println(sb.String())
	// Output:
	// a quick brown fox
}

// Collect N maximum elements from a channel.
func ExamplePartialSortCopyBy() {
	ch := make(chan int)
	go func() {
		values := make([]int, 100)
		algo.Iota(slices.Begin(values), slices.End(values), 1)
		slices.Shuffle(values, rand.New(rand.NewSource(1)))
		algo.Copy(slices.Begin(values), slices.End(values), iter.ChanWriter(ch))
		close(ch)
	}()

	top := make([]int, 5)
	algo.PartialSortCopyBy(iter.ChanReader(ch), nil, slices.Begin(top), slices.End(top),
		func(x, y int) bool { return x > y })
	fmt.Println(top)
	// Output:
	// [100 99 98 97 96]
}

// Print all permutations of ["a", "b", "c"].
func ExampleNextPermutation() {
	s := []string{"a", "b", "c"}
	for ok := true; ok; ok = algo.NextPermutation[string](slices.Begin(s), slices.End(s)) {
		fmt.Println(s)
	}
	// Output:
	// [a b c]
	// [a c b]
	// [b a c]
	// [b c a]
	// [c a b]
	// [c b a]
}
