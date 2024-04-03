/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package uuid

import (
	"strconv"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkNewSec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSec()
	}
}

func BenchmarkSecGenerator(b *testing.B) {
	for i := 10; i <= 10_000; i *= 10 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			benchmarkSecGenerator(b, i)
		})
	}
}

func benchmarkSecGenerator(b *testing.B, batch int) {
	gen := SecGenerator(batch)
	for i := 0; i < b.N; i++ {
		gen()
	}
}

func BenchmarkSecConGenerator(b *testing.B) {
	for i := 10; i <= 10_000; i *= 10 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			benchmarkSecConGenerator(b, i)
		})
	}
}

func benchmarkSecConGenerator(b *testing.B, batch int) {
	gen := SecConGenerator(batch)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gen()
		}
	})
}
