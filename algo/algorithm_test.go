package algo_test

import (
	"container/heap"
	"container/list"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	. "github.com/disksing/iter/v2"
	. "github.com/disksing/iter/v2/algo"
	"github.com/disksing/iter/v2/lists"
	"github.com/disksing/iter/v2/slices"
	"github.com/disksing/iter/v2/strs"
	"github.com/stretchr/testify/assert"
)

var (
	fuzzTime = flag.Duration("fuzz-time", 0, "fuzz test timeout")
)

func TestMain(m *testing.M) {
	if code := m.Run(); code != 0 {
		os.Exit(code)
	}
	if *fuzzTime > 0 {
		start := time.Now()
		for time.Since(start) < *fuzzTime {
			if code := m.Run(); code != 0 {
				os.Exit(code)
			}
		}
	}
}

const randN = 20

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	_first_int   = slices.Begin[int]
	_first_int_r = slices.RBegin[int]
	_last_int    = slices.End[int]
	_last_int_r  = slices.REnd[int]

	_first_byte   = slices.Begin[byte]
	_first_byte_r = slices.RBegin[byte]
	_last_byte    = slices.End[byte]
	_last_byte_r  = slices.REnd[byte]

	_first_bool   = slices.Begin[bool]
	_first_bool_r = slices.RBegin[bool]
	_last_bool    = slices.End[bool]
	_last_bool_r  = slices.REnd[bool]

	_first_str   = strs.Begin
	_first_str_r = strs.RBegin
	_last_str    = strs.End
	_last_str_r  = strs.REnd

	_head_int   = lists.Begin[int]
	_head_int_r = lists.RBegin[int]
	_tail_int   = lists.End[int]
	_tail_int_r = lists.REnd[int]
)

func randInt() int {
	return r.Intn(randN)
}

func randIntSlice() []int {
	lh := []int{randInt(), randInt()}
	Sort[int](_first_int(lh), _last_int(lh))
	var s []int
	GenerateN(slices.Appender(&s), randInt(), func() int { return lh[0] + r.Intn(lh[1]-lh[0]+1) })
	return s
}

func randString() string {
	alphabets := make([]byte, 26)
	Iota(_first_byte(alphabets), _last_byte(alphabets), byte('a'))
	var bs strs.StringBuilderInserter[byte]
	GenerateN(&bs, randInt(), RandomGenerator(alphabets, r))
	return bs.String()
}

func sliceEqual(assert *assert.Assertions, a, b []int) {
	if len(a) == 0 && len(b) == 0 {
		return
	}
	assert.Equal(a, b)
}

func __eq[T Comparable[T]](x, y T) bool {
	return x.Eq(y)
}

func TestAllAnyNoneOf(t *testing.T) {
	assert := assert.New(t)
	pred := func(x int) bool { return x%2 == 0 }
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
	assert.Equal(AllOf(_first_int(s), _last_int(s), pred), allOf(s))
	assert.Equal(AnyOf(_first_int(s), _last_int(s), pred), anyOf(s))
	assert.Equal(NoneOf(_first_int(s), _last_int(s), pred), noneOf(s))
}

func TestForEach(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	var b []int
	f := func(x int) { b = append(b, x) }
	ForEach(_first_int(a), _last_int(a), f)
	sliceEqual(assert, a, b)
	n := r.Intn(len(a) + 1)
	b = nil
	ForEachN(_first_int(a), n, f)
	sliceEqual(assert, a[:n], b)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	count := make([]int, randN)
	ForEach(_first_int(a), _last_int(a), func(x int) { count[x]++ })
	for i := 0; i < randN; i++ {
		assert.Equal(Count(_first_int(a), _last_int(a), i), count[i])
	}
}

