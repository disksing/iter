package slices_test

import (
	"slices"
	"strings"
	"testing"

	iter "github.com/disksing/iter/v2"
	iterslices "github.com/disksing/iter/v2/slices"
)

var _ iter.RandomReadWriter[int, iterslices.Iterator[int]] = iterslices.Iterator[int]{}

func TestSearchFacade(t *testing.T) {
	values := []int{1, 2, 2, 3, 2, 2, 3, 4}

	if got := iterslices.FindEnd(values, []int{2, 2, 3}); got != 4 {
		t.Fatalf("FindEnd() = %d, want 4", got)
	}
	if got := iterslices.FindEndBy(values, []string{"2", "3"}, func(x int, y string) bool {
		return string(rune('0'+x)) == y
	}); got != 5 {
		t.Fatalf("FindEndBy() = %d, want 5", got)
	}
	if got := iterslices.SearchN(values, 2, 2); got != 1 {
		t.Fatalf("SearchN() = %d, want 1", got)
	}
	if got := iterslices.SearchNBy(values, 2, "2", func(x int, y string) bool {
		return string(rune('0'+x)) == y
	}); got != 1 {
		t.Fatalf("SearchNBy() = %d, want 1", got)
	}
}

func TestPartitionAndSortFacade(t *testing.T) {
	partitioned := []int{2, 4, 6, 1, 3, 5}
	if got := iterslices.PartitionPoint(partitioned, func(x int) bool { return x%2 == 0 }); got != 3 {
		t.Fatalf("PartitionPoint() = %d, want 3", got)
	}

	values := []int{9, 1, 5, 3, 7, 2, 8, 6, 4}
	iterslices.Sort(values)
	if !slices.IsSorted(values) || !iterslices.IsSorted(values) {
		t.Fatalf("Sort() result is not sorted: %v", values)
	}
	if got := iterslices.IsSortedUntil([]int{1, 3, 2, 4}); got != 2 {
		t.Fatalf("IsSortedUntil() = %d, want 2", got)
	}

	desc := []int{1, 4, 2, 3}
	greater := func(a, b int) bool { return a > b }
	iterslices.SortBy(desc, greater)
	if !slices.Equal(desc, []int{4, 3, 2, 1}) || !iterslices.IsSortedBy(desc, greater) {
		t.Fatalf("SortBy() result = %v", desc)
	}
	if got := iterslices.IsSortedUntilBy([]int{4, 3, 1, 2}, greater); got != 3 {
		t.Fatalf("IsSortedUntilBy() = %d, want 3", got)
	}
}

func TestStablePartialAndNthElementFacade(t *testing.T) {
	type item struct {
		key int
		id  string
	}
	values := []item{{2, "a"}, {1, "b"}, {2, "c"}, {1, "d"}}
	iterslices.StableSortBy(values, func(a, b item) bool { return a.key < b.key })
	want := []item{{1, "b"}, {1, "d"}, {2, "a"}, {2, "c"}}
	if !slices.Equal(values, want) {
		t.Fatalf("StableSortBy() = %v, want %v", values, want)
	}

	partial := []int{8, 3, 7, 4, 2, 6, 5, 1}
	sorted := slices.Clone(partial)
	slices.Sort(sorted)
	iterslices.PartialSort(partial, 4)
	if !slices.Equal(partial[:4], sorted[:4]) {
		t.Fatalf("PartialSort() prefix = %v, want %v", partial[:4], sorted[:4])
	}

	partialDesc := []int{1, 6, 2, 5, 3, 4}
	iterslices.PartialSortBy(partialDesc, 3, func(a, b int) bool { return a > b })
	if !slices.Equal(partialDesc[:3], []int{6, 5, 4}) {
		t.Fatalf("PartialSortBy() prefix = %v", partialDesc[:3])
	}

	nth := []int{8, 3, 7, 4, 2, 6, 5, 1}
	iterslices.NthElement(nth, 3)
	if nth[3] != sorted[3] {
		t.Fatalf("NthElement() value = %d, want %d", nth[3], sorted[3])
	}
	for i := range nth {
		if i < 3 && nth[i] > nth[3] {
			t.Fatalf("NthElement() left partition invalid: %v", nth)
		}
		if i > 3 && nth[i] < nth[3] {
			t.Fatalf("NthElement() right partition invalid: %v", nth)
		}
	}

	nthDesc := []int{1, 6, 2, 5, 3, 4}
	iterslices.NthElementBy(nthDesc, 2, func(a, b int) bool { return a > b })
	if nthDesc[2] != 4 {
		t.Fatalf("NthElementBy() value = %d, want 4", nthDesc[2])
	}
}

