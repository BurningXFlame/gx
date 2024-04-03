/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package guard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/reliable/backoff"
)

func TestUntilSuccess(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	WithGuard(ctx, Conf{
		Tag: "dummy",
		Fn:  failUntil(3),
		Bf:  bf,
	})

	as.Nil(ctx.Err())
}

var bf = backoff.Conf{
	Min:        time.Millisecond,
	Max:        time.Millisecond * 100,
	Unit:       time.Millisecond * 10,
	Strategy:   backoff.Linear,
	ResetAfter: time.Millisecond * 100,
}

func failUntil(times int) func(ctx context.Context) error {
	cnt := 0

	return func(ctx context.Context) error {
		cnt++

		if cnt < times {
			return errDummy
		}

		return nil
	}
}

var errDummy = errors.New("dummy")

func TestNeverSuccess(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	WithGuard(ctx, Conf{
		Tag: "dummy",
		Fn: func(_ context.Context) error {
			return errDummy
		},
		Bf: bf,
	})

	as.Error(ctx.Err())
}

func TestRetryOnSuccess(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	WithGuard(ctx, Conf{
		Tag: "dummy",
		Fn: func(_ context.Context) error {
			return nil
		},
		Bf:                 bf,
		AlsoRetryOnSuccess: true,
	})

	as.Error(ctx.Err())
}
