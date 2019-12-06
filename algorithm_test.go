package iter_test

import (
	"container/heap"
	"container/list"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	. "github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

const randN = 100

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randInt() int {
	return r.Intn(randN)
}

func randIntSlice() []int {
	lh := []int{randInt(), randInt()}
	Sort(begin(lh), end(lh))
	var s []int
	GenerateN(SliceBackInserter(&s), randInt(), func() Any { return lh[0] + r.Intn(lh[1]-lh[0]+1) })
	return s
}

func randString() string {
	alphabets := make([]byte, 26)
	Iota(begin(alphabets), end(alphabets), byte('a'))
	var bs StringBuilderInserter
	GenerateN(&bs, randInt(), RandomGenerator(alphabets, r))
	return bs.String()
}

var (
	begin    = SliceBegin
	end      = SliceEnd
	strBegin = StringBegin
	strEnd   = StringEnd
)

func sliceEqual(assert *assert.Assertions, a, b []int) {
	if len(a) == 0 && len(b) == 0 {
		return
	}
	assert.Equal(a, b)
}

func _eq(x, y Any) bool {
	type ieq interface{ Eq(Any) bool }
	if e, ok := x.(ieq); ok {
		return e.Eq(y)
	}
	return x == y
}

func TestAllAnyNoneOf(t *testing.T) {
	assert := assert.New(t)
	pred := func(x Any) bool { return x.(int)%2 == 0 }
	allOf := func(x []int) bool {
		for _, v := range x {
			if !pred(v) {
				return false
			}
		}
		return true
	}
	anyOf := func(x []int) bool {
		for _, v := range x {
			if pred(v) {
				return true
			}
		}
		return false
	}
	noneOf := func(x []int) bool {
		for _, v := range x {
			if pred(v) {
				return false
			}
		}
		return true
	}
	s := randIntSlice()
	assert.Equal(AllOf(begin(s), end(s), pred), allOf(s))
	assert.Equal(AnyOf(begin(s), end(s), pred), anyOf(s))
	assert.Equal(NoneOf(begin(s), end(s), pred), noneOf(s))
}

func TestForEach(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	var b []int
	f := func(x Any) { b = append(b, x.(int)) }
	ForEach(begin(a), end(a), f)
	sliceEqual(assert, a, b)
	n := r.Intn(len(a) + 1)
	b = nil
	ForEachN(begin(a), n, f)
	sliceEqual(assert, a[:n], b)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	count := make([]int, randN)
	ForEach(begin(a), end(a), func(x Any) { count[x.(int)]++ })
	for i := 0; i < 100; i++ {
		assert.Equal(Count(begin(a), end(a), i), count[i])
	}
}

func TestMismatch(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	var last2 ForwardReader
	if len(b) <= len(a) || r.Intn(2) == 0 {
		last2 = end(b)
	}
	it1, it2 := Mismatch(begin(a), end(a), begin(b), last2)
	n1, n2 := Distance(begin(a), it1), Distance(begin(b), it2)
	assert.Equal(n1, n2)
	sliceEqual(assert, a[:n1], b[:n1])
	assert.True((n1 >= len(a)) || (n1 >= len(b)) || a[n1] != b[n1])
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	f := func(x Any) bool { return x.(int)%2 == 0 }
	v := randInt()
	it := Find(begin(a), end(a), v)
	assert.True(NoneOf(begin(a), it, func(x Any) bool { return x.(int) == v }))
	if n := Distance(begin(a), it); n < len(a) {
		assert.Equal(a[n], v)
	}
	it = FindIf(begin(a), end(a), f)
	assert.True(NoneOf(begin(a), it, f))
	if n := Distance(begin(a), it); n < len(a) {
		assert.True(f(a[n]))
	}
	it = FindIfNot(begin(a), end(a), f)
	assert.True(AllOf(begin(a), it, f))
	if n := Distance(begin(a), it); n < len(a) {
		assert.False(f(a[n]))
	}
}

func TestFindEnd(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	it := FindEnd(begin(a), end(a), begin(b), end(b))
	if _eq(it, end(a)) {
		if len(b) > 0 {
			it = Search(begin(a), end(a), begin(b), end(b))
		}
	} else {
		assert.True(Equal(begin(b), end(b), it, nil))
		it = FindEnd(NextForwardReader(it), end(a), begin(b), end(b))
	}
	assert.True(_eq(it, end(a)))
}

func TestFindFirstOf(t *testing.T) {
	a, b := randString(), randString()
	i := strings.IndexAny(a, b)
	if i == -1 {
		i = len(a)
	}
	assert.New(t).True(_eq(
		FindFirstOf(strBegin(a), strEnd(a), strBegin(b), strEnd(b)),
		AdvanceN(strBegin(a), i),
	))
}

func TestAdjacentFind(t *testing.T) {
	a := randIntSlice()
	res := len(a)
	for i := 0; i < len(a)-1; i++ {
		if a[i] == a[i+1] {
			res = i
			break
		}
	}
	assert.New(t).True(_eq(
		AdjacentFind(begin(a), end(a)),
		AdvanceN(begin(a), res),
	))
}

func TestSearch(t *testing.T) {
	a, b := randString(), randString()
	i := strings.Index(a, b)
	if i == -1 {
		i = len(a)
	}
	assert.New(t).True(_eq(
		Search(strBegin(a), strEnd(a), strBegin(b), strEnd(b)),
		AdvanceN(strBegin(a), i),
	))
}

func TestSearchN(t *testing.T) {
	a := randString()
	c := byte('a' + byte(r.Intn(26)))
	n := r.Intn(10)
	b := strings.Repeat(string(c), n)
	i := strings.Index(a, b)
	if i == -1 {
		i = len(a)
	}
	assert.New(t).True(_eq(
		SearchN(strBegin(a), strEnd(a), n, c),
		AdvanceN(strBegin(a), i),
	))
}

func TestCopy(t *testing.T) {
	a := randIntSlice()
	var b []int
	Copy(begin(a), end(a), SliceBackInserter(&b))
	sliceEqual(assert.New(t), a, b)
}

func TestCopyIf(t *testing.T) {
	a := randIntSlice()
	var b []int
	f := func(x Any) bool { return x.(int)%2 == 0 }
	var c []int
	for _, x := range a {
		if f(x) {
			c = append(c, x)
		}
	}
	CopyIf(begin(a), end(a), SliceBackInserter(&b), f)
	sliceEqual(assert.New(t), b, c)
}

func TestCopyN(t *testing.T) {
	a := randIntSlice()
	n := r.Intn(len(a) + 1)
	var b []int
	CopyN(begin(a), n, SliceBackInserter(&b))
	sliceEqual(assert.New(t), b, a[:n])
}

func TestCopyBackward(t *testing.T) {
	a := randIntSlice()
	n := randInt()
	b := make([]int, len(a)+n)
	CopyBackward(begin(a), end(a), end(b))
	sliceEqual(assert.New(t), a, b[n:])
}

func TestFill(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	x := randInt()
	Fill(begin(a), end(a), x)
	assert.True(AllOf(begin(a), end(a), func(v Any) bool { return v.(int) == x }))
}

func TestFillN(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	for len(a) == 0 {
		a = randIntSlice()
	}
	b := append(a[:0:0], a...)
	n := r.Intn(len(a)) - r.Intn(len(a))
	x := randInt()
	FillN(begin(a), n, x)
	for i, v := range a {
		if i < n {
			assert.Equal(v, x)
		} else {
			assert.Equal(v, b[i])
		}
	}
}

func TestTransform(t *testing.T) {
	a := randIntSlice()
	var b []int
	for _, x := range a {
		b = append(b, x*2)
	}
	Transform(begin(a), end(a), begin(a), func(x Any) Any { return x.(int) * 2 })
	sliceEqual(assert.New(t), a, b)
}

func TestTransformBinary(t *testing.T) {
	a, b := randIntSlice(), randIntSlice()
	if len(a) > len(b) {
		a, b = b, a
	}
	c := make([]int, len(a))
	TransformBinary(begin(a), end(a), begin(b), begin(c), func(x, y Any) Any { return x.(int) * y.(int) })
	for i := range a {
		a[i] *= b[i]
	}
	sliceEqual(assert.New(t), a, c)
}

func TestGenerate(t *testing.T) {
	assert := assert.New(t)
	var i int
	g := func() Any { i++; return i }
	a := randIntSlice()
	Generate(begin(a), end(a), g)
	for i := range a {
		assert.Equal(i+1, a[i])
	}
}

func TestGenerateN(t *testing.T) {
	assert := assert.New(t)
	var i int
	g := func() Any { i++; return i }
	a := randIntSlice()
	b := append(a[:0:0], a...)
	n := r.Intn(len(a) + 1)
	GenerateN(begin(a), n, g)
	for i := range a {
		if i < n {
			assert.Equal(i+1, a[i])
		} else {
			assert.Equal(a[i], b[i])
		}
	}
}

func TestRemove(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	b := append(a[:0:0], a...)
	c := append(a[:0:0], a...)
	var d, e []int
	f := func(x Any) bool { return x.(int)%2 == 0 }

	count1 := Count(begin(a), end(a), 1)
	countf := CountIf(begin(a), end(a), f)
	Erase(&a, Remove(begin(a), end(a), 1))
	Erase(&b, RemoveIf(begin(b), end(b), f))
	RemoveCopy(begin(c), end(c), SliceBackInserter(&d), 1)
	RemoveCopyIf(begin(c), end(c), SliceBackInserter(&e), f)

	assert.Equal(Count(begin(a), end(a), 1), 0)
	assert.True(NoneOf(begin(b), end(b), f))
	assert.Equal(Count(begin(d), end(d), 1), 0)
	assert.True(NoneOf(begin(e), end(e), f))
	assert.Equal(len(a), len(c)-count1)
	assert.Equal(len(b), len(c)-countf)
	assert.Equal(len(d), len(c)-count1)
	assert.Equal(len(e), len(c)-countf)
}

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	b := append(a[:0:0], a...)
	c := append(a[:0:0], a...)
	var d, e []int
	f := func(x Any) bool { return x.(int)%2 == 0 }

	Replace(begin(a), end(a), 1, 2)
	ReplaceIf(begin(b), end(b), f, 1)
	ReplaceCopy(begin(c), end(c), SliceBackInserter(&d), 1, 2)
	ReplaceCopyIf(begin(c), end(c), SliceBackInserter(&e), f, 1)

	for i := range c {
		if c[i] == 1 {
			assert.Equal(a[i], 2)
			assert.Equal(d[i], 2)
		} else {
			assert.Equal(a[i], c[i])
			assert.Equal(d[i], c[i])
		}
		if f(c[i]) {
			assert.Equal(b[i], 1)
			assert.Equal(e[i], 1)
		} else {
			assert.Equal(b[i], c[i])
			assert.Equal(e[i], c[i])
		}
	}
}

