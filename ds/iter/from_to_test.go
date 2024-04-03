/*
GX (https://github.com/BurningXFlame/gx).
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

func TestFromToSlice(t *testing.T) {
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
			l := FromSlice(tc.in).ToSlice()
			as.Equal(tc.in, l)
		})
	}
}

func TestFromToMap(t *testing.T) {
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
			m := make(map[int]string)
			for _, v := range tc.in {
				m[v] = strconv.Itoa(v)
			}

			it := FromMap(m)
			m2 := ToMap(it)

			as.Equal(m, m2)
		})
	}
}

func TestFromToChan(t *testing.T) {
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
			ch := make(chan int)
			go func() {
				for _, v := range tc.in {
					ch <- v
				}
				close(ch)
			}()

			ch2 := FromChan((<-chan int)(ch)).ToChan()

			var l []int
			for e := range ch2 {
				l = append(l, e)
			}

			as.Equal(tc.in, l)
		})
	}
}

func TestFromRangerToCollector(t *testing.T) {
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
			r := &sliceRanger[int]{l: tc.in}
			it := FromRanger[int](r)

			c := new(sliceCollector[int])
			it.To(c)

			as.Equal(tc.in, c.l)
		})
	}
}

type sliceRanger[E any] struct {
	l []E
}

func (r *sliceRanger[E]) ForEach(fn func(E)) {
	for _, e := range r.l {
		fn(e)
	}
}

type sliceCollector[E any] struct {
	l []E
}

func (c *sliceCollector[E]) Add(e E) {
	c.l = append(c.l, e)
}

func TestFromIterToCollector(t *testing.T) {
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
			it := &sliceIter[int]{src: tc.in}

			c := new(sliceCollector[int])

			From[int](it).To(c)

			as.Equal(tc.in, c.l)
		})
	}
}

func TestFromMultiIterToCollector(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		in2 []int
	}{
		{nil, nil},
		{nil, []int{8257}},
		{[]int{3129, 752, 4994}, nil},
		{[]int{8257}, []int{3129, 752, 4994}},
		{[]int{3129, 752, 4994}, []int{8257}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			it := &sliceIter[int]{src: tc.in}

			ch := make(chan int)
			go func() {
				for _, v := range tc.in2 {
					ch <- v
				}
				close(ch)
			}()
			it2 := &chanIter[int]{src: ch}

			c := new(sliceCollector[int])

			From[int](it, it2).To(c)

			var expect []int
			size := len(tc.in) + len(tc.in2)
			if size > 0 {
				expect = make([]int, size)
			}
			n := copy(expect, tc.in)
			copy(expect[n:], tc.in2)

			as.Equal(expect, c.l)
		})
	}
}
