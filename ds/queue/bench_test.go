/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package queue

import "testing"

func BenchmarkEnq(b *testing.B) {
	q := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Enq(i)
	}
}

func BenchmarkDeq(b *testing.B) {
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Enq(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Deq()
	}
}