func TestSwapRanges(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	a0 := append(a[:0:0], a...)
	b0 := append(b[:0:0], b...)
	l := Min(len(a), len(b)).(int)
	l = r.Intn(l + 1)
	s1 := r.Intn(len(a) - l + 1)
	s2 := r.Intn(len(b) - l + 1)
	SwapRanges(AdvanceNReadWriter(begin(a), s1), AdvanceNReadWriter(begin(a), s1+l), AdvanceNReadWriter(begin(b), s2))
	for i := range a {
		if i < s1 || i > s1+l {
			assert.Equal(a[i], a0[i])
		}
	}
	for i := range b {
		if i < s2 || i > s2+l {
			assert.Equal(b[i], b0[i])
		}
	}
	for i := 0; i < l; i++ {
		assert.Equal(a[s1+i], b0[s2+i])
		assert.Equal(b[s2+i], a0[s1+i])
	}
}

func TestReverse(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	b := append(a[:0:0], a...)
	c := append(a[:0:0], a...)
	Reverse(begin(a), end(a))
	ReverseCopy(begin(b), end(b), begin(c))
	for i := range a {
		assert.Equal(c[i], a[i])
		assert.Equal(b[len(b)-i-1], a[i])
	}
}

func TestRotate(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	b := append(a[:0:0], a...)
	n := r.Intn(len(a) + 1)
	d := append(a[:0:0], a[n:]...)
	d = append(d, a[:n]...)
	var c []int
	Rotate(begin(a), AdvanceNReadWriter(begin(a), n), end(a))
	RotateCopy(begin(b), AdvanceNReadWriter(begin(b), n), end(b), SliceBackInserter(&c))
	sliceEqual(assert, d, a)
	sliceEqual(assert, d, c)
}

