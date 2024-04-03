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
	"time"

	"github.com/stretchr/testify/require"
)

const ca = 3
const timeout = time.Second / 10

func TestAcquireAvailable(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s := New(ca)

	for i := 0; i < ca-1; i++ {
		as.Nil(s.Acquire(ctx))
	}

	as.Nil(s.Acquire(ctx))
}

func TestAcquireUnavailable(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s := New(ca)

	for i := 0; i < ca; i++ {
		as.Nil(s.Acquire(ctx))
	}

	as.Error(s.Acquire(ctx))
}

func TestTryAcquire(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s := New(ca)

	for i := 0; i < ca-1; i++ {
		as.Nil(s.Acquire(ctx))
	}

	as.True(s.TryAcquire())
	as.False(s.TryAcquire())
}

func TestRelease(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s := New(ca)

	for i := 0; i < ca; i++ {
		as.Nil(s.Acquire(ctx))
	}

	as.False(s.TryAcquire())

	s.Release()

	as.Nil(s.Acquire(ctx))
}

func TestAvailable(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s := New(ca)

	for i := 0; i < ca; i++ {
		as.Nil(s.Acquire(ctx))
		as.Equal(ca-1-i, s.Available())
	}

	s.Release()
	as.Equal(1, s.Available())

	as.Nil(s.Acquire(ctx))
	as.Equal(0, s.Available())
}

func TestCapLow(t *testing.T) {
	as := require.New(t)

	s := New(0)
	as.Equal(defCap, s.Available())
}
