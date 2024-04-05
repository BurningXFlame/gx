/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package stack

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
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
			s := New[int]()

			for _, e := range tc {
				s.Push(e)
			}

			as.Equal(n, s.Len())

			e, ok := s.Peek()
			if n > 0 {
				as.True(ok)
				as.Equal(tc[n-1], e)
			} else {
				as.False(ok)
			}

			for i := 0; i < n; i++ {
				e, ok := s.Pop()
				as.True(ok)
				as.Equal(tc[n-1-i], e)
			}

			as.Equal(0, s.Len())
			_, ok = s.Peek()
			as.False(ok)
			_, ok = s.Pop()
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
			s := New[int]()
			for _, e := range tc.in {
				s.Push(e)
			}

			sum := 0
			s.Range(func(e int) bool {
				sum += e
				return true
			})
			as.Equal(tc.sum, sum)

			var even bool
			s.Range(func(e int) bool {
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
			s := New[int]()
			for _, e := range tc.in {
				s.Push(e)
			}

			sum := 0
			s.ForEach(func(e int) {
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
			s := New[int]()
			for _, e := range tc.in {
				s.Push(e)
			}

			var l []int
			s.ForEach(func(e int) {
				l = append(l, e)
			})
			as.ElementsMatch(tc.in, l)

			n := len(tc.in)
			for i := 0; i < n; i++ {
				as.Equal(tc.in[i], l[n-1-i])
			}
		})
	}
}
