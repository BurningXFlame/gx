/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package queue

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	as := require.New(t)

	tcs := [][]int{
		{},
		{8257},
		{3129, 752, 4994},
		{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			n := len(tc)
			q := New[int]()

			for _, e := range tc {
				q.Enq(e)
			}

			as.Equal(n, q.Len())

			e, ok := q.Peek()
			if n > 0 {
				as.True(ok)
				as.Equal(tc[0], e)
			} else {
				as.False(ok)
			}

			for i := 0; i < n; i++ {
				e, ok := q.Deq()
				as.True(ok)
				as.Equal(tc[i], e)
			}

			as.Equal(0, q.Len())
			_, ok = q.Peek()
			as.False(ok)
			_, ok = q.Deq()
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

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			q := New[int]()
			for _, e := range tc.in {
				q.Enq(e)
			}

			sum := 0
			q.Range(func(e int) bool {
				sum += e
				return true
			})
			as.Equal(tc.sum, sum)

			var even bool
			q.Range(func(e int) bool {
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

func TestForEach(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		sum int
	}{
		{[]int{}, 0},
		{[]int{8257}, 8257},
		{[]int{3129, 752, 4994}, 8875},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, 50648},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			q := New[int]()
			for _, e := range tc.in {
				q.Enq(e)
			}

			sum := 0
			q.ForEach(func(e int) {
				sum += e
			})
			as.Equal(tc.sum, sum)
		})
	}
}

func TestForEachOrder(t *testing.T) {
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
			q := New[int]()
			for _, e := range tc.in {
				q.Enq(e)
			}

			var l []int
			q.ForEach(func(e int) {
				l = append(l, e)
			})
			as.Equal(tc.in, l)
		})
	}
}
