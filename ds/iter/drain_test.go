/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForEach(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in []int
	}{
		{nil},
		{[]int{8257}},
		{[]int{3129, 752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			it := FromSlice(tc.in)

			c := new(sliceCollector[int])

			it.ForEach(c.Add)

			as.Equal(tc.in, c.l)
		})
	}
}

func TestMin(t *testing.T) {
	as := require.New(t)

	less := func(a, b int) bool {
		return a < b
	}

	tcs := []struct {
		in  []int
		out Option[int]
	}{
		{nil, Option[int]{}},
		{[]int{8257}, Some(8257)},
		{[]int{3129, 752, 4994}, Some(752)},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, Some(2615)},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			val := FromSlice(tc.in).Min(less)
			as.Equal(tc.out, val)
		})
	}
}

func TestMax(t *testing.T) {
	as := require.New(t)

	less := func(a, b int) bool {
		return a < b
	}

	tcs := []struct {
		in  []int
		out Option[int]
	}{
		{nil, Option[int]{}},
		{[]int{8257}, Some(8257)},
		{[]int{3129, 752, 4994}, Some(4994)},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, Some(9715)},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			val := FromSlice(tc.in).Max(less)
			as.Equal(tc.out, val)
		})
	}
}

func TestSort(t *testing.T) {
	as := require.New(t)

	less := func(a, b int) bool {
		return a < b
	}

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, []int{8257}},
		{[]int{3129, 752, 4994}, []int{752, 3129, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{2615, 3690, 4421, 5015, 5613, 5955, 6038, 7586, 9715}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).Sort(less)
			as.Equal(tc.out, l)
		})
	}
}

func TestReduce(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out int
	}{
		{nil, 0},
		{[]int{8257}, 8257},
		{[]int{3129, 752, 4994}, 8875},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, 50648},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			val := FromSlice(tc.in).Reduce(0, func(acc int, e int) int {
				return acc + e
			})
			as.Equal(tc.out, val)
		})
	}
}

func TestReduceToAnotherType(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out string
	}{
		{nil, ""},
		{[]int{8257}, "8257"},
		{[]int{3129, 752, 4994}, "31297524994"},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, "758656133690261544219715595560385015"},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			it := FromSlice(tc.in)
			val := Reduce(it, "", func(acc string, e int) string {
				return acc + strconv.Itoa(e)
			})
			as.Equal(tc.out, val)
		})
	}
}

func TestRange(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, nil},
		{[]int{3129, 752, 4994}, nil},
		{[]int{752, 3129, 4994}, []int{752}},
		{[]int{752, 4994, 3129}, []int{752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586}},
		{[]int{7586, 3690, 5613, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 3690}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			it := FromSlice(tc.in)

			c := new(sliceCollector[int])

			it.Range(func(e int) bool {
				if e%2 == 1 {
					return false
				}

				c.Add(e)
				return true
			})

			as.Equal(tc.out, c.l)
		})
	}
}

func TestAny(t *testing.T) {
	as := require.New(t)

	pred := func(e int) bool {
		return e%2 == 0
	}

	tcs := []struct {
		in  []int
		out bool
	}{
		{nil, false},
		{[]int{8257}, false},
		{[]int{3129, 752, 4994}, true},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, true},
		{[]int{5613, 2615, 4421, 9715, 5955, 5015}, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			ok := FromSlice(tc.in).Any(pred)
			as.Equal(tc.out, ok)
		})
	}
}

func TestEvery(t *testing.T) {
	as := require.New(t)

	pred := func(e int) bool {
		return e%2 == 0
	}

	tcs := []struct {
		in  []int
		out bool
	}{
		{nil, true},
		{[]int{8257}, false},
		{[]int{752}, true},
		{[]int{3129, 752, 4994}, false},
		{[]int{752, 4994}, true},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, false},
		{[]int{5613, 2615, 4421, 9715, 5955, 5015}, false},
		{[]int{7586, 3690, 6038}, true},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			ok := FromSlice(tc.in).Every(pred)
			as.Equal(tc.out, ok)
		})
	}
}
