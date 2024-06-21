/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package guard

import (
	"context"
	"time"

	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/backoff"
)

type Conf struct {
	// The func to be guarded.
	// Fn should return ASAP when ctx.Done channel is closed, which usually means an exit signal is sent.
	Fn func(ctx context.Context) error
	// Backoff strategy determines how long to wait between retries.
	Bf backoff.Conf
	// If true, re-run Fn even if it returns nil error.
	AlsoRetryOnSuccess bool
	// Used to tag log messages
	Tag string
	// A TagLogger used to log messages
	Log log.TagLogger
}

// Auto re-run a function until it succeeds (aka, returns nil error) or ctx.Done channel is closed.
// If AlsoRetryOnSuccess is true, auto re-run a function until ctx.Done channel is closed.
func WithGuard(ctx context.Context, cf Conf) {
	if cf.Log == nil {
		cf.Log = log.WithTag("")
	}
	lg := cf.Log.WithTag("guard " + cf.Tag)
	lg.Info("starting")

	bf := backoff.New(cf.Bf)

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			lg.Info("received exit signal, exiting")
			return

		case <-timer.C:
		}

		err := cf.Fn(ctx)
		if err == nil && !cf.AlsoRetryOnSuccess {
			lg.Info("completed")
			return
		}
		if err != nil && ctx.Err() != nil {
			lg.Info("received exit signal, exiting")
			return
		}

		dur := bf.Next()

		if err != nil {
			lg.Warn("re-run in %v because of error: %v", dur, err)
		} else {
			lg.Warn("re-run in %v", dur)
		}

		timer.Reset(dur)
	}
}