func TestBinarySearchFacade(t *testing.T) {
	values := []int{1, 2, 2, 2, 4, 6}

	if got := iterslices.LowerBound(values, 2); got != 1 {
		t.Fatalf("LowerBound() = %d, want 1", got)
	}
	if got := iterslices.UpperBound(values, 2); got != 4 {
		t.Fatalf("UpperBound() = %d, want 4", got)
	}
	if !iterslices.BinarySearch(values, 4) || iterslices.BinarySearch(values, 3) {
		t.Fatal("BinarySearch() returned an unexpected result")
	}
	first, last := iterslices.EqualRange(values, 2)
	if first != 1 || last != 4 {
		t.Fatalf("EqualRange() = (%d, %d), want (1, 4)", first, last)
	}

	desc := []int{6, 4, 2, 2, 2, 1}
	greater := func(a, b int) bool { return a > b }
	if got := iterslices.LowerBoundBy(desc, 2, greater); got != 2 {
		t.Fatalf("LowerBoundBy() = %d, want 2", got)
	}
	if got := iterslices.UpperBoundBy(desc, 2, greater); got != 5 {
		t.Fatalf("UpperBoundBy() = %d, want 5", got)
	}
	if !iterslices.BinarySearchBy(desc, 4, greater) {
		t.Fatal("BinarySearchBy() did not find 4")
	}
	first, last = iterslices.EqualRangeBy(desc, 2, greater)
	if first != 2 || last != 5 {
		t.Fatalf("EqualRangeBy() = (%d, %d), want (2, 5)", first, last)
	}
}

func TestHeapFacade(t *testing.T) {
	values := []int{3, 1, 4, 2, 5}
	iterslices.MakeHeap(values)
	if !iterslices.IsHeap(values) || iterslices.IsHeapUntil(values) != len(values) {
		t.Fatalf("MakeHeap() result is not a heap: %v", values)
	}

	values = append(values, 6)
	iterslices.PushHeap(values)
	if !iterslices.IsHeap(values) || values[0] != 6 {
		t.Fatalf("PushHeap() result is invalid: %v", values)
	}

	iterslices.PopHeap(values)
	if values[len(values)-1] != 6 || !iterslices.IsHeap(values[:len(values)-1]) {
		t.Fatalf("PopHeap() result is invalid: %v", values)
	}

	heap := values[:len(values)-1]
	iterslices.SortHeap(heap)
	if !slices.IsSorted(heap) {
		t.Fatalf("SortHeap() result is not sorted: %v", heap)
	}

	minHeap := []int{5, 3, 4, 1, 2}
	greater := func(a, b int) bool { return a > b }
	iterslices.MakeHeapBy(minHeap, greater)
	if !iterslices.IsHeapBy(minHeap, greater) ||
		iterslices.IsHeapUntilBy(minHeap, greater) != len(minHeap) {
		t.Fatalf("MakeHeapBy() result is not a min heap: %v", minHeap)
	}
	iterslices.SortHeapBy(minHeap, greater)
	if !slices.IsSortedFunc(minHeap, func(a, b int) int { return b - a }) {
		t.Fatalf("SortHeapBy() result is not descending: %v", minHeap)
	}
}

func TestIsPermutationByFacade(t *testing.T) {
	left := []string{"A", "b", "C"}
	right := []string{"c", "B", "a"}
	if !iterslices.IsPermutationBy(left, right, strings.EqualFold) {
		t.Fatal("IsPermutationBy() = false, want true")
	}
}