func TestMismatch(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	l2 := _last_int(b)
	last2 := &l2
	if len(b) > len(a) && r.Intn(2) == 0 {
		last2 = nil
	}
	it1, it2 := Mismatch[int](_first_int(a), _last_int(a), _first_int(b), last2)
	n1, n2 := Distance[int](_first_int(a), it1), Distance[int](_first_int(b), it2)
	assert.Equal(n1, n2)
	sliceEqual(assert, a[:n1], b[:n1])
	assert.True((n1 >= len(a)) || (n1 >= len(b)) || a[n1] != b[n1])
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	f := func(x int) bool { return x%2 == 0 }
	v := randInt()
	it := Find(_first_int(a), _last_int(a), v)
	assert.True(NoneOf(_first_int(a), it, func(x int) bool { return x == v }))
	if n := Distance[int](_first_int(a), it); n < len(a) {
		assert.Equal(a[n], v)
	}
	it = FindIf(_first_int(a), _last_int(a), f)
	assert.True(NoneOf(_first_int(a), it, f))
	if n := Distance[int](_first_int(a), it); n < len(a) {
		assert.True(f(a[n]))
	}
	it = FindIfNot(_first_int(a), _last_int(a), f)
	assert.True(AllOf(_first_int(a), it, f))
	if n := Distance[int](_first_int(a), it); n < len(a) {
		assert.False(f(a[n]))
	}
}

func TestFindEnd(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	it := FindEnd[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b))
	if __eq(it, _last_int(a)) {
		if len(b) > 0 {
			it = Search[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b))
		}
	} else {
		assert.True(Equal[int](_first_int(b), _last_int(b), it, nil))
		it = FindEnd[int](it.Next(), _last_int(a), _first_int(b), _last_int(b))
	}
	assert.True(__eq(it, _last_int(a)))
}