func TestShuffle(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	count := make([]int, randN)
	for _, x := range a {
		count[x]++
	}
	Shuffle(begin(a), end(a), r)
	for _, x := range a {
		count[x]--
	}
	for _, x := range count {
		assert.Equal(x, 0)
	}
}

func TestSampleSelection(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	n := randInt()
	var b []int
	Sample(begin(a), end(a), SliceBackInserter(&b), n, r)
	count := make([]int, randN)
	for _, x := range a {
		count[x]++
	}
	assert.Equal(len(b), Min(n, len(a)).(int))
	for i := 0; i < len(b) && i < len(a); i++ {
		count[b[i]]--
		assert.GreaterOrEqual(count[b[i]], 0)
	}
}

func TestSampleReservoir(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	n := randInt()
	b := list.New()
	Copy(begin(a), end(a), ListBackInserter(b))
	c := make([]int, n)
	Sample(ListBegin(b), ListEnd(b), begin(c), n, r)
	count := make([]int, randN)
	for _, x := range a {
		count[x]++
	}
	for i := 0; i < len(c) && i < len(a); i++ {
		count[c[i]]--
		assert.GreaterOrEqual(count[c[i]], 0)
	}
}

func TestUnique(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	b := append(a[:0:0], a...)
	c := make([]int, len(a))
	Erase(&b, Unique(begin(b), end(b)))
	Erase(&c, UniqueCopy(begin(a), end(a), begin(c)))
	sliceEqual(assert, b, c)
	for i := 0; i < len(b)-1; i++ {
		assert.NotEqual(b[i], b[i+1])
	}
}

