package slices

import (
	"math/rand"

	"github.com/disksing/iter/v2"
	"github.com/disksing/iter/v2/algo"
)

// AllOf checks if unary predicate pred returns true for all elements.
func AllOf[T any](s []T, pred algo.UnaryPredicate[T]) bool {
	return algo.AllOf(Begin(s), End(s), pred)
}

// AnyOf checks if unary predicate pred returns true for at least one element.
func AnyOf[T any](s []T, pred algo.UnaryPredicate[T]) bool {
	return algo.AnyOf(Begin(s), End(s), pred)
}

// NoneOf checks if unary predicate pred returns true for no elements.
func NoneOf[T any](s []T, pred algo.UnaryPredicate[T]) bool {
	return algo.NoneOf(Begin(s), End(s), pred)
}

// ForEach applies the given function f to each element of the slice.
func ForEach[T any](s []T, f algo.IteratorFunction[T]) algo.IteratorFunction[T] {
	return algo.ForEach(Begin(s), End(s), f)
}

// Count counts the elements that are equal to value.
func Count[T comparable](s []T, v T) int {
	return algo.Count(Begin(s), End(s), v)
}

// CountIf counts elements for which predicate pred returns true.
func CountIf[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return algo.CountIf(Begin(s), End(s), pred)
}

func __end_ref[T any](s []T) *Iterator[T] {
	e := End(s)
	return &e
}

func __idx[T any](it Iterator[T]) int {
	return it.i
}

// Mismatch returns the first mismatching position of the two slices.
func Mismatch[T comparable](s1, s2 []T) (int, int) {
	it1, it2 := algo.Mismatch[T](Begin(s1), End(s1), Begin(s2), __end_ref(s2))
	return __idx(it1), __idx(it2)
}

// MismatchBy returns the first mismatching position of the two slices. Elements
// are compared using the given comparer eq.
func MismatchBy[T1, T2 any](s1 []T1, s2 []T2, eq algo.EqComparer[T1, T2]) (int, int) {
	it1, it2 := algo.MismatchBy(Begin(s1), End(s1), Begin(s2), __end_ref(s2), eq)
	return __idx(it1), __idx(it2)
}

// Find returns the first position in the slice that is equal to value.
func Find[T comparable](s []T, x T) int {
	return __idx(algo.Find(Begin(s), End(s), x))
}

// FindIf returns the first position in the slice which predicate pred returns
// true.
func FindIf[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return __idx(algo.FindIf(Begin(s), End(s), pred))
}

// FindIfNot returns the first position in the slice which predicate pred returns
// false.
func FindIfNot[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return __idx(algo.FindIfNot(Begin(s), End(s), pred))
}

// FindFirstOf searches s1 for any of the elements in s2.
func FindFirstOf[T comparable](s1, s2 []T) int {
	return __idx(algo.FindFirstOf[T](Begin(s1), End(s1), Begin(s2), End(s2)))
}

// FindFirstOfBy searches s1 for any of the elements in s2.
//
// Elements are compared using the given binary comparer eq.
func FindFirstOfBy[T1, T2 any](s1 []T1, s2 []T2, eq algo.EqComparer[T1, T2]) int {
	return __idx(algo.FindFirstOfBy(Begin(s1), End(s1), Begin(s2), End(s2), eq))
}

// AdjacentFind searches the slice for two consecutive identical elements.
func AdjacentFind[T comparable](s []T) int {
	return __idx(algo.AdjacentFind[T](Begin(s), End(s)))
}

// AdjacentFindBy searches the slice for two consecutive identical elements.
//
// Elements are compared using the given binary comparer eq.
func AdjacentFindBy[T any](s []T, eq algo.EqComparer[T, T]) int {
	return __idx(algo.AdjacentFindBy(Begin(s), End(s), eq))
}

// Search searches for the first occurrence of s2 in s1.
func Search[T comparable](s1, s2 []T) int {
	return __idx(algo.Search[T](Begin(s1), End(s1), Begin(s2), End(s2)))
}

// SearchBy searches for the first occurrence of s2 in s1.
//
// Elements are compared using the given binary comparer eq.
func SearchBy[T1, T2 any](s1 []T1, s2 []T2, eq algo.EqComparer[T1, T2]) int {
	return __idx(algo.SearchBy(Begin(s1), End(s1), Begin(s2), End(s2), eq))
}