func TestFindFirstOf(t *testing.T) {
	a, b := randString(), randString()
	i := strings.IndexAny(a, b)
	if i == -1 {
		i = len(a)
	}
	assert.New(t).True(__eq(
		FindFirstOf[byte](_first_str(a), _last_str(a), _first_str(b), _last_str(b)),
		_first_str(a).AdvanceN(i),
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
	assert.New(t).True(__eq(
		AdjacentFind[int](_first_int(a), _last_int(a)),
		AdvanceN[int](_first_int(a), res),
	))
}

func TestSearch(t *testing.T) {
	a, b := randString(), randString()
	i := strings.Index(a, b)
	if i == -1 {
		i = len(a)
	}
	assert.New(t).True(__eq(
		Search[byte](_first_str(a), _last_str(a), _first_str(b), _last_str(b)),
		AdvanceN[byte](_first_str(a), i),
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
	assert.New(t).True(__eq(
		SearchN(_first_str(a), _last_str(a), n, c),
		AdvanceN[byte](_first_str(a), i),
	))
}

func TestCopy(t *testing.T) {
	a := randIntSlice()
	var b []int
	Copy[int](_first_int(a), _last_int(a), slices.Appender(&b))
	sliceEqual(assert.New(t), a, b)
}

func TestCopyIf(t *testing.T) {
	a := randIntSlice()
	var b []int
	f := func(x int) bool { return x%2 == 0 }
	var c []int
	for _, x := range a {
		if f(x) {
			c = append(c, x)
		}
	}
	CopyIf(_first_int(a), _last_int(a), slices.Appender(&b), f)
	sliceEqual(assert.New(t), b, c)
}

func TestCopyN(t *testing.T) {
	a := randIntSlice()
	n := r.Intn(len(a) + 1)
	var b []int
	CopyN[int](_first_int(a), n, slices.Appender(&b))
	sliceEqual(assert.New(t), b, a[:n])
}

func TestCopyBackward(t *testing.T) {
	a := randIntSlice()
	n := randInt()
	b := make([]int, len(a)+n)
	CopyBackward[int](_first_int(a), _last_int(a), _last_int(b))
	sliceEqual(assert.New(t), a, b[n:])
}

func TestFill(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	x := randInt()
	Fill(_first_int(a), _last_int(a), x)
	assert.True(AllOf(_first_int(a), _last_int(a), func(v int) bool { return v == x }))
}

func TestFillN(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	for len(a) == 0 {
		a = randIntSlice()
	}
	b := append(a[:0:0], a...)
	n := r.Intn(len(a)) - randInt()
	x := randInt()
	FillN(_first_int(a), n, x)
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
	Transform(_first_int(a), _last_int(a), _first_int(a), func(x int) int { return x * 2 })
	sliceEqual(assert.New(t), a, b)
}

func TestTransformBinary(t *testing.T) {
	a, b := randIntSlice(), randIntSlice()
	if len(a) > len(b) {
		a, b = b, a
	}
	c := make([]int, len(a))
	TransformBinary(_first_int(a), _last_int(a), _first_int(b), _first_int(c), func(x, y int) int { return x * y })
	for i := range a {
		a[i] *= b[i]
	}
	sliceEqual(assert.New(t), a, c)
}

func TestGenerate(t *testing.T) {
	assert := assert.New(t)
	var i int
	g := func() int { i++; return i }
	a := randIntSlice()
	Generate(_first_int(a), _last_int(a), g)
	for i := range a {
		assert.Equal(i+1, a[i])
	}
}

func TestGenerateN(t *testing.T) {
	assert := assert.New(t)
	var i int
	g := func() int { i++; return i }
	a := randIntSlice()
	b := append(a[:0:0], a...)
	n := r.Intn(len(a) + 1)
	GenerateN(_first_int(a), n, g)
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
	f := func(x int) bool { return x%2 == 0 }

	count1 := Count(_first_int(a), _last_int(a), 1)
	countf := CountIf(_first_int(a), _last_int(a), f)
	a = a[:Distance[int](_first_int(a), Remove(_first_int(a), _last_int(a), 1))]
	b = b[:Distance[int](_first_int(b), RemoveIf(_first_int(b), _last_int(b), f))]
	RemoveCopy(_first_int(c), _last_int(c), slices.Appender(&d), 1)
	RemoveCopyIf(_first_int(c), _last_int(c), slices.Appender(&e), f)

	assert.Equal(Count(_first_int(a), _last_int(a), 1), 0)
	assert.True(NoneOf(_first_int(b), _last_int(b), f))
	assert.Equal(Count(_first_int(d), _last_int(d), 1), 0)
	assert.True(NoneOf(_first_int(e), _last_int(e), f))
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
	f := func(x int) bool { return x%2 == 0 }

	Replace(_first_int(a), _last_int(a), 1, 2)
	ReplaceIf(_first_int(b), _last_int(b), f, 1)
	ReplaceCopy(_first_int(c), _last_int(c), slices.Appender(&d), 1, 2)
	ReplaceCopyIf(_first_int(c), _last_int(c), slices.Appender(&e), f, 1)

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
	l := Min(len(a), len(b))
	l = r.Intn(l + 1)
	s1 := r.Intn(len(a) - l + 1)
	s2 := r.Intn(len(b) - l + 1)
	SwapRanges[int](_first_int(a).AdvanceN(s1), _first_int(a).AdvanceN(s1+l), _first_int(b).AdvanceN(s2))
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
	Reverse[int](_first_int(a), _last_int(a))
	ReverseCopy[int](_first_int(b), _last_int(b), _first_int(c))
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
	Rotate[int](_first_int(a), _first_int(a).AdvanceN(n), _last_int(a))
	RotateCopy[int](_first_int(b), _first_int(b).AdvanceN(n), _last_int(b), slices.Appender(&c))
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
	Shuffle[int](_first_int(a), _last_int(a), r)
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
	Sample[int](_first_int(a), _last_int(a), slices.Appender(&b), n, r)
	count := make([]int, randN)
	for _, x := range a {
		count[x]++
	}
	assert.Equal(len(b), Min(n, len(a)))
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
	Copy[int](_first_int(a), _last_int(a), lists.ListBackInserter[int](b))
	c := make([]int, n)
	Sample[int](_head_int(b), _tail_int(b), _first_int(c), n, r)
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
	b = b[:Distance[int](_first_int(b), Unique[int](_first_int(b), _last_int(b)))]
	c = c[:Distance[int](_first_int(c), UniqueCopy[int](_first_int(a), _last_int(a), _first_int(c)))]
	sliceEqual(assert, b, c)
	for i := 0; i < len(b)-1; i++ {
		assert.NotEqual(b[i], b[i+1])
	}
}

func TestPartition(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]bool, l)
	GenerateN(_first_bool(a), l, func() bool { return r.Intn(2) == 0 })
	f := func(x bool) bool { return x }
	checkPartition := func(a []bool) bool {
		var i int
		for ; i < len(a) && a[i]; i++ {
		}
		if i == len(a) {
			assert.True(IsPartitioned(_first_bool(a), slices.End(a), f))
			assert.Equal(i, Distance[bool](_first_bool(a), PartitionPoint(_first_bool(a), slices.End(a), f)))
			return true
		}
		m := i
		for i++; i < len(a); i++ {
			if a[i] {
				assert.False(IsPartitioned(_first_bool(a), _last_bool(a), f))
				return false
			}
		}
		assert.True(IsPartitioned(_first_bool(a), _last_bool(a), f))
		assert.Equal(m, Distance[bool](_first_bool(a), PartitionPoint(_first_bool(a), _last_bool(a), f)))
		return true
	}
	checkPartition(a)

	var b, c []bool
	PartitionCopy(_first_bool(a), _last_bool(a), slices.Appender(&b), slices.Appender(&c), f)
	ita := Partition(_first_bool(a), _last_bool(a), f)
	assert.True(checkPartition(a))
	assert.True(AllOf(_first_bool(b), _last_bool(b), f))
	assert.True(NoneOf(_first_bool(c), _last_bool(c), f))
	assert.Equal(len(b), Distance[int](_first_bool(a), ita))
}

type compareItem struct {
	a, b int
}

func (ci *compareItem) Equal(x any) bool {
	return ci.a == x.(*compareItem).a
}

func (ci *compareItem) Less(x any) bool {
	return ci.a < x.(*compareItem).a
}

func (ci *compareItem) Less2(x any) bool {
	return ci.Less(x) ||
		(ci.a == x.(*compareItem).a && ci.b < x.(*compareItem).b)
}

func (ci *compareItem) String() string {
	return fmt.Sprintf("{a=%v,b=%v}", ci.a, ci.b)
}

type forwardListIter[T any] struct {
	l *list.List
	e *list.Element
}

func forwardListBegin[T any](l *list.List) *forwardListIter[T] {
	return &forwardListIter[T]{
		l: l,
		e: l.Front(),
	}
}

func forwardListEnd[T any](l *list.List) *forwardListIter[T] {
	return &forwardListIter[T]{
		l: l,
		e: l.Back(),
	}
}

func (l *forwardListIter[T]) Eq(x *forwardListIter[T]) bool {
	return l.e == x.e
}

func (l *forwardListIter[T]) AllowMultiplePass() {}

func (l *forwardListIter[T]) Next() *forwardListIter[T] {
	return &forwardListIter[T]{
		l: l.l,
		e: l.e.Next(),
	}
}

func (l *forwardListIter[T]) Read() T {
	return l.e.Value.(T)
}

func (l *forwardListIter[T]) Write(x T) {
	l.e.Value = x
}

func TestStablePartition(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]*compareItem, l)
	var id int
	GenerateN(slices.Begin(a), l, func() *compareItem {
		id++
		return &compareItem{
			a: r.Intn(2),
			b: id,
		}
	})
	f := func(x *compareItem) bool { return x.a > 0 }
	b := list.New()
	Copy[*compareItem](slices.Begin(a), slices.End(a), lists.ListBackInserter[*compareItem](b))

	{
		StablePartitionBidi(slices.Begin(a), slices.End(a), f)
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
		StablePartition(forwardListBegin[*compareItem](b), forwardListEnd[*compareItem](b), f)
		var ele *list.Element
		for mb := 0; ele != nil && f(ele.Value.(*compareItem)); ele = ele.Next() {
			cb := ele.Value.(*compareItem).b
			assert.Greater(cb, mb)
			mb = cb
		}
		for mb := 0; ele != nil; ele = ele.Next() {
			assert.False(f(ele.Value.(*compareItem)))
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
	assert.Equal(IsSorted[int](_first_int(a), _last_int(a)), sort.IsSorted(is))
	it := IsSortedUntil[int](_first_int(a), _last_int(a))
	if __eq(it, _last_int(a)) {
		assert.True(sort.IsSorted(is))
	} else {
		n := Distance[int](_first_int(a), it)
		assert.True(sort.IsSorted(is[:n]))
		assert.False(sort.IsSorted(is[:n+1]))
	}

	Sort[int](_first_int(a), _last_int(a))
	assert.True(sort.IsSorted(is))

	if len(a) == 0 {
		return
	}
	n := r.Intn(len(a)) + 1
	nth := _first_int(a).AdvanceN(n - 1)
	nv := nth.Read()
	nth1 := nth.Next()

	Shuffle[int](_first_int(a), _last_int(a), r)
	b := make([]int, n)
	PartialSortCopy[int](_first_int(a), _last_int(a), _first_int(b), _first_int(b))
	sliceEqual(assert, b, make([]int, n))
	PartialSortCopy[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b))
	assert.True(sort.IsSorted(sort.IntSlice(b)))

	Shuffle[int](_first_int(a), _last_int(a), r)
	PartialSort[int](_first_int(a), nth1, _last_int(a))
	sliceEqual(assert, a[:n], b)
	assert.GreaterOrEqual(MinElement[int](_first_int(a[n-1:]), _last_int(a[n-1:])).Read(), a[n-1])

	Shuffle[int](_first_int(a), _last_int(a), r)
	b = append(b[:0:0], a...)
	NthElement[int](_first_int(a), _last_int(a), _last_int(a))
	sliceEqual(assert, a, b)
	NthElement[int](_first_int(a), nth, _last_int(a))
	assert.Equal(nth.Read(), nv)
}

func TestNthElement(t *testing.T) {
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
			NthElement[int](_first_int(b), _first_int(b).AdvanceN(i), _last_int(b))
			if i < len(a) {
				assert.Equal(b[i], s[i])
			}
		}
	}
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Sort[int](_first_int(a), _last_int(a))
	Sort[int](_first_int(b), _last_int(b))
	ab := append(a[:len(a):len(a)], b...)

	c := make([]int, len(a)+len(b))
	PartialSortCopy[int](_first_int(ab), _last_int(ab), _first_int(c), _last_int(c))
	var d []int
	Merge[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b), slices.Appender(&d))
	sliceEqual(assert, c, d)

	middle := _first_int(ab).AdvanceN(len(a))
	InplaceMerge[int](_first_int(ab), middle, _last_int(ab))
	sliceEqual(assert, c, ab)
}

