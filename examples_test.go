package iter_test

import (
	"fmt"

	"github.com/disksing/iter/v2/algo"
	"github.com/disksing/iter/v2/slices"
	"github.com/disksing/iter/v2/strs"
)

// // Print all list items to console.
// func ExampleIOWriter() {
// 	l := list.New()
// 	GenerateN(ListBackInserter(l), 5, IotaGenerator(1))
// 	Copy(lBegin(l), lEnd(l), IOWriter(os.Stdout, "->"))
// 	// Output:
// 	// 1->2->3->4->5
// }

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

// // Sum all integers received from a channel.
// func ExampleChanReader() {
// 	ch := make(chan int)
// 	go func() {
// 		CopyN(IotaReader(1), 100, ChanWriter(ch))
// 		close(ch)
// 	}()
// 	fmt.Println(Accumulate(ChanReader(ch), ChanEOF, 0))
// 	// Output:
// 	// 5050
// }

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

// // Collect N maximum elements from a channel.
// func ExamplePartialSortCopyBy() {
// 	ch := make(chan int)
// 	go func() {
// 		n := make([]int, 100)
// 		Iota(begin(n), end(n), 1)
// 		Shuffle(begin(n), end(n), r)
// 		Copy(begin(n), end(n), ChanWriter(ch))
// 		close(ch)
// 	}()
// 	top := make([]int, 5)
// 	PartialSortCopyBy(ChanReader(ch), ChanEOF, begin(top), end(top),
// 		func(x, y any) bool { return x.(int) > y.(int) })
// 	Copy(begin(top), end(top), IOWriter(os.Stdout, ", "))
// 	// Output:
// 	// 100, 99, 98, 97, 96
// }

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
