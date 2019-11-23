package iter_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	. "github.com/disksing/iter"
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

func randString() string {
	l := randInt()
	var s strings.Builder
	for i := 0; i < l; i++ {
		s.WriteByte('a' + byte(rand.Intn(26)))
	}
	return s.String()
}

var (
	begin    = SliceBegin
	end      = SliceEnd
	strBegin = StringBegin
	strEnd   = StringEnd
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
	n := rand.Intn(len(a) + 1)
	b = nil
	ForEachN(begin(a), n, f)
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
		assert.Equal(Count(begin(a), end(a), i), m[i])
	}
}

func TestMismatch(t *testing.T) {
	assert := assert.New(t)
	a, b := randIntSlice(), randIntSlice()
	var last2 ForwardReader
	if len(b) <= len(a) || rand.Intn(2) == 0 {
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
	c := byte('a' + byte(rand.Intn(26)))
	n := rand.Intn(10)
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
