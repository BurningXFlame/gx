/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package ringbuf

import "testing"

func BenchmarkPushBack(b *testing.B) {
	r := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.PushBack(i)
	}
}

func BenchmarkPushFront(b *testing.B) {
	r := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.PushFront(i)
	}
}

func BenchmarkPopFront(b *testing.B) {
	r := New[int]()
	for i := 0; i < b.N; i++ {
		r.PushBack(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.PopFront()
	}
}

func BenchmarkPopBack(b *testing.B) {
	r := New[int]()
	for i := 0; i < b.N; i++ {
		r.PushFront(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.PopBack()
	}
}