// FindEnd searches for the last occurrence of s2 in s1.
func FindEnd[T comparable](s1, s2 []T) int {
	return __idx(algo.FindEnd[T](Begin(s1), End(s1), Begin(s2), End(s2)))
}

// FindEndBy searches for the last occurrence of s2 in s1.
//
// Elements are compared using the given binary comparer eq.
func FindEndBy[T1, T2 any](s1 []T1, s2 []T2, eq algo.EqComparer[T1, T2]) int {
	return __idx(algo.FindEndBy(Begin(s1), End(s1), Begin(s2), End(s2), eq))
}

// SearchN searches for the first sequence of count elements equal to v.
func SearchN[T comparable](s []T, count int, v T) int {
	return __idx(algo.SearchN(Begin(s), End(s), count, v))
}

// SearchNBy searches for the first sequence of count matching elements.
//
// Elements are compared with v using the given binary comparer eq.
func SearchNBy[T1, T2 any](s []T1, count int, v T2, eq algo.EqComparer[T1, T2]) int {
	return __idx(algo.SearchNBy(Begin(s), End(s), count, v, eq))
}

// Remove removes all elements equal to v from the slice and returns the new
// slice.
func Remove[T comparable](s []T, v T) []T {
	return s[:__idx(algo.Remove(Begin(s), End(s), v))]
}

// RemoveIf removes all elements equal to v from the slice and returns the new
// slice.
func RemoveIf[T any](s []T, pred algo.UnaryPredicate[T]) []T {
	return s[:__idx(algo.RemoveIf(Begin(s), End(s), pred))]
}

// Replace replaces all elements equal to old with new in the slice.
func Replace[T comparable](s []T, old, new T) {
	algo.Replace(Begin(s), End(s), old, new)
}

// ReplaceIf replaces all elements satisfy pred with new in the slice.
func ReplaceIf[T any](s []T, pred algo.UnaryPredicate[T], v T) {
	algo.ReplaceIf(Begin(s), End(s), pred, v)
}

// SwapRanges exchanges elements between two slices. The length of the second
// slice should be not less than the first.
func SwapRanges[T any](s1, s2 []T) {
	algo.SwapRanges[T](Begin(s1), End(s1), Begin(s2))
}

// Reverse reverses the order of the elements in the slice.
func Reverse[T any](s []T) {
	algo.Reverse[T](Begin(s), End(s))
}

// Rotate performs a left rotation on a slice in such a way, that the element
// newFirst becomes the first element of the new range and newFirst - 1 becomes
// the last element.
func Rotate[T any](s []T, newFirst int) int {
	return __idx(algo.Rotate[T](Begin(s), Begin(s).AdvanceN(newFirst), End(s)))
}

// Shuffle reorders the elements in the list such that each possible permutation
// of those elements has equal probability of appearance.
func Shuffle[T any](s []T, r *rand.Rand) {
	algo.Shuffle[T](Begin(s), End(s), r)
}

// Unique eliminates all but the first element from every consecutive group of
// equivalent elements from the slice.
func Unique[T comparable](s []T) []T {
	return s[:__idx(algo.Unique[T](Begin(s), End(s)))]
}

// UniqueIf eliminates all but the first element from every consecutive group of
// equivalent elements from the slice.
//
// Elements are compared using the given binary comparer eq.
func UniqueIf[T any](s []T, eq algo.EqComparer[T, T]) []T {
	return s[:__idx(algo.UniqueIf(Begin(s), End(s), eq))]
}

// IsPartitioned returns true if all elements in the slice that satisfy the
// predicate pred appear before all elements that don't. Also returns true if the
// slice is empty.
func IsPartitioned[T any](s []T, pred algo.UnaryPredicate[T]) bool {
	return algo.IsPartitioned(Begin(s), End(s), pred)
}

// Partition reorders the elements in the slice in such a way that
// all elements for which the predicate pred returns true precede the elements
// for which predicate pred returns false.
//
// Relative order of the elements is not preserved.
func Partition[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return __idx(algo.Partition(Begin(s), End(s), pred))
}

// StablePartition reorders the elements in the slice in such a way that all
// elements for which the predicate pred returns true precede the elements for
// which predicate pred returns false. Relative order of the elements is
// preserved.
func StablePartition[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return __idx(algo.StablePartitionBidi(Begin(s), End(s), pred))
}

