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
)

func BenchmarkGlsGet(b *testing.B) {
	Put(&key, rand.Int())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get(&key)
	}
}

func BenchmarkGlsCostA(b *testing.B) {
	Put(&key, rand.Int())

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		Go(func() {
			wg.Done()
		})
	}

	wg.Wait()
}

// In contrast to BenchmarkGlsCostA
func BenchmarkGlsCostB(b *testing.B) {
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}

	wg.Wait()
}
