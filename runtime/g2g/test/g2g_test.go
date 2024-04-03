/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package test

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	as := require.New(t)

	tcs := []func() int32{c1, c2, c3, c4, c5}
	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			as.Equal(total, tc())
		})
	}
}

const (
	n           = 100
	total int32 = 5050
)

func c1() int32 {
	var total int32

	var wg sync.WaitGroup
	wg.Add(n)
	for e := 1; e <= n; e++ {
		e := e
		go func() {
			defer wg.Done()
			atomic.AddInt32(&total, int32(e))
		}()
	}
	wg.Wait()

	return total
}

func c2() int32 {
	var total int32

	var wg sync.WaitGroup
	wg.Add(n)
	for e := 1; e <= n; e++ {
		go func(e int) {
			defer wg.Done()
			atomic.AddInt32(&total, int32(e))
		}(e)
	}
	wg.Wait()

	return total
}

func c3() int32 {
	var total int32

	var wg sync.WaitGroup
	wg.Add(n)
	for e := 1; e <= n; e++ {
		go _c3(&total, int32(e), &wg)
	}
	wg.Wait()

	return total
}

func _c3(total *int32, e int32, wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt32(total, e)
}

func c4() int32 {
	var total int32

	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i + 1
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for _, e := range l {
		go _c3(&total, int32(e), &wg)
	}
	wg.Wait()

	return total
}

func c5() int32 {
	var total int32

	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i + 1
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i, e := range l {
		go _c3(&total, int32(e), &wg)
		_ = i
	}
	wg.Wait()

	return total
}
