package iter_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

func randInt() int {
	return rand.Intn(100)
}

func randIntSlice() []int {
	l, h := randInt(), randInt()
	if l > h {
		l, h = h, l
	}
	s := make([]int, randInt())
	for i := range s {
		s[i] = l + rand.Intn(h-l+1)
	}
	return s
}

var (
	begin = iter.SliceBegin
	end   = iter.SliceEnd
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func sliceEqual(assert *assert.Assertions, a, b []int) {
	if len(a) == 0 && len(b) == 0 {
		return
	}
	assert.Equal(a, b)
}

func iterEqual(assert *assert.Assertions, a, b interface{}) {
	assert.True(a.(iter.Equalable).Equal(b), "a=%v\nb=%v", a, b)
}

func TestAllAnyNoneOf(t *testing.T) {
	assert := assert.New(t)
	pred := func(x iter.Any) bool { return x.(int)%2 == 0 }
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
	assert.Equal(iter.AllOf(begin(s), end(s), pred), allOf(s))
	assert.Equal(iter.AnyOf(begin(s), end(s), pred), anyOf(s))
	assert.Equal(iter.NoneOf(begin(s), end(s), pred), noneOf(s))
}

func TestForEach(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	var b []int
	f := func(x iter.Iter) { b = append(b, x.(iter.Readable).Read().(int)) }
	iter.ForEach(begin(a), end(a), f)
	sliceEqual(assert, a, b)
	n := rand.Intn(len(a) + 1)
	b = nil
	iter.ForEachN(begin(a), n, f)
	sliceEqual(assert, a[:n], b)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	m := make(map[int]int)
	for _, x := range a {
		m[x]++
	}
	for i := 0; i < 100; i++ {
		assert.Equal(iter.Count(begin(a), end(a), i), m[i])
	}
}

func TestMismatch(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	var last2 iter.ForwardReader
	if len(b) <= len(a) || rand.Intn(2) == 0 {
		last2 = end(b)
	}
	it1, it2 := iter.Mismatch(begin(a), end(a), begin(b), last2)
	n1, n2 := iter.Distance(begin(a), it1), iter.Distance(begin(b), it2)
	assert.Equal(n1, n2)
	assert.Equal(a[:n1], b[:n1])
	assert.True((n1 >= len(a)) || (n1 >= len(b)) || a[n1] != b[n1])
}

func TestFind(t *testing.T) {
	assert := assert.New(t)
	a := randIntSlice()
	f := func(x iter.Any) bool { return x.(int)%2 == 0 }
	v := randInt()
	it := iter.Find(begin(a), end(a), v)
	assert.True(iter.NoneOf(begin(a), it, func(x iter.Any) bool { return x.(int) == v }))
	if n := iter.Distance(begin(a), it); n < len(a) {
		assert.Equal(a[n], v)
	}
	it = iter.FindIf(begin(a), end(a), f)
	assert.True(iter.NoneOf(begin(a), it, f))
	if n := iter.Distance(begin(a), it); n < len(a) {
		assert.True(f(a[n]))
	}
	it = iter.FindIfNot(begin(a), end(a), f)
	assert.True(iter.AllOf(begin(a), it, f))
	if n := iter.Distance(begin(a), it); n < len(a) {
		assert.False(f(a[n]))
	}
}

func TestFindEnd(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	it := iter.FindEnd(begin(a), end(a), begin(b), end(b))
	if it.(iter.Equalable).Equal(end(a)) {
		if len(b) > 0 {
			it = iter.Search(begin(a), end(a), begin(b), end(b))
		}
	} else {
		assert.True(iter.Equal(begin(b), end(b), it, nil))
		it = iter.FindEnd(iter.NextReader(it), end(a), begin(b), end(b))
	}
	iterEqual(assert, it, end(a))
}

func TestMinmax(t *testing.T) {
	assert := assert.New(t)
	a, b := randInt(), randInt()
	min, max := iter.Minmax(a, b)
	assert.LessOrEqual(min, max)
	assert.Equal(iter.Max(a, b), max)
	assert.Equal(iter.Min(a, b), min)
}

func TestMinmaxElement(t *testing.T) {
	assert := assert.New(t)
	s := randIntSlice()
	min, max := iter.MinmaxElement(begin(s), end(s))
	min2, max2 := iter.MinElement(begin(s), end(s)), iter.MaxElement(begin(s), end(s))
	assert.True(iter.NoneOf(begin(s), end(s), func(v iter.Any) bool { return v.(int) > max.Read().(int) || v.(int) < min.Read().(int) }))
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
	l, h := iter.Minmax(randInt(), randInt())
	v := randInt()
	c := iter.Clamp(v, l, h)
	if c != v {
		assert.True(v < l.(int) || v > h.(int))
	}
	assert.GreaterOrEqual(c, l)
	assert.LessOrEqual(c, h)
}