func TestSet(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	countA := make([]int, randN)
	ForEach(_first_int(a), _last_int(a), func(x int) { countA[x]++ })
	countB := make([]int, randN)
	ForEach(_first_int(b), _last_int(b), func(x int) { countB[x]++ })
	Sort[int](_first_int(a), _last_int(a))
	Sort[int](_first_int(b), _last_int(b))
	assert.Equal(
		Includes[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b)),
		InnerProductBy(_first_int(countA), _last_int(countA), _first_int(countB),
			true,
			func(acc, cur bool) bool { return acc && cur },
			func(a, b int) bool { return a >= b }),
	)
	var diff, intersection, symmetric, union []int
	SetDifference[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b), slices.Appender(&diff))
	SetIntersection[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b), slices.Appender(&intersection))
	SetSymmetricDifference[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b), slices.Appender(&symmetric))
	SetUnion[int](_first_int(a), _last_int(a), _first_int(b), _last_int(b), slices.Appender(&union))

	var diff2, intersection2, symmetric2, union2 []int
	for i := range countA {
		FillN(slices.Appender(&diff2), countA[i]-countB[i], i)
		FillN(slices.Appender(&intersection2), Min(countA[i], countB[i]), i)
		FillN(slices.Appender(&symmetric2), Max(countA[i]-countB[i], countB[i]-countA[i]), i)
		FillN(slices.Appender(&union2), Max(countA[i], countB[i]), i)
	}

	sliceEqual(assert, diff, diff2)
	sliceEqual(assert, intersection, intersection2)
	sliceEqual(assert, symmetric, symmetric2)
	sliceEqual(assert, union, union2)
}

