/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package nanoid

import (
	"strconv"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkGenerator(b *testing.B) {
	for i := 10; i <= 10_000; i *= 10 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			benchmarkGenerator(b, i)
		})
	}
}

func benchmarkGenerator(b *testing.B, batch int) {
	gen := Generator(batch)
	for i := 0; i < b.N; i++ {
		gen()
	}
}

func BenchmarkConGenerator(b *testing.B) {
	for i := 10; i <= 10_000; i *= 10 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			benchmarkConGenerator(b, i)
		})
	}
}

func benchmarkConGenerator(b *testing.B, batch int) {
	gen := ConGenerator(batch)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gen()
		}
	})
}
