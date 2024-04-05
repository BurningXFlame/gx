/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package heap

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinHeap(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, []int{8257}},
		{[]int{3129, 752, 4994}, []int{752, 3129, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{2615, 3690, 4421, 5015, 5613, 5955, 6038, 7586, 9715}},
	}

	less := func(a, b int) bool {
		return a < b
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			n := len(tc.in)
			h := New(less, n)
			for _, e := range tc.in {
				h.Push(e)
			}
			as.Equal(n, h.Len())

			var actual []int
			for i := 0; i < n; i++ {
				e, ok := h.Pop()
				as.True(ok)
				actual = append(actual, e)
			}
			as.Equal(tc.out, actual)
			as.Equal(0, h.Len())

			_, ok := h.Pop()
			as.False(ok)
		})
	}
}

func TestMaxHeap(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, []int{8257}},
		{[]int{3129, 752, 4994}, []int{4994, 3129, 752}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{9715, 7586, 6038, 5955, 5613, 5015, 4421, 3690, 2615}},
	}

	less := func(a, b int) bool {
		return a > b
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			n := len(tc.in)
			h := New(less, n)
			for _, e := range tc.in {
				h.Push(e)
			}
			as.Equal(n, h.Len())

			var actual []int
			for i := 0; i < n; i++ {
				e, ok := h.Pop()
				as.True(ok)
				actual = append(actual, e)
			}
			as.Equal(tc.out, actual)
			as.Equal(0, h.Len())

			_, ok := h.Pop()
			as.False(ok)
		})
	}
}

func TestRange(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in   []int
		sum  int
		even bool
	}{
		{[]int{}, 0, false},
		{[]int{8257}, 8257, false},
		{[]int{3129, 752, 4994}, 8875, true},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, 50648, true},
	}

	less := func(a, b int) bool {
		return a < b
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			h := New(less)
			for _, e := range tc.in {
				h.Push(e)
			}

			sum := 0
			h.Range(func(e int) bool {
				sum += e
				return true
			})
			as.Equal(tc.sum, sum)

			var even bool
			h.Range(func(e int) bool {
				if e%2 == 0 {
					even = true
					return false
				}
				return true
			})
			as.Equal(tc.even, even)
		})
	}
}