func TestBinarySearch(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	Sort[int](_first_int(a), _last_int(a))
	x := randInt()
	l, h := LowerBound(_first_int(a), _last_int(a), x), UpperBound(_first_int(a), _last_int(a), x)
	l2, h2 := EqualRange(_first_int(a), _last_int(a), x)
	assert.True(__eq(l, l2))
	assert.True(__eq(h, h2))
	ok := BinarySearch(_first_int(a), _last_int(a), x)
	assert.Equal(ok, !Find(_first_int(a), _last_int(a), x).Eq(_last_int(a)))
	if l.Eq(_last_int(a)) {
		assert.True(__eq(h, _last_int(a)))
		if len(a) > 0 {
			assert.Less(a[len(a)-1], x)
		}
	} else {
		assert.GreaterOrEqual(l.Read(), x)
		if !l.Eq(_first_int(a)) {
			assert.Less(l.Prev().Read(), x)
		}
		if !h.Eq(_last_int(a)) {
			assert.Greater(h.Read(), x)
		}
		if !h.Eq(_first_int(a)) {
			assert.LessOrEqual(h.Prev().Read(), x)
		}
	}
}

func TestStableSort(t *testing.T) {
	assert := assert.New(t)
	l := randInt()
	a := make([]*compareItem, l)
	var id int
	GenerateN(slices.Begin(a), l, func() *compareItem {
		id++
		return &compareItem{
			a: randInt(),
			b: id,
		}
	})
	StableSortBy(slices.Begin(a), slices.End(a), func(x, y *compareItem) bool { return x.Less(y) })
	assert.True(IsSortedBy(slices.Begin(a), slices.End(a), func(x, y *compareItem) bool { return x.Less2(y) }))
}

