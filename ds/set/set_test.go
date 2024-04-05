/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package set

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		add    []int
		delete []int
		len    int
	}{
		{[]int{}, []int{}, 0},
		{[]int{}, []int{8257}, 0},
		{[]int{8257}, []int{}, 1},
		{[]int{8257}, []int{8257}, 0},
		{[]int{3129, 752, 4994}, []int{}, 3},
		{[]int{3129, 752, 4994}, []int{752}, 2},
		{[]int{3129, 752, 4994}, []int{8257}, 3},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{}, 9},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 2615, 752}, 7},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}, 9},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.add {
				s.Add(e)
			}
			as.Equal(len(tc.add), s.Len())
			for _, e := range tc.add {
				as.True(s.Contain(e))
			}

			for _, e := range tc.add {
				s.Add(e)
			}
			as.Equal(len(tc.add), s.Len())

			for _, e := range tc.delete {
				s.Delete(e)
			}
			as.Equal(tc.len, s.Len())
			for _, e := range tc.delete {
				as.False(s.Contain(e))
			}

			for _, e := range tc.delete {
				s.Delete(e)
			}
			as.Equal(tc.len, s.Len())
		})
	}
}

func TestEqual(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		a     []int
		b     []int
		equal bool
	}{
		{[]int{}, []int{}, true},
		{[]int{}, []int{8257}, false},
		{[]int{8257}, []int{}, false},
		{[]int{8257}, []int{8257}, true},
		{[]int{3129, 752, 4994}, []int{}, false},
		{[]int{3129, 752, 4994}, []int{752}, false},
		{[]int{3129, 752, 4994}, []int{4421, 9715, 5955}, false},
		{[]int{3129, 752, 4994}, []int{3129, 752, 4994}, true},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 752, 5015}, false},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 5015, 6038}, []int{9715, 5613, 7586, 2615, 3690, 4421, 5955, 6038, 5015}, true},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.a {
				s.Add(e)
			}
			for _, e := range tc.a {
				s.Add(e)
			}

			b := New[int]()
			for _, e := range tc.b {
				b.Add(e)
			}

			as.Equal(tc.equal, s.Equal(b))
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
				s.Add(e)
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
				s.Add(e)
			}

			sum := 0
			s.ForEach(func(e int) {
				sum += e
			})
			as.Equal(tc.sum, sum)
		})
	}
}

func TestUnion(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		a   []int
		b   []int
		len int
	}{
		{[]int{}, []int{}, 0},
		{[]int{}, []int{8257}, 1},
		{[]int{8257}, []int{}, 1},
		{[]int{8257}, []int{8257}, 1},
		{[]int{3129, 752, 4994}, []int{}, 3},
		{[]int{3129, 752, 4994}, []int{752}, 3},
		{[]int{3129, 752, 4994}, []int{8257}, 4},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 2615, 752}, 10},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}, 12},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 752, 5015}, 10},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 5015, 6038}, []int{9715, 5613, 7586, 2615, 3690, 4421, 5955, 6038, 5015}, 9},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.a {
				s.Add(e)
			}

			b := New[int]()
			for _, e := range tc.b {
				b.Add(e)
			}

			r := s.Union(b)
			as.Equal(tc.len, r.Len())

			for _, e := range tc.a {
				as.True(r.Contain(e))
			}
			for _, e := range tc.b {
				as.True(r.Contain(e))
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		a   []int
		b   []int
		len int
	}{
		{[]int{}, []int{}, 0},
		{[]int{}, []int{8257}, 0},
		{[]int{8257}, []int{}, 0},
		{[]int{8257}, []int{8257}, 1},
		{[]int{3129, 752, 4994}, []int{}, 0},
		{[]int{3129, 752, 4994}, []int{752}, 1},
		{[]int{3129, 752, 4994}, []int{8257}, 0},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 2615, 752}, 2},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}, 0},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 752, 5015}, 8},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 5015, 6038}, []int{9715, 5613, 7586, 2615, 3690, 4421, 5955, 6038, 5015}, 9},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.a {
				s.Add(e)
			}

			b := New[int]()
			for _, e := range tc.b {
				b.Add(e)
			}

			r := s.Intersect(b)
			as.Equal(tc.len, r.Len())

			r.Range(func(e int) bool {
				as.True(s.Contain(e))
				as.True(b.Contain(e))
				return true
			})
		})
	}
}

func TestDiff(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		a   []int
		b   []int
		len int
	}{
		{[]int{}, []int{}, 0},
		{[]int{}, []int{8257}, 0},
		{[]int{8257}, []int{}, 1},
		{[]int{8257}, []int{8257}, 0},
		{[]int{3129, 752, 4994}, []int{}, 3},
		{[]int{3129, 752, 4994}, []int{752}, 2},
		{[]int{3129, 752, 4994}, []int{8257}, 3},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 2615, 752}, 7},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}, 9},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 752, 5015}, 1},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 5015, 6038}, []int{9715, 5613, 7586, 2615, 3690, 4421, 5955, 6038, 5015}, 0},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.a {
				s.Add(e)
			}

			b := New[int]()
			for _, e := range tc.b {
				b.Add(e)
			}

			r := s.Diff(b)
			as.Equal(tc.len, r.Len())

			r.Range(func(e int) bool {
				as.True(s.Contain(e))
				as.False(b.Contain(e))
				return true
			})
		})
	}
}

func TestUnionIntersectDiff(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		a []int
		b []int
	}{
		{[]int{}, []int{}},
		{[]int{}, []int{8257}},
		{[]int{8257}, []int{}},
		{[]int{8257}, []int{8257}},
		{[]int{3129, 752, 4994}, []int{}},
		{[]int{3129, 752, 4994}, []int{752}},
		{[]int{3129, 752, 4994}, []int{8257}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 2615, 752}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{3129, 752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 752, 5015}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 5015, 6038}, []int{9715, 5613, 7586, 2615, 3690, 4421, 5955, 6038, 5015}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			s := New[int]()
			for _, e := range tc.a {
				s.Add(e)
			}

			b := New[int]()
			for _, e := range tc.b {
				b.Add(e)
			}

			ru := s.Union(b)
			ru.Intersect(s).Equal(s)
			ru.Intersect(b).Equal(b)

			ri := s.Intersect(b)
			ri.Intersect(s).Equal(ri)
			ri.Intersect(b).Equal(ri)
			ri.Union(s).Equal(s)
			ri.Union(b).Equal(b)

			rd := s.Diff(b)
			rd.Intersect(s).Equal(rd)
			as.Equal(0, rd.Intersect(b).Len())
			rd.Union(ri).Equal(s)
		})
	}
}
