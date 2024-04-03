/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import (
	"testing"
)

func TestBufWriter(t *testing.T) {
	test(t, NewWriter)
}

func BenchmarkBufWriter(b *testing.B) {
	bench(b, NewWriter)
}
