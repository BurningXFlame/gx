/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package set

import "testing"

func BenchmarkAdd(b *testing.B) {
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkDelete(b *testing.B) {
	s := New[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Delete(i)
	}
}
