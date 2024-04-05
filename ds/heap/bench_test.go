/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package heap

import "testing"

var less Less[int] = func(a, b int) bool {
	return a < b
}

func BenchmarkPush(b *testing.B) {
	q := New(less)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	q := New(less)
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
