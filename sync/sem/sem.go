/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package sem

import (
	"context"
)

// Semaphore is commonly used for limiting max concurrency, e.g. limiting max number of concurrent connections.
type Sem struct {
	ch chan struct{}
}

// Create a semaphore.
// The ca specifies the capacity of the semaphore.
func New(ca int) *Sem {
	if ca < 1 {
		ca = defCap
	}

	return &Sem{
		ch: make(chan struct{}, ca),
	}
}

const defCap = 1

// Acquire a permit from the semaphore.
// If none is available, block until one is available or ctx.Done channel is closed.
func (s *Sem) Acquire(ctx context.Context) error {
	select {
	case s.ch <- s0:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var s0 = struct{}{}

// Try to acquire a permit from the semaphore.
// Return true if available, false otherwise.
func (s *Sem) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release a permit to the semaphore.
func (s *Sem) Release() {
	select {
	case <-s.ch:
	default:
	}
}

// Return the number of available permits.
func (s *Sem) Available() int {
	return cap(s.ch) - len(s.ch)
}
