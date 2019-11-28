package iter_test

import (
	"container/list"
	"fmt"
	"math/rand"
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
	l, h := randInt(), randInt()
	if l > h {
		l, h = h, l
	}
	s := make([]int, randInt())
	for i := range s {
		s[i] = l + r.Intn(h-l+1)
	}
	return s
}

func randString() string {
	l := randInt()
	var s strings.Builder
	for i := 0; i < l; i++ {
		s.WriteByte('a' + byte(r.Intn(26)))
	}
	return s.String()
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

func iterEqual(assert *assert.Assertions, a, b Iter) {
	assert.True(a.(Equalable).Equal(b), "a=%v\nb=%v", a, b)
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
	f := func(x Iter) { b = append(b, x.(Reader).Read().(int)) }
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
	for _, x := range a {
		count[x]++
	}
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
	assert.Equal(a[:n1], b[:n1])
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
	if it.(Equalable).Equal(end(a)) {
		if len(b) > 0 {
			it = Search(begin(a), end(a), begin(b), end(b))
		}
	} else {
		assert.True(Equal(begin(b), end(b), it, nil))
		it = FindEnd(NextReader(it), end(a), begin(b), end(b))
	}
	iterEqual(assert, it, end(a))
}

func TestFindFirstOf(t *testing.T) {
	a, b := randString(), randString()
	i := strings.IndexAny(a, b)
	if i == -1 {
		i = len(a)
	}
	iterEqual(assert.New(t),
		FindFirstOf(strBegin(a), strEnd(a), strBegin(b), strEnd(b)),
		AdvanceN(strBegin(a), i),
	)
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
	iterEqual(assert.New(t),
		AdjacentFind(begin(a), end(a)),
		AdvanceN(begin(a), res),
	)
}

func TestSearch(t *testing.T) {
	a, b := randString(), randString()
	i := strings.Index(a, b)
	if i == -1 {
		i = len(a)
	}
	iterEqual(assert.New(t),
		Search(strBegin(a), strEnd(a), strBegin(b), strEnd(b)),
		AdvanceN(strBegin(a), i),
	)
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
	iterEqual(assert.New(t),
		SearchN(strBegin(a), strEnd(a), n, c),
		AdvanceN(strBegin(a), i),
	)
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
	SliceErase(&a, Remove(begin(a), end(a), 1))
	SliceErase(&b, RemoveIf(begin(b), end(b), f))
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
	SliceErase(&b, Unique(begin(b), end(b)))
	SliceErase(&c, UniqueCopy(begin(a), end(a), begin(c)))
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

func (l *forwardListIter) Equal(x Any) bool {
	return l.e == x.(*forwardListIter).e
}

func (l *forwardListIter) Next() ForwardIter {
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

func TestCompare(t *testing.T) {
	assert := assert.New(t)
	a, b := randString(), randString()
	x, y, z, w := StringBegin(a), StringEnd(a), StringBegin(b), StringEnd(b)
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

func TestEqual(t *testing.T) {
	assert := assert.New(t)
	a, b := randString(), randString()
	if len(a) > len(b) {
		a, b = b, a
	}
	assert.Equal(Equal(StringBegin(a), StringEnd(a), StringBegin(b), nil), a == b[:len(a)])
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
		iterEqual(assert, min, end(s))
		iterEqual(assert, max, end(s))
		iterEqual(assert, min2, end(s))
		iterEqual(assert, max2, end(s))
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
