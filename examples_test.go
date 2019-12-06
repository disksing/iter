package iter_test

import (
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"time"

	. "github.com/disksing/iter"
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
	Copy(ListBegin(l), ListEnd(l), IOWriter(os.Stdout, ", "))
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
	UniqueCopyIf(StringBegin(str), StringEnd(str), &sb,
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
		Iota(SliceBegin(n), SliceEnd(n), 1)
		Shuffle(SliceBegin(n), SliceEnd(n), rand.New(rand.NewSource(time.Now().UnixNano())))
		Copy(SliceBegin(n), SliceEnd(n), ChanWriter(ch))
		close(ch)
	}()
	top := make([]int, 5)
	PartialSortCopyBy(ChanReader(ch), ChanEOF, SliceBegin(top), SliceEnd(top),
		func(x, y Any) bool { return x.(int) > y.(int) })
	Copy(SliceBegin(top), SliceEnd(top), IOWriter(os.Stdout, ", "))
	// Output:
	// 100, 99, 98, 97, 96
}

// Print all permutations of ["a", "b", "c"].
func ExampleNextPermutation() {
	s := []string{"a", "b", "c"}
	for ok := true; ok; ok = NextPermutation(SliceBegin(s), SliceEnd(s)) {
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