func TestPartition(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]bool, l)
	GenerateN(begin(a), l, func() Any { return r.Intn(2) == 0 })
	f := func(x Any) bool { return x.(bool) }
	checkPartition := func(a []bool) bool {
		var i int
		for ; i < len(a) && a[i]; i++ {
		}
		if i == len(a) {
			assert.True(IsPartitioned(begin(a), end(a), f))
			assert.Equal(i, Distance(begin(a), PartitionPoint(begin(a), end(a), f)))
			return true
		}
		m := i
		for i++; i < len(a); i++ {
			if a[i] {
				assert.False(IsPartitioned(begin(a), end(a), f))
				return false
			}
		}
		assert.True(IsPartitioned(begin(a), end(a), f))
		assert.Equal(m, Distance(begin(a), PartitionPoint(begin(a), end(a), f)))
		return true
	}
	checkPartition(a)

	var b, c []bool
	PartitionCopy(begin(a), end(a), SliceBackInserter(&b), SliceBackInserter(&c), f)
	ita := Partition(begin(a), end(a), f)
	assert.True(checkPartition(a))
	assert.True(AllOf(begin(b), end(b), f))
	assert.True(NoneOf(begin(c), end(c), f))
	assert.Equal(len(b), Distance(begin(a), ita))
}

type compareItem struct {
	a, b int
}

func (ci *compareItem) Equal(x Any) bool {
	return ci.a == x.(*compareItem).a
}

func (ci *compareItem) Less(x Any) bool {
	return ci.a < x.(*compareItem).a
}

func (ci *compareItem) Less2(x Any) bool {
	return ci.Less(x) ||
		(ci.a == x.(*compareItem).a && ci.b < x.(*compareItem).b)
}

func (ci *compareItem) String() string {
	return fmt.Sprintf("{a=%v,b=%v}", ci.a, ci.b)
}

type forwardListIter struct {
	l *list.List
	e *list.Element
}

func forwardListBegin(l *list.List) *forwardListIter {
	return &forwardListIter{
		l: l,
		e: l.Front(),
	}
}

func forwardListEnd(l *list.List) *forwardListIter {
	return &forwardListIter{
		l: l,
		e: l.Back(),
	}
}

func (l *forwardListIter) Eq(x Any) bool {
	return l.e == x.(*forwardListIter).e
}

func (l *forwardListIter) AllowMultiplePass() {}

func (l *forwardListIter) Next() Incrementable {
	return &forwardListIter{
		l: l.l,
		e: l.e.Next(),
	}
}

func (l *forwardListIter) Read() Any {
	return l.e.Value
}

func (l *forwardListIter) Write(x Any) {
	l.e.Value = x
}

func TestStablePartition(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]*compareItem, l)
	var id int
	GenerateN(begin(a), l, func() Any {
		id++
		return &compareItem{
			a: r.Intn(2),
			b: id,
		}
	})
	f := func(x Any) bool { return x.(*compareItem).a > 0 }
	b := list.New()
	Copy(begin(a), end(a), ListBackInserter(b))

	{
		StablePartition(begin(a), end(a), f)
		var i int
		for mb := 0; i < len(a) && f(a[i]); i++ {
			cb := a[i].b
			assert.Greater(cb, mb)
			mb = cb
		}
		for mb := 0; i < len(a); i++ {
			assert.False(f(a[i]))
			cb := a[i].b
			assert.Greater(cb, mb)
			mb = cb
		}
	}

	{
		StablePartition(forwardListBegin(b), forwardListEnd(b), f)
		var ele *list.Element
		for mb := 0; ele != nil && f(ele.Value); ele = ele.Next() {
			cb := ele.Value.(*compareItem).b
			assert.Greater(cb, mb)
			mb = cb
		}
		for mb := 0; ele != nil; ele = ele.Next() {
			assert.False(f(ele.Value))
			cb := ele.Value.(*compareItem).b
			assert.Greater(cb, mb)
			mb = cb
		}
	}
}

