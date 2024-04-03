/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package gid

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGid(t *testing.T) {
	as := require.New(t)

	g1 := Gid()

	t.Run("", func(t *testing.T) {
		g2 := Gid()
		as.NotEqual(g1, g2)
		fmt.Printf("g1: %v, g2: %v\n", g1, g2)
	})
}

func TestGidCollision(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := 10000
	var m sync.Map

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			id := Gid()
			m.Store(id, struct{}{})

			wg.Done()
			<-ctx.Done()
		}()
	}

	wg.Wait()

	cnt := 0
	m.Range(func(key, value any) bool {
		cnt++
		return true
	})

	as.Equal(n, cnt)
}

func BenchmarkGid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Gid()
	}
}
