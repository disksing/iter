package iter_test

import (
	"sync/atomic"
	"testing"

	. "github.com/disksing/iter/v2"
	"github.com/stretchr/testify/assert"
)

// check if need skip the test. A test will be run count/count+1 times.
func skipAfter(t *testing.T, count int) {
	if atomic.LoadInt32(&testCounter) > int32(count) {
		t.SkipNow()
	}
}

var testCounter int32

func TestTouchSkipCounter(t *testing.T) {
	atomic.AddInt32(&testCounter, 1)
}

func TestMisc(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	it := IotaReader(100)
	assert.Equal(it.Read(), 100)
	assert.False(it.Eq(IotaReader(100)))
	it = it.Next()
	assert.Equal(it.Read(), 101)

	it2 := RepeatReader(100)
	assert.Equal(it2.Read(), 100)
	assert.False(__eq(it2, RepeatReader(100)))
	it2 = it2.Next()
	assert.Equal(it2.Read(), 100)

	g := IotaGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 101)
	g = RepeatGenerator(100)
	assert.Equal(g(), 100)
	assert.Equal(g(), 100)
}

func TestMakeString(t *testing.T) {
	skipAfter(t, 1)
	assert := assert.New(t)

	bs := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!'}
	assert.Equal(MakeString[byte](_first_byte(bs), _last_byte(bs)), "hello world!")

	rs := []rune{'改', '革', '春', '风', '吹', '满', '地'}
	assert.Equal(MakeString[rune](SliceBegin(rs), SliceEnd(rs)), "改革春风吹满地")
}
