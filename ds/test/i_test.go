/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/ds/heap"
	"github.com/burningxflame/gx/ds/iter"
	"github.com/burningxflame/gx/ds/queue"
	"github.com/burningxflame/gx/ds/ringbuf"
	"github.com/burningxflame/gx/ds/set"
	"github.com/burningxflame/gx/ds/stack"
)

func TestIteratorFromToRingbuf(t *testing.T) {
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
			r := ringbuf.New[int](len(tc.in))
			for _, v := range tc.in {
				r.PushBack(v)
			}

			r2 := ringbuf.New[int](len(tc.in))

			iter.FromRanger[int](r).To(r2)

			var l []int
			r2.ForEach(func(e int) {
				l = append(l, e)
			})

			as.Equal(tc.in, l)
		})
	}
}

func TestIteratorFromToQueue(t *testing.T) {
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
			q := queue.New[int](len(tc.in))
			for _, v := range tc.in {
				q.Enq(v)
			}

			q2 := queue.New[int](len(tc.in))

			iter.FromRanger[int](q).To(q2)

			var l []int
			q2.ForEach(func(e int) {
				l = append(l, e)
			})

			as.Equal(tc.in, l)
		})
	}
}

func TestIteratorFromToStack(t *testing.T) {
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
			s := stack.New[int](len(tc.in))
			for _, v := range tc.in {
				s.Push(v)
			}

			s2 := stack.New[int](len(tc.in))

			iter.FromRanger[int](s).To(s2)

			var l []int
			s2.ForEach(func(e int) {
				l = append(l, e)
			})

			as.Equal(tc.in, l)
		})
	}
}

func TestIteratorFromToSet(t *testing.T) {
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
			s := set.New[int](len(tc.in))
			for _, v := range tc.in {
				s.Add(v)
			}

			s2 := set.New[int](len(tc.in))

			iter.FromRanger[int](s).To(s2)

			var l []int
			s2.ForEach(func(e int) {
				l = append(l, e)
			})

			as.ElementsMatch(tc.in, l)
		})
	}
}

func TestIteratorFromToHeap(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in []int
	}{
		{nil},
		{[]int{8257}},
		{[]int{3129, 752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}},
	}

	less := func(a, b int) bool {
		return a < b
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			h := heap.New[int](less, len(tc.in))
			for _, v := range tc.in {
				h.Push(v)
			}

			h2 := heap.New[int](less, len(tc.in))

			iter.FromRanger[int](h).To(h2)

			var l []int
			h2.ForEach(func(e int) {
				l = append(l, e)
			})

			as.ElementsMatch(tc.in, l)
		})
	}
}
