/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import (
	"context"
	"time"
)

type flusher struct {
	WriteFlusher
	exit  <-chan struct{}
	chErr chan<- error
}

// Wrap a WriteFlusher and return another WriteFlusher which auto flushes the wrapped WriteFlusher at certain intervals until ctx.Done channel is closed.
func WithAutoFlush(
	ctx context.Context,
	w WriteFlusher, // the WriteFlusher to be wrapped
	interval time.Duration, // the flush interval
	// Specify a buffer chan if you want to receive background flush errors if any. Leave it nil otherwise.
	chErr chan<- error,
) WriteFlusher {
	f := &flusher{
		WriteFlusher: w,
		exit:         ctx.Done(),
		chErr:        chErr,
	}

	f.startFlush(interval)

	return f
}

func (f *flusher) startFlush(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-f.exit:
				return

			case <-ticker.C:
				err := f.WriteFlusher.Flush()
				if err != nil && f.chErr != nil {
					// If chErr is full, ignore the error.
					select {
					case f.chErr <- err:
					default:
					}
				}
			}
		}
	}()
}