func TestSort(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	is := sort.IntSlice(a)
	assert.Equal(IsSorted(begin(a), end(a)), sort.IsSorted(is))
	it := IsSortedUntil(begin(a), end(a))
	if it.Eq(end(a)) {
		assert.True(sort.IsSorted(is))
	} else {
		n := Distance(begin(a), it)
		assert.True(sort.IsSorted(is[:n]))
		assert.False(sort.IsSorted(is[:n+1]))
	}

	Sort(begin(a), end(a))
	assert.True(sort.IsSorted(is))

	if len(a) == 0 {
		return
	}
	n := r.Intn(len(a)) + 1
	nth := AdvanceNReadWriter(begin(a), n-1)
	nv := nth.Read().(int)
	nth1 := NextRandomReadWriter(nth)

	Shuffle(begin(a), end(a), r)
	b := make([]int, n)
	PartialSortCopy(begin(a), end(a), begin(b), begin(b))
	sliceEqual(assert, b, make([]int, n))
	PartialSortCopy(begin(a), end(a), begin(b), end(b))
	assert.True(sort.IsSorted(sort.IntSlice(b)))

	Shuffle(begin(a), end(a), r)
	PartialSort(begin(a), nth1, end(a))
	sliceEqual(assert, a[:n], b)
	assert.GreaterOrEqual(MinElement(begin(a[n-1:]), end(a[n-1:])).Read().(int), a[n-1])

	Shuffle(begin(a), end(a), r)
	b = append(b[:0:0], a...)
	NthElement(begin(a), end(a), end(a))
	sliceEqual(assert, a, b)
	NthElement(begin(a), nth, end(a))
	assert.Equal(nth.Read().(int), nv)
}

func TestNthElement(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	cases := []string{
		"",
		"a",
		"aaaaaaaaa",
		"aaaaaaaaX",
		"baaaaaaaX",
		"baaaXaaaa",
		"abcamaaaz",
		"abcdefghi",
		"acccmxccz",
		"mmmmmmmcz",
		"aaabaaaaa",
		"aaaaaaaba",
		"aaabaaabaaa",
		"aaabaaabbaa",
	}
	for _, c := range cases {
		a := make([]int, 0, len(c))
		for _, x := range c {
			a = append(a, int(x))
		}
		s := append(a[:0:0], a...)
		sort.Ints(s)
		for i := 0; i <= len(a); i++ {
			b := append(a[:0:0], a...)
			NthElement(begin(b), AdvanceNReadWriter(begin(b), i), end(b))
			if i < len(a) {
				assert.Equal(b[i], s[i])
			}
		}
	}
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Sort(begin(a), end(a))
	Sort(begin(b), end(b))
	ab := append(a[:len(a):len(a)], b...)

	c := make([]int, len(a)+len(b))
	PartialSortCopy(begin(ab), end(ab), begin(c), end(c))
	var d []int
	Merge(begin(a), end(a), begin(b), end(b), SliceBackInserter(&d))
	sliceEqual(assert, c, d)

	middle := AdvanceN(begin(ab), len(a)).(BidiReadWriter)
	InplaceMerge(begin(ab), middle, end(ab))
	sliceEqual(assert, c, ab)
}

