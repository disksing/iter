package strs_test

import (
	"fmt"
	"testing"

	"github.com/disksing/iter/v2"
	"github.com/disksing/iter/v2/algo"
	"github.com/disksing/iter/v2/slices"
	. "github.com/disksing/iter/v2/strs"
	"github.com/stretchr/testify/assert"
)

func TestStringIterator(t *testing.T) {
	assert := assert.New(t)
	s := "abcdefg"
	assert.Equal(MakeString[byte](RBegin(s), REnd(s)), "gfedcba")

	Begin(s).AllowMultiplePass()

	assert.Contains(fmt.Sprint(Begin(s)), "->")
	assert.Contains(fmt.Sprint(RBegin(s)), "<-")

	assert.Equal(End(s).Prev().Read(), byte('g'))
	assert.Equal(REnd(s).Prev().Read(), byte('a'))

	assert.Equal(
		Begin(s).AdvanceN(3).Read(),
		byte('d'),
	)
	assert.Equal(
		RBegin(s).AdvanceN(3).Read(),
		byte('d'),
	)
	assert.Equal(
		End(s).AdvanceN(-3).Read(),
		byte('e'),
	)
	assert.Equal(
		REnd(s).AdvanceN(-3).Read(),
		byte('c'),
	)

	assert.Equal(Begin(s).Distance(End(s)), len(s))
	assert.Equal(End(s).Distance(Begin(s)), -len(s))
	assert.Equal(RBegin(s).Distance(REnd(s)), len(s))
	assert.Equal(REnd(s).Distance(RBegin(s)), -len(s))

	assert.True(Begin(s).Less(End(s)))
	assert.False(End(s).Less(Begin(s)))
	assert.True(RBegin(s).Less(REnd(s)))
	assert.False(REnd(s).Less(RBegin(s)))
}

func TestStringBuilder(t *testing.T) {
	assert := assert.New(t)

	var bs StringBuilderInserter[any]
	algo.FillN[any](&bs, 3, 'a')
	assert.Equal(bs.String(), "aaa")

	bs.Reset()
	algo.FillN[any](&bs, 3, rune('香'))
	assert.Equal(bs.String(), "香香香")

	bs.Reset()
	bs.Delimiter = ","
	algo.FillN[any](&bs, 3, []byte("abc"))
	assert.Equal(bs.String(), "abc,abc,abc")

	bs.Reset()
	bs.Delimiter = "-->"
	algo.FillN[any](&bs, 3, "abc")
	assert.Equal(bs.String(), "abc-->abc-->abc")

	var bsi StringBuilderInserter[int]
	bs.Delimiter = ""
	algo.CopyN[int](iter.IotaReader(1), 5, &bsi)
	assert.Equal(bsi.String(), "12345")
}

func TestMakeString(t *testing.T) {
	assert := assert.New(t)

	bs := []byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '!'}
	assert.Equal(MakeString[byte](slices.Begin(bs), slices.End(bs)), "hello world!")

	rs := []rune{'改', '革', '春', '风', '吹', '满', '地'}
	assert.Equal(MakeString[rune](slices.Begin(rs), slices.End(rs)), "改革春风吹满地")
}