// PartitionPoint returns the first position for which pred is false in an
// already partitioned slice.
func PartitionPoint[T any](s []T, pred algo.UnaryPredicate[T]) int {
	return __idx(algo.PartitionPoint(Begin(s), End(s), pred))
}

// IsSorted reports whether s is sorted in ascending order.
func IsSorted[T iter.Ordered](s []T) bool {
	return algo.IsSorted[T](Begin(s), End(s))
}

// IsSortedBy reports whether s is sorted according to less.
func IsSortedBy[T any](s []T, less algo.LessComparer[T]) bool {
	return algo.IsSortedBy(Begin(s), End(s), less)
}

// IsSortedUntil returns the first position at which s stops being sorted.
// It returns len(s) when s is sorted.
func IsSortedUntil[T iter.Ordered](s []T) int {
	return __idx(algo.IsSortedUntil[T](Begin(s), End(s)))
}

// IsSortedUntilBy returns the first position at which s stops being sorted
// according to less. It returns len(s) when s is sorted.
func IsSortedUntilBy[T any](s []T, less algo.LessComparer[T]) int {
	return __idx(algo.IsSortedUntilBy(Begin(s), End(s), less))
}

// Sort sorts s in ascending order.
func Sort[T iter.Ordered](s []T) {
	algo.Sort[T](Begin(s), End(s))
}

// SortBy sorts s according to less.
func SortBy[T any](s []T, less algo.LessComparer[T]) {
	algo.SortBy(Begin(s), End(s), less)
}

// StableSort stably sorts s in ascending order.
func StableSort[T iter.Ordered](s []T) {
	algo.StableSort[T](Begin(s), End(s))
}

// StableSortBy stably sorts s according to less.
func StableSortBy[T any](s []T, less algo.LessComparer[T]) {
	algo.StableSortBy(Begin(s), End(s), less)
}

// PartialSort rearranges s so that s[:middle] contains the smallest middle
// elements in ascending order. middle must be in [0, len(s)].
func PartialSort[T iter.Ordered](s []T, middle int) {
	algo.PartialSort[T](Begin(s), Begin(s).AdvanceN(middle), End(s))
}

// PartialSortBy rearranges s so that s[:middle] contains the first middle
// elements according to less, in sorted order.
func PartialSortBy[T any](s []T, middle int, less algo.LessComparer[T]) {
	algo.PartialSortBy(Begin(s), Begin(s).AdvanceN(middle), End(s), less)
}

// NthElement rearranges s so the element at nth is the one that would occur
// there in sorted order. nth must be in [0, len(s)).
func NthElement[T iter.Ordered](s []T, nth int) {
	algo.NthElement[T](Begin(s), Begin(s).AdvanceN(nth), End(s))
}

// NthElementBy rearranges s so the element at nth is the one that would occur
// there when ordered by less.
func NthElementBy[T any](s []T, nth int, less algo.LessComparer[T]) {
	algo.NthElementBy(Begin(s), Begin(s).AdvanceN(nth), End(s), less)
}

// LowerBound returns the first position at which v could be inserted without
// violating ascending order.
func LowerBound[T iter.Ordered](s []T, v T) int {
	return __idx(algo.LowerBound[T](Begin(s), End(s), v))
}

// LowerBoundBy returns the first insertion position for v according to less.
func LowerBoundBy[T any](s []T, v T, less algo.LessComparer[T]) int {
	return __idx(algo.LowerBoundBy(Begin(s), End(s), v, less))
}

// UpperBound returns the first position after all elements equal to v.
func UpperBound[T iter.Ordered](s []T, v T) int {
	return __idx(algo.UpperBound[T](Begin(s), End(s), v))
}

// UpperBoundBy returns the first position after all elements equivalent to v
// according to less.
func UpperBoundBy[T any](s []T, v T, less algo.LessComparer[T]) int {
	return __idx(algo.UpperBoundBy(Begin(s), End(s), v, less))
}

// BinarySearch reports whether sorted s contains v.
func BinarySearch[T iter.Ordered](s []T, v T) bool {
	return algo.BinarySearch[T](Begin(s), End(s), v)
}

// BinarySearchBy reports whether sorted s contains v according to less.
func BinarySearchBy[T any](s []T, v T, less algo.LessComparer[T]) bool {
	return algo.BinarySearchBy(Begin(s), End(s), v, less)
}