func TestHeap(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	isMaxHeap := func(a []int) bool {
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
	assert.Equal(IsHeap[int](_first_int(a), _last_int(a)), isMaxHeap(a))
	it := IsHeapUntil[int](_first_int(a), _last_int(a))
	if __eq(it, _last_int(a)) {
		assert.True(isMaxHeap(a))
	} else {
		n := Distance[int](_first_int(a), it)
		assert.True(isMaxHeap(a[:n]))
		assert.False(isMaxHeap(a[:n+1]))
	}
	MakeHeap[int](_first_int(a), _last_int(a))
	assert.True(isMaxHeap(a))
	SortHeap[int](_first_int(a), _last_int(a))
	assert.True(IsSorted[int](_first_int(a), _last_int(a)))
}

type intHeap []int

func (h intHeap) Len() int           { return len(h) }
func (h intHeap) Less(i, j int) bool { return h[j] < h[i] }
func (h intHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() any {
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
	MakeHeap[int](_first_int(a), _last_int(a))
	sliceEqual(assert, a, *b)
	n := randInt()
	for i := 0; i < n; i++ {
		if r.Intn(2) == 0 || len(a) == 0 {
			x := randInt()
			a = append(a, x)
			PushHeap[int](_first_int(a), _last_int(a))
			heap.Push(b, x)
		} else {
			PopHeap[int](_first_int(a), _last_int(a))
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
	min, max := MinmaxElement[int](_first_int(s), _last_int(s))
	min2, max2 := MinElement[int](_first_int(s), _last_int(s)), MaxElement[int](_first_int(s), _last_int(s))
	assert.True(NoneOf(_first_int(s), _last_int(s), func(v int) bool { return v > max.Read() || v < min.Read() }))
	if len(s) > 0 {
		assert.Equal(min.Read(), min2.Read())
		assert.Equal(max.Read(), max2.Read())
	} else {
		assert.True(__eq(min, _last_int(s)))
		assert.True(__eq(max, _last_int(s)))
		assert.True(__eq(min2, _last_int(s)))
		assert.True(__eq(max2, _last_int(s)))
	}
}

func TestClamp(t *testing.T) {
	assert := assert.New(t)
	l, h := Minmax(randInt(), randInt())
	v := randInt()
	c := Clamp(v, l, h)
	if c != v {
		assert.True(v < l || v > h)
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
	assert.Equal(Equal[byte](_first_str(a), _last_str(a), _first_str(b), nil), a == b[:len(a)])
	a, b = randString(), randString()
	lastb := _last_str(b)
	assert.Equal(Equal[byte](_first_str(a), _last_str(a), _first_str(b), &lastb), a == b)
}

func TestCompare(t *testing.T) {
	assert := assert.New(t)
	a, b := randString(), randString()
	if randInt() == 0 {
		b = a
	}
	x, y, z, w := _first_str(a), _last_str(a), _first_str(b), _last_str(b)
	if a == b {
		assert.True(Equal[byte](x, y, z, &w))
		assert.False(LexicographicalCompare[byte](x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay[byte](x, y, z, w), 0)
	} else if a < b {
		assert.False(Equal[byte](x, y, z, &w))
		assert.True(LexicographicalCompare[byte](x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay[byte](x, y, z, w), -1)
	} else {
		assert.False(Equal[byte](x, y, z, &w))
		assert.False(LexicographicalCompare[byte](x, y, z, w))
		assert.Equal(LexicographicalCompareThreeWay[byte](x, y, z, w), 1)
	}
}

func TestIsPermutation(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Generate(_first_int(b), _last_int(b), func() int {
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
	lastb := _last_int(b)
	assert.Equal(
		IsPermutation[int](_first_int(a), _last_int(a), _first_int(b), &lastb),
		Count(_first_int(count), _last_int(count), 0) == randN,
	)

}

func TestIsPermutation2(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	Generate(_first_int(b), _last_int(b), func() int {
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
		IsPermutation[int](_first_int(a), _last_int(a), _first_int(b), nil),
		Count(_first_int(count), _last_int(count), 0) == randN,
	)
}

func TestPermutation(t *testing.T) {
	assert := assert.New(t)
	total := []int{0: 0, 1: 1, 2: 2, 3: 6, 4: 24, 5: 120}
	a := randIntSlice()
	ml := Min(len(a), 5)
	a = a[:r.Intn(ml+1)]
	b := append(a[:0:0], a...)
	c := make([]int, len(a))
	for i := 0; ; i++ {
		Copy[int](_first_int(a), _last_int(a), _first_int(c))
		ok := NextPermutation[int](_first_int(a), _last_int(a))
		assert.Equal(LexicographicalCompare[int](_first_int(c), _last_int(c), _first_int(a), _last_int(a)), ok)
		last := _last_int(a)
		if Equal[int](_first_int(a), _last_int(a), _first_int(b), &last) {
			break
		}
		assert.Less(i, total[len(a)])
	}
	for i := 0; ; i++ {
		Copy[int](_first_int(a), _last_int(a), _first_int(c))
		ok := PrevPermutation[int](_first_int(a), _last_int(a))
		assert.Equal(LexicographicalCompare[int](_first_int(a), _last_int(a), _first_int(c), _last_int(c)), ok)
		last := _last_int(b)
		if Equal[int](_first_int(a), _last_int(a), _first_int(b), &last) {
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
	Iota(_first_int(a), _last_int(a), s+1)
	Generate(_first_int(b), _last_int(b), func() int { s++; return s })
	sliceEqual(assert.New(t), a, b)
}

func TestAccumulate(t *testing.T) {
	a := randIntSlice()
	sum := Accumulate(_first_int(a), _last_int(a), 0)
	sum2 := 0
	ForEach(_first_int(a), _last_int(a), func(it int) {
		sum2 += it
	})
	assert.New(t).Equal(sum, sum2)
}

func TestInnerProduct(t *testing.T) {
	a, b := randIntSlice(), randIntSlice()
	l := Min(len(a), len(b))
	p := InnerProduct(_first_int(a), _last_int(a[:l]), _first_int(b), 0)
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
	AdjacentDifference[int](_first_int(a), _last_int(a), _first_int(g))
	sliceEqual(assert, g, diff)
	PartialSum[int](_first_int(a), _last_int(a), _first_int(g))
	sliceEqual(assert, g, ps)
	ExclusiveScan(_first_int(a), _last_int(a), _first_int(g), 1)
	sliceEqual(assert, g, exc)
	InclusiveScan(_first_int(a), _last_int(a), _first_int(g), 2)
	sliceEqual(assert, g, inc)
	TransformExclusiveScan(_first_int(a), _last_int(a), _first_int(g), 3, func(x int) int { return x * 2 })
	sliceEqual(assert, g, exct)
	TransformInclusiveScan(_first_int(a), _last_int(a), _first_int(g), 4, func(x int) int { return x * x })
	sliceEqual(assert, g, inct)
}