func TestSet(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	countA := make([]int, randN)
	ForEach(begin(a), end(a), func(x Any) { countA[x.(int)]++ })
	countB := make([]int, randN)
	ForEach(begin(b), end(b), func(x Any) { countB[x.(int)]++ })
	Sort(begin(a), end(a))
	Sort(begin(b), end(b))
	assert.Equal(
		Includes(begin(a), end(a), begin(b), end(b)),
		InnerProductBy(begin(countA), end(countA), begin(countB),
			true,
			func(acc, cur Any) Any { return acc.(bool) && cur.(bool) },
			func(a, b Any) Any { return a.(int) >= b.(int) }),
	)
	var diff, intersection, symmetric, union []int
	SetDifference(begin(a), end(a), begin(b), end(b), SliceBackInserter(&diff))
	SetIntersection(begin(a), end(a), begin(b), end(b), SliceBackInserter(&intersection))
	SetSymmetricDifference(begin(a), end(a), begin(b), end(b), SliceBackInserter(&symmetric))
	SetUnion(begin(a), end(a), begin(b), end(b), SliceBackInserter(&union))

	var diff2, intersection2, symmetric2, union2 []int
	for i := range countA {
		FillN(SliceBackInserter(&diff2), countA[i]-countB[i], i)
		FillN(SliceBackInserter(&intersection2), Min(countA[i], countB[i]).(int), i)
		FillN(SliceBackInserter(&symmetric2), Max(countA[i]-countB[i], countB[i]-countA[i]).(int), i)
		FillN(SliceBackInserter(&union2), Max(countA[i], countB[i]).(int), i)
	}

	sliceEqual(assert, diff, diff2)
	sliceEqual(assert, intersection, intersection2)
	sliceEqual(assert, symmetric, symmetric2)
	sliceEqual(assert, union, union2)
}

func TestBinarySearch(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	Sort(begin(a), end(a))
	x := randInt()
	l, h := LowerBound(begin(a), end(a), x), UpperBound(begin(a), end(a), x)
	l2, h2 := EqualRange(begin(a), end(a), x)
	assert.True(_eq(l, l2))
	assert.True(_eq(h, h2))
	ok := BinarySearch(begin(a), end(a), x)
	assert.Equal(ok, !Find(begin(a), end(a), x).Eq(end(a)))
	if l.Eq(end(a)) {
		assert.True(_eq(h, end(a)))
		if len(a) > 0 {
			assert.Less(a[len(a)-1], x)
		}
	} else {
		assert.GreaterOrEqual(l.Read(), x)
		if !l.Eq(begin(a)) {
			assert.Less(PrevBidiReader(l.(BidiReader)).Read(), x)
		}
		if !h.Eq(end(a)) {
			assert.Greater(h.Read(), x)
		}
		if !h.Eq(begin(a)) {
			assert.LessOrEqual(PrevBidiReader(h.(BidiReader)).Read(), x)
		}
	}
}

func TestStableSort(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]*compareItem, l)
	var id int
	GenerateN(begin(a), l, func() Any {
		id++
		return &compareItem{
			a: randInt(),
			b: id,
		}
	})
	StableSort(begin(a), end(a))
	assert.True(IsSortedBy(begin(a), end(a), func(x, y Any) bool { return x.(*compareItem).Less2(y) }))
}

func TestHeap(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	isHeap := func(a []int) bool {
		for i := 0; i < len(a); i++ {
			if 2*i+1 < len(a) && a[i] < a[2*i+1] {
				return false
			}
			if 2*i+2 < len(a) && a[i] < a[2*i+2] {
				return false
			}
		}
		return true
	}
	assert.Equal(IsHeap(begin(a), end(a)), isHeap(a))
	it := IsHeapUntil(begin(a), end(a))
	if it.Eq(end(a)) {
		assert.True(isHeap(a))
	} else {
		n := Distance(begin(a), it)
		assert.True(isHeap(a[:n]))
		assert.False(isHeap(a[:n+1]))
	}
	MakeHeap(begin(a), end(a))
	assert.True(isHeap(a))
	SortHeap(begin(a), end(a))
	assert.True(IsSorted(begin(a), end(a)))
}

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[j] < h[i] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TestHeapPP(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	aa := append(a[:0:0], a...)
	b := (*intHeap)(&aa)
	heap.Init(b)
	MakeHeap(begin(a), end(a))
	sliceEqual(assert, a, *b)
	n := randInt()
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 || len(a) == 0 {
			x := randInt()
			a = append(a, x)
			PushHeap(begin(a), end(a))
			heap.Push(b, x)
		} else {
			PopHeap(begin(a), end(a))
			a = a[:len(a)-1]
			heap.Pop(b)
		}
		sliceEqual(assert, a, *b)
	}
}

func TestMinmax(t *testing.T) {
	assert := assert.New(t)
	a, b := randInt(), randInt()
	min, max := Minmax(a, b)
	assert.LessOrEqual(min, max)
	assert.Equal(Max(a, b), max)
	assert.Equal(Min(a, b), min)
}

