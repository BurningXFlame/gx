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

func TestFilter(t *testing.T) {
	as := require.New(t)

	pred := func(e int) bool {
		return e%2 == 0
	}

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, nil},
		{[]int{3129, 752, 4994}, []int{752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{7586, 3690, 6038}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).Filter(pred).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestMap(t *testing.T) {
	as := require.New(t)

	fn := func(e int) int {
		return e * 2
	}

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, []int{16514}},
		{[]int{3129, 752, 4994}, []int{6258, 1504, 9988}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{15172, 11226, 7380, 5230, 8842, 19430, 11910, 12076, 10030}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).Map(fn).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestMapToAnotherType(t *testing.T) {
	as := require.New(t)

	fn := func(e int) string {
		return strconv.Itoa(e * 2)
	}

	tcs := []struct {
		in  []int
		out []string
	}{
		{nil, nil},
		{[]int{8257}, []string{"16514"}},
		{[]int{3129, 752, 4994}, []string{"6258", "1504", "9988"}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []string{"15172", "11226", "7380", "5230", "8842", "19430", "11910", "12076", "10030"}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			it := FromSlice(tc.in)
			l := Map(it, fn).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestTake(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		n   int
		out []int
	}{
		{nil, 0, nil},
		{nil, 2, nil},
		{[]int{8257}, 0, nil},
		{[]int{8257}, 1, []int{8257}},
		{[]int{8257}, 2, []int{8257}},
		{[]int{8257}, 100, []int{8257}},
		{[]int{3129, 752, 4994}, 0, nil},
		{[]int{3129, 752, 4994}, 1, []int{3129}},
		{[]int{3129, 752, 4994}, 2, []int{3129, 752}},
		{[]int{3129, 752, 4994}, 3, []int{3129, 752, 4994}},
		{[]int{3129, 752, 4994}, 4, []int{3129, 752, 4994}},
		{[]int{3129, 752, 4994}, 100, []int{3129, 752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, 6, []int{7586, 5613, 3690, 2615, 4421, 9715}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, -1, nil},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).Take(tc.n).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestDrop(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		n   int
		out []int
	}{
		{nil, 0, nil},
		{nil, 2, nil},
		{[]int{8257}, 0, []int{8257}},
		{[]int{8257}, 1, nil},
		{[]int{8257}, 100, nil},
		{[]int{8257}, -1, []int{8257}},
		{[]int{8257}, -100, []int{8257}},
		{[]int{3129, 752, 4994}, 3, nil},
		{[]int{3129, 752, 4994}, 2, []int{4994}},
		{[]int{3129, 752, 4994}, 1, []int{752, 4994}},
		{[]int{3129, 752, 4994}, 0, []int{3129, 752, 4994}},
		{[]int{3129, 752, 4994}, -1, []int{3129, 752, 4994}},
		{[]int{3129, 752, 4994}, -100, []int{3129, 752, 4994}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, 6, []int{5955, 6038, 5015}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, -1, []int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).Drop(tc.n).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestTakeWhile(t *testing.T) {
	as := require.New(t)

	pred := func(e int) bool {
		return e%2 == 0
	}

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
			l := FromSlice(tc.in).TakeWhile(pred).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestDropWhile(t *testing.T) {
	as := require.New(t)

	pred := func(e int) bool {
		return e%2 == 0
	}

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, []int{8257}},
		{[]int{3129, 752, 4994}, []int{3129, 752, 4994}},
		{[]int{752, 3129, 4994}, []int{3129, 4994}},
		{[]int{752, 4994, 3129}, []int{3129}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}},
		{[]int{7586, 3690, 5613, 2615, 4421, 9715, 5955, 6038, 5015}, []int{5613, 2615, 4421, 9715, 5955, 6038, 5015}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).DropWhile(pred).ToSlice()
			as.Equal(tc.out, l)
		})
	}
}

func TestChain(t *testing.T) {
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
			it := FromSlice(tc.in)

			ch := make(chan int)
			go func() {
				for _, v := range tc.in2 {
					ch <- v
				}
				close(ch)
			}()
			it2 := FromChan((<-chan int)(ch))

			l := Chain(it, it2).ToSlice()

			var expect []int
			size := len(tc.in) + len(tc.in2)
			if size > 0 {
				expect = make([]int, size)
			}
			n := copy(expect, tc.in)
			copy(expect[n:], tc.in2)

			as.Equal(expect, l)
		})
	}
}

func TestPipeline(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  []int
		out []int
	}{
		{nil, nil},
		{[]int{8257}, nil},
		{[]int{3129, 752, 4994}, nil},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{5230, 8842, 19430, 11910, 10030}},
		{[]int{8257, 3129, 752, 4994}, nil},
		{[]int{8257, 7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{5230, 8842, 19430, 11910, 10030}},
		{[]int{3129, 752, 4994, 7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{11226, 5230, 8842, 19430, 11910, 10030}},
		{[]int{8257, 3129, 752, 4994, 7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015}, []int{11226, 5230, 8842, 19430, 11910, 10030}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015, 3129, 752, 4994}, []int{5230, 8842, 19430, 11910, 10030, 6258}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015, 8257}, []int{5230, 8842, 19430, 11910, 10030, 16514}},
		{[]int{7586, 5613, 3690, 2615, 4421, 9715, 5955, 6038, 5015, 3129, 752, 4994, 8257}, []int{5230, 8842, 19430, 11910, 10030, 6258}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			l := FromSlice(tc.in).
				Drop(3).
				Filter(func(e int) bool {
					return e%2 == 1
				}).
				Map(func(e int) int {
					return e * 2
				}).
				Take(6).
				ToSlice()

			as.Equal(tc.out, l)
		})
	}
}
