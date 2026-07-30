package algo_test

import (
	"slices"
	"testing"

	"github.com/disksing/iter/v2/algo"
	iterslices "github.com/disksing/iter/v2/slices"
)

func FuzzSort(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1})
	f.Add([]byte{3, 1, 2, 1, 0, 255})

	f.Fuzz(func(t *testing.T, input []byte) {
		if len(input) > 4096 {
			t.Skip()
		}

		got := slices.Clone(input)
		want := slices.Clone(input)
		algo.Sort[byte](iterslices.Begin(got), iterslices.End(got))
		slices.Sort(want)

		if !slices.Equal(got, want) {
			t.Fatalf("Sort() = %v, want %v", got, want)
		}
	})
}

func FuzzNthElement(f *testing.F) {
	f.Add([]byte{3, 1, 2}, uint64(0))
	f.Add([]byte{3, 1, 2}, uint64(1))
	f.Add([]byte{3, 1, 2}, uint64(2))
	f.Add([]byte{5, 5, 1, 9, 0, 5}, uint64(99))

	f.Fuzz(func(t *testing.T, input []byte, rawIndex uint64) {
		if len(input) == 0 || len(input) > 4096 {
			t.Skip()
		}

		nth := int(rawIndex % uint64(len(input)))
		got := slices.Clone(input)
		want := slices.Clone(input)
		slices.Sort(want)

		algo.NthElement[byte](
			iterslices.Begin(got),
			iterslices.Begin(got).AdvanceN(nth),
			iterslices.End(got),
		)

		if got[nth] != want[nth] {
			t.Fatalf("NthElement()[%d] = %d, want %d; result %v", nth, got[nth], want[nth], got)
		}
		for i := range got {
			if i < nth && got[i] > got[nth] {
				t.Fatalf("left partition contains %d > pivot %d: %v", got[i], got[nth], got)
			}
			if i > nth && got[i] < got[nth] {
				t.Fatalf("right partition contains %d < pivot %d: %v", got[i], got[nth], got)
			}
		}
	})
}
