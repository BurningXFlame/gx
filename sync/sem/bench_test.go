/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package sem

import (
	"context"
	"testing"
)

func BenchmarkAcquireRelease(b *testing.B) {
	s := New(b.N)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Acquire(ctx)
		s.Release()
	}
}