// EqualRange returns the lower and upper bounds of v in sorted s.
func EqualRange[T iter.Ordered](s []T, v T) (int, int) {
	first, last := algo.EqualRange[T](Begin(s), End(s), v)
	return __idx(first), __idx(last)
}

// EqualRangeBy returns the lower and upper bounds of v according to less.
func EqualRangeBy[T any](s []T, v T, less algo.LessComparer[T]) (int, int) {
	first, last := algo.EqualRangeBy(Begin(s), End(s), v, less)
	return __idx(first), __idx(last)
}

// MaxElement returns the position largest element in a slice.
func MaxElement[T iter.Ordered](s []T) int {
	return __idx(algo.MaxElement[T](Begin(s), End(s)))
}

// MaxElementBy returns the position of the largest element in a slice.
//
// Values are compared using the given binary comparer less.
func MaxElementBy[T any](s []T, less algo.LessComparer[T]) int {
	return __idx(algo.MaxElementBy(Begin(s), End(s), less))
}

// MinElement returns the position smallest element in a slice.
func MinElement[T iter.Ordered](s []T) int {
	return __idx(algo.MinElement[T](Begin(s), End(s)))
}

// MinElementBy returns the position of the smallest element in a slice.
//
// Values are compared using the given binary comparer less.
func MinElementBy[T any](s []T, less algo.LessComparer[T]) int {
	return __idx(algo.MinElementBy(Begin(s), End(s), less))
}

// MinmaxElement returns the positions of the smallest and largest elements.
func MinmaxElement[T iter.Ordered](s []T) (int, int) {
	a, b := algo.MinmaxElement[T](Begin(s), End(s))
	return __idx(a), __idx(b)
}

// MinmaxElementBy returns the position of the smallest and the latest element in
// a slice.
//
// Values are compared using the given binary comparer less.
func MinmaxElementBy[T any](s []T, less algo.LessComparer[T]) (int, int) {
	a, b := algo.MinmaxElementBy(Begin(s), End(s), less)
	return __idx(a), __idx(b)
}

// Equal returns true if two slices are equal.
func Equal[T comparable](s1, s2 []T) bool {
	return algo.Equal[T](Begin(s1), End(s1), Begin(s2), __end_ref(s2))
}

// EqualBy returns true if two slices are equal.
//
// Elements are compared using the given binary comparer eq.
func EqualBy[T any](s1, s2 []T, eq algo.EqComparer[T, T]) bool {
	return algo.EqualBy(Begin(s1), End(s1), Begin(s2), __end_ref(s2), eq)
}

// LexicographicalCompare checks if the first slice is
// lexicographically less than the second slice.
func LexicographicalCompare[T iter.Ordered](s1, s2 []T) bool {
	return algo.LexicographicalCompare[T](Begin(s1), End(s1), Begin(s2), End(s2))
}

// LexicographicalCompareBy checks if the first slice is lexicographically less
// than the second slice.
//
// Elements are compared using the given binary comparer less.
func LexicographicalCompareBy[T any](s1, s2 []T, less algo.LessComparer[T]) bool {
	return algo.LexicographicalCompareBy(Begin(s1), End(s1), Begin(s2), End(s2), less)
}

// LexicographicalCompareThreeWay lexicographically compares two ranges s1 and s2
// using three-way comparison. The result will be 0 if s1 == s2, -1 if s1 < s2 ,
// 1 if s1 > s2.
func LexicographicalCompareThreeWay[T iter.Ordered](s1, s2 []T) int {
	return algo.LexicographicalCompareThreeWay[T](Begin(s1), End(s1), Begin(s2), End(s2))
}

// LexicographicalCompareThreeWayBy lexicographically compares two ranges s1 and
// s2 using three-way comparison. The result will be 0 if s1 == s2, -1 if s1 < s2
// , 1 if s1 > s2.
//
// Elements are compared using the given binary predicate cmp.
func LexicographicalCompareThreeWayBy[T any](s1, s2 []T, cmp algo.ThreeWayComparer[T, T]) int {
	return algo.LexicographicalCompareThreeWayBy(Begin(s1), End(s1), Begin(s2), End(s2), cmp)
}