func TestMinmaxElement(t *testing.T) {
	assert := assert.New(t)
	s := randIntSlice()
	min, max := MinmaxElement(begin(s), end(s))
	min2, max2 := MinElement(begin(s), end(s)), MaxElement(begin(s), end(s))
	assert.True(NoneOf(begin(s), end(s), func(v Any) bool { return v.(int) > max.Read().(int) || v.(int) < min.Read().(int) }))
	if len(s) > 0 {
		assert.Equal(min.Read(), min2.Read())
		assert.Equal(max.Read(), max2.Read())
	} else {
		assert.True(_eq(min, end(s)))
		assert.True(_eq(max, end(s)))
		assert.True(_eq(min2, end(s)))
		assert.True(_eq(max2, end(s)))
	}
}

func TestClamp(t *testing.T) {
	assert := assert.New(t)
	l, h := Minmax(randInt(), randInt())
	v := randInt()
	c := Clamp(v, l, h)
	if c != v {
		assert.True(v < l.(int) || v > h.(int))
	}
	assert.GreaterOrEqual(c, l)
	assert.LessOrEqual(c, h)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)
	a, b := randString(), randString()
	if len(a) > len(b) {
		a, b = b, a
	}
	assert.Equal(Equal(strBegin(a), strEnd(a), strBegin(b), nil), a == b[:len(a)])
	a, b = randString(), randString()
	assert.Equal(Equal(strBegin(a), strEnd(a), strBegin(b), strEnd(b)), a == b)
}

