/*
GX (https://github.com/BurningXFlame/gx).
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
	Tag string // Used to tag log messages
	// The func to be guarded. Fn must be canceled when ctx.Done channel is closed.
	Fn func(ctx context.Context) error
	// Backoff strategy determines how long to wait between retries.
	Bf backoff.Conf
	// If true, re-run Fn even if it returns nil error.
	AlsoRetryOnSuccess bool
	Log                log.TagLogger
}

// Auto re-run a function until it succeeds (aka, returns nil error) or ctx.Done channel is closed.
// If AlsoRetryOnSuccess is true, auto re-run a function until ctx.Done channel is closed.
func WithGuard(ctx context.Context, cf Conf) {
	if cf.Log == nil {
		cf.Log = log.WithTag("")
	}
	log := cf.Log.WithTag("guard " + cf.Tag)
	log.Info("starting")

	bf := backoff.New(cf.Bf)

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("received exit signal, exiting")
			return

		case <-timer.C:
		}

		err := cf.Fn(ctx)
		if err == nil && !cf.AlsoRetryOnSuccess {
			log.Info("completed")
			return
		}

		dur := bf.Next()

		if err != nil {
			log.Warn("re-run in %v because of error: %v", dur, err)
		} else {
			log.Warn("re-run in %v", dur)
		}

		timer.Reset(dur)
	}
}