// IsPermutation returns true if there exists a permutation of the elements in
// the slice s1 that makes that range equal to s2.
func IsPermutation[T comparable](s1, s2 []T) bool {
	return algo.IsPermutation[T](Begin(s1), End(s1), Begin(s2), __end_ref(s2))
}

// IsPermutationBy reports whether s1 is a permutation of s2 according to eq.
func IsPermutationBy[T any](s1, s2 []T, eq algo.EqComparer[T, T]) bool {
	return algo.IsPermutationBy(Begin(s1), End(s1), Begin(s2), __end_ref(s2), eq)
}

// NextPermutation transforms the slice into the next permutation from the set of
// all permutations that are lexicographically ordered. Returns true if such
// permutation exists, otherwise transforms the range into the first permutation
// (as if by Sort(s)) and returns false.
func NextPermutation[T iter.Ordered](s []T) bool {
	return algo.NextPermutation[T](Begin(s), End(s))
}

// PrevPermutation transforms the slice into the previous permutation from the
// set of all permutations that are lexicographically ordered. Returns true if
// such permutation exists, otherwise transforms the range into the last
// permutation (as if by Sort(s) + Reverse(s)) and returns false.
func PrevPermutation[T iter.Ordered](s []T) bool {
	return algo.PrevPermutation[T](Begin(s), End(s))
}

// Accumulate computes the sum of the given value v and the elements in the
// slice, using v+=x.
func Accumulate[T iter.Numeric](s []T, v T) T {
	return algo.Accumulate(Begin(s), End(s), v)
}

// InnerProduct computes inner product (i.e. sum of products) or performs ordered
// map/reduce operation on the slice, using v=v+x*y.
func InnerProduct[T iter.Numeric](s1, s2 []T, v T) T {
	return algo.InnerProduct(Begin(s1), End(s1), Begin(s2), v)
}

// IsHeap reports whether s is a max heap.
func IsHeap[T iter.Ordered](s []T) bool {
	return algo.IsHeap[T](Begin(s), End(s))
}

// IsHeapBy reports whether s is a heap according to less.
func IsHeapBy[T any](s []T, less algo.LessComparer[T]) bool {
	return algo.IsHeapBy(Begin(s), End(s), less)
}

// IsHeapUntil returns the first position at which s stops being a max heap.
// It returns len(s) when s is a heap.
func IsHeapUntil[T iter.Ordered](s []T) int {
	return __idx(algo.IsHeapUntil[T](Begin(s), End(s)))
}

// IsHeapUntilBy returns the first position at which s stops being a heap
// according to less.
func IsHeapUntilBy[T any](s []T, less algo.LessComparer[T]) int {
	return __idx(algo.IsHeapUntilBy(Begin(s), End(s), less))
}

// MakeHeap rearranges s into a max heap.
func MakeHeap[T iter.Ordered](s []T) {
	algo.MakeHeap[T](Begin(s), End(s))
}

// MakeHeapBy rearranges s into a heap according to less.
func MakeHeapBy[T any](s []T, less algo.LessComparer[T]) {
	algo.MakeHeapBy(Begin(s), End(s), less)
}

// PushHeap adds the final element of s to the max heap in s[:len(s)-1].
func PushHeap[T iter.Ordered](s []T) {
	algo.PushHeap[T](Begin(s), End(s))
}

// PushHeapBy adds the final element of s to the heap in s[:len(s)-1].
func PushHeapBy[T any](s []T, less algo.LessComparer[T]) {
	algo.PushHeapBy(Begin(s), End(s), less)
}

// PopHeap moves the largest element to the end of s and restores the max heap
// in s[:len(s)-1].
func PopHeap[T iter.Ordered](s []T) {
	algo.PopHeap[T](Begin(s), End(s))
}

// PopHeapBy moves the first heap element to the end of s according to less.
func PopHeapBy[T any](s []T, less algo.LessComparer[T]) {
	algo.PopHeapBy(Begin(s), End(s), less)
}

// SortHeap turns a max heap into an ascending sorted slice.
func SortHeap[T iter.Ordered](s []T) {
	algo.SortHeap[T](Begin(s), End(s))
}

// SortHeapBy turns a heap into a sorted slice according to less.
func SortHeapBy[T any](s []T, less algo.LessComparer[T]) {
	algo.SortHeapBy(Begin(s), End(s), less)
}
