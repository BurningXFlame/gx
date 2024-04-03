/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type newer func(w io.Writer, bufSize int) WriteFlusher

var (
	msg    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-\n")
	msgLen = len(msg)
)

func test(t *testing.T, fn newer) {
	as := require.New(t)

	const n = 100

	var w bytes.Buffer
	bw := fn(&w, msgLen*n/3)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			n, err := bw.Write(msg)
			as.Nil(err)
			as.Equal(msgLen, n)
		}()
	}

	wg.Wait()
	bw.Flush()
	as.Equal(bytes.Repeat(msg, n), w.Bytes())
}

func bench(b *testing.B, fn newer) {
	const bufSize = 1 << 20

	f, err := os.CreateTemp(b.TempDir(), "")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	w := fn(f, bufSize)
	defer w.Flush()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := w.Write(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.StopTimer()
}
