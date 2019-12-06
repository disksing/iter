package iter_test

import (
	"container/list"
	"fmt"
	"os"

	. "github.com/disksing/iter"
)

// shortcuts to make life easier.
var (
	begin  = SliceBegin
	end    = SliceEnd
	lBegin = ListBegin
	lEnd   = ListEnd
	sBegin = StringBegin
	sEnd   = StringEnd
)

// Reverse a string.
func ExampleMakeString() {
	s := "!dlrow olleH"
	fmt.Println(MakeString(StringRBegin(s), StringREnd(s)))
	// Output:
	// Hello world!
}

// Print all list items to console.
func ExampleIOWriter() {
	l := list.New()
	GenerateN(ListBackInserter(l), 5, IotaGenerator(1))
	Copy(lBegin(l), lEnd(l), IOWriter(os.Stdout, ", "))
	// Output:
	// 1, 2, 3, 4, 5
}

// Sum all integers received from a channel.
func ExampleChanReader() {
	ch := make(chan int)
	go func() {
		CopyN(IotaReader(1), 100, ChanWriter(ch))
		close(ch)
	}()
	fmt.Println(Accumulate(ChanReader(ch), ChanEOF, 0))
	// Output:
	// 5050
}

// Remove consecutive spaces in a string.
func ExampleUniqueCopyIf() {
	str := "  a  quick   brown  fox    jumps  over the   lazy dog.  "
	var sb StringBuilderInserter
	UniqueCopyIf(sBegin(str), sEnd(str), &sb,
		func(x, y Any) bool { return x.(byte) == ' ' && y.(byte) == ' ' })
	fmt.Println(sb.String())
	// Output:
	// a quick brown fox jumps over the lazy dog.
}

// Collect N maximum elements from a channel.
func ExamplePartialSortCopyBy() {
	ch := make(chan int)
	go func() {
		n := make([]int, 100)
		Iota(begin(n), end(n), 1)
		Shuffle(begin(n), end(n), r)
		Copy(begin(n), end(n), ChanWriter(ch))
		close(ch)
	}()
	top := make([]int, 5)
	PartialSortCopyBy(ChanReader(ch), ChanEOF, begin(top), end(top),
		func(x, y Any) bool { return x.(int) > y.(int) })
	Copy(begin(top), end(top), IOWriter(os.Stdout, ", "))
	// Output:
	// 100, 99, 98, 97, 96
}

// Print all permutations of ["a", "b", "c"].
func ExampleNextPermutation() {
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
}
