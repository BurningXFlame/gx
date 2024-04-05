/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// Concurrent Buffer Writer
package conbuf

import (
	"bufio"
	"io"
	"sync"
)

type writer struct {
	w  io.Writer
	bw *bufio.Writer
	mu sync.Mutex
}

// Create a concurrent Buffer Writer
func NewWriter(w io.Writer, bufSize int) WriteFlusher {
	return &writer{
		w:  w,
		bw: bufio.NewWriterSize(w, bufSize),
	}
}

func (w *writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// make sure to write p as a whole
	if len(p) > w.bw.Available() {
		err := w.bw.Flush()
		if err != nil {
			return 0, err
		}
	}

	n, err := w.bw.Write(p)
	if err != nil {
		// According to the spec,
		// if an error occurs writing to a bufio.Writer, no more data will be accepted,
		// and all subsequent writes and Flush will return the error.
		// So, reset bufio.Writer on error.
		w.bw.Reset(w.w)
	}

	return n, err
}

func (w *writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.bw.Flush()
}
