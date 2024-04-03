/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import (
	"io"
	"sync"
	"testing"
)

func BenchmarkNoBuf(b *testing.B) {
	bench(b, newNoBuf)
}

func newNoBuf(w io.Writer, _ int) WriteFlusher {
	return &nobuf{w: w}
}

type nobuf struct {
	w  io.Writer
	mu sync.Mutex
}

func (w *nobuf) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.w.Write(p)
}

func (w *nobuf) Flush() error {
	return nil
}
