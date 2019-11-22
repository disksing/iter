package iter_test

import (
	"testing"

	"github.com/disksing/iter"
	"github.com/stretchr/testify/assert"
)

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
	if len(s) > 0 {
		assert.True(iter.NoneOf(begin(s), end(s), func(v iter.Any) bool { return v.(int) > max.Read().(int) || v.(int) < min.Read().(int) }))
		assert.Equal(iter.MinElement(begin(s), end(s)).Read(), min.Read())
		assert.Equal(iter.MaxElement(begin(s), end(s)).Read(), max.Read())
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
