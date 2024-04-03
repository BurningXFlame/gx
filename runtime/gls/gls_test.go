/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package gls

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGls(t *testing.T) {
	as := require.New(t)

	span0 := span{
		traceId: rand.Int(),
		pid:     0,
		id:      rand.Int(),
	}
	Put(&key, span0)

	spawn(as, span0, 8)

	cnt := 0
	mStore.Range(func(key, value any) bool {
		cnt++
		return true
	})
	as.Equal(1, cnt)
}

type span struct {
	traceId int
	pid     int
	id      int
}

var key int

func spawn(as *require.Assertions, parent span, n int) {
	if n < 2 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		Go(func() {
			defer wg.Done()

			sp, ok := Get(&key).(span)
			as.True(ok)
			as.Equal(parent, sp)

			sp.pid = parent.id
			sp.id = rand.Int()

			Put(&key, sp)

			spawn(as, sp, n/2)
		})
	}

	wg.Wait()
}
