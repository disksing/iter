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

func __end_ref[T any](s []T) *sliceIter[T] {
	e := End(s)
	return &e
}

func __idx[T any](it sliceIter[T]) int {
	if it.i == len(it.s) {
		return -1
	}
	return it.i
}

// Mismatch returns the first mismatching position of the two slices. If a slice
// is prefix of the other, it returns -1.
func Mismatch[T comparable](s1, s2 []T) (int, int) {
	it1, it2 := algo.Mismatch[T](Begin(s1), End(s1), Begin(s2), __end_ref(s2))
	return __idx(it1), __idx(it2)
}

// MismatchBy returns the first mismatching position of the two slices. If a
// slice is prefix of the other, it returns -1. Elements are compared using the
// given comparer eq.
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

// MinmaxElement returns the position smallest and the lagest element in a slice.
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
func IsPermutation[T iter.Ordered](s1, s2 []T) bool {
	return algo.IsPermutation[T](Begin(s1), End(s1), Begin(s2), __end_ref(s2))
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