func TestCompare(t *testing.T) {
	assert := assert.New(t)
	a, b := randString(), randString()
	if randInt() == 0 {
		b = a
	}
	x, y, z, w := strBegin(a), strEnd(a), strBegin(b), strEnd(b)
	if a == b {
		assert.True(Equal(x, y, z, w))
		assert.False(LexicographicalCompare(x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay(x, y, z, w), 0)
	} else if a < b {
		assert.False(Equal(x, y, z, w))
		assert.True(LexicographicalCompare(x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay(x, y, z, w), -1)
	} else {
		assert.False(Equal(x, y, z, w))
		assert.False(LexicographicalCompare(x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay(x, y, z, w), 1)
	}
}

func TestIsPermutation(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Generate(begin(b), end(b), func() Any {
		if len(a) > 0 {
			return a[r.Intn(len(a))]
		}
		return 0
	})

	count := make([]int, randN)
	for _, x := range a {
		count[x]++
	}
	for _, x := range b {
		count[x]--
	}
	assert.Equal(
		IsPermutation(begin(a), end(a), begin(b), end(b)),
		Count(begin(count), end(count), 0) == randN,
	)

}

func TestIsPermutation2(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Generate(begin(b), end(b), func() Any {
		if len(a) > 0 {
			return a[r.Intn(len(a))]
		}
		return 0
	})
	if len(a) > len(b) {
		a, b = b, a
	}
	count := make([]int, randN)
	for i := range a {
		count[a[i]]++
		count[b[i]]--
	}
	assert.Equal(
		IsPermutation(begin(a), end(a), begin(b), nil),
		Count(begin(count), end(count), 0) == randN,
	)
}

func TestPermutation(t *testing.T) {
	assert := assert.New(t)
	total := []int{0: 0, 1: 1, 2: 2, 3: 6, 4: 24, 5: 120}
	a := randIntSlice()
	ml := Min(len(a), 5).(int)
	a = a[:r.Intn(ml+1)]
	b := append(a[:0:0], a...)
	c := make([]int, len(a))
	for i := 0; ; i++ {
		Copy(begin(a), end(a), begin(c))
		ok := NextPermutation(begin(a), end(a))
		assert.Equal(LexicographicalCompare(begin(c), end(c), begin(a), end(a)), ok)
		if Equal(begin(a), end(a), begin(b), end(b)) {
			break
		}
		assert.Less(i, total[len(a)])
	}
	for i := 0; ; i++ {
		Copy(begin(a), end(a), begin(c))
		ok := PrevPermutation(begin(a), end(a))
		assert.Equal(LexicographicalCompare(begin(a), end(a), begin(c), end(c)), ok)
		if Equal(begin(a), end(a), begin(b), end(b)) {
			break
		}
		assert.Less(i, total[len(a)])
	}
}

func TestIota(t *testing.T) {
	l := randInt()
	a := make([]int, l)
	b := make([]int, l)
	s := randInt()
	Iota(begin(a), end(a), s+1)
	Generate(begin(b), end(b), func() Any { s++; return s })
	sliceEqual(assert.New(t), a, b)
}

func TestAccumulate(t *testing.T) {
	a := randIntSlice()
	sum := Accumulate(begin(a), end(a), 0)
	sum2 := 0
	ForEach(begin(a), end(a), func(it Any) {
		sum2 += it.(int)
	})
	assert.New(t).Equal(sum, sum2)
}

func TestInnerProduct(t *testing.T) {
	a, b := randIntSlice(), randIntSlice()
	l := Min(len(a), len(b)).(int)
	p := InnerProduct(begin(a), end(a[:l]), begin(b), 0)
	var p2 int
	for i := 0; i < l; i++ {
		p2 += a[i] * b[i]
	}
	assert.New(t).Equal(p, p2)
}

func TestPartialSum(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	diff := make([]int, len(a))
	ps := make([]int, len(a))
	exc := make([]int, len(a))
	inc := make([]int, len(a))
	exct := make([]int, len(a))
	inct := make([]int, len(a))
	for i := range a {
		if i == 0 {
			diff[i] = a[i]
			ps[i] = a[i]
			exc[i] = 1
			inc[i] = 2 + a[i]
			exct[i] = 3
			inct[i] = 4 + a[i]*a[i]
		} else {
			diff[i] = a[i] - a[i-1]
			ps[i] = ps[i-1] + a[i]
			exc[i] = exc[i-1] + a[i-1]
			exct[i] = exct[i-1] + a[i-1]*2
			inc[i] = inc[i-1] + a[i]
			inct[i] = inct[i-1] + a[i]*a[i]
		}
	}
	g := make([]int, len(a))
	AdjacentDifference(begin(a), end(a), begin(g))
	sliceEqual(assert, g, diff)
	PartialSum(begin(a), end(a), begin(g))
	sliceEqual(assert, g, ps)
	ExclusiveScan(begin(a), end(a), begin(g), 1)
	sliceEqual(assert, g, exc)
	InclusiveScan(begin(a), end(a), begin(g), 2)
	sliceEqual(assert, g, inc)
	TransformExclusiveScan(begin(a), end(a), begin(g), 3, func(x Any) Any { return x.(int) * 2 })
	sliceEqual(assert, g, exct)
	TransformInclusiveScan(begin(a), end(a), begin(g), 4, func(x Any) Any { return x.(int) * x.(int) })
	sliceEqual(assert, g, inct)
}

type dummyObj struct {
	s string
}

func (obj dummyObj) Eq(x Any) bool {
	return len(obj.s) == len(x.(dummyObj).s)
}

func (obj dummyObj) Less(x Any) bool {
	return len(obj.s) < len(x.(dummyObj).s)
}

func (obj dummyObj) Cmp(x Any) int {
	if obj.Eq(x) {
		return 0
	}
	if obj.Less(x) {
		return -1
	}
	return 1
}

func (obj dummyObj) Inc() Any {
	return dummyObj{s: obj.s + "a"}
}

func (obj dummyObj) Add(x Any) Any {
	return dummyObj{s: obj.s + x.(dummyObj).s}
}

func (obj dummyObj) Sub(x Any) Any {
	return dummyObj{s: obj.s[:len(obj.s)-len(x.(dummyObj).s)]}
}

func (obj dummyObj) Mul(x Any) Any {
	return dummyObj{s: strings.Repeat(obj.s, len(x.(dummyObj).s))}
}

func TestCustomType(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	a, b, c := []dummyObj{{"abc"}}, []dummyObj{{"xxxx"}}, []dummyObj{{"xyz"}}
	// _eq
	assert.False(Equal(begin(a), end(a), begin(b), nil))
	assert.True(Equal(begin(a), end(a), begin(c), nil))
	// _cmp
	assert.Equal(LexicographicalCompareThreeWay(begin(a), end(a), begin(b), end(b)), -1)
	assert.Equal(LexicographicalCompareThreeWay(begin(a), end(a), begin(c), end(c)), 0)
	// _inc
	d := make([]dummyObj, 2)
	Iota(begin(d), end(d), a[0])
	assert.True(Equal(begin(b), end(b), begin(d[1:]), nil))
	// _add
	assert.True(Accumulate(begin(d), end(d), dummyObj{""}).(dummyObj).Eq(dummyObj{"1234567"}))
	// _sub
	e := make([]dummyObj, 2)
	AdjacentDifference(begin(d), end(d), begin(e))
	assert.True(e[1].Eq(dummyObj{"b"}))
	// _mul
	assert.True(InnerProduct(begin(d), end(d), begin(e), dummyObj{}).(dummyObj).Eq(dummyObj{strings.Repeat("x", 13)}))
}
