package iter_test

import (
	"testing"

	. "github.com/disksing/iter/v2"
	"github.com/stretchr/testify/assert"
)

func TestMisc(t *testing.T) {
	assert := assert.New(t)

	it := IotaReader(100)
	assert.Equal(it.Read(), 100)
	assert.False(it.Eq(IotaReader(100)))
	it = it.Next()
	assert.Equal(it.Read(), 101)

	it2 := RepeatReader(100)
	assert.Equal(it2.Read(), 100)
	assert.False(it2.Eq(RepeatReader(100)))
	it2 = it2.Next()
	assert.Equal(it2.Read(), 100)

	g := IotaGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 101)
	g = RepeatGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 100)
}
