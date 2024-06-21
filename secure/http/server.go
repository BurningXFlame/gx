/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package http

import (
	"context"
	"net/http"
	"time"

	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/sync/sem"
)

// Secure HTTP Server
type Server struct {
	// http.Server in std lib
	Std http.Server
	// Used to limit max number of concurrent requests.
	// Default to no limit.
	Limiter *sem.Sem
	// If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
	ShutdownTimeout time.Duration
	// Used to tag log messages
	Tag string
	// A TagLogger used to log messages
	Log log.TagLogger
}

// Start the Server
func (s *Server) Serve(ctx context.Context) error {
	if s.Log == nil {
		s.Log = log.WithTag("")
	}
	lg := s.Log.WithTag(s.Tag)

	if s.Limiter != nil {
		origHandler := s.Std.Handler
		defer func() {
			s.Std.Handler = origHandler
		}()

		s.Std.Handler = LimitHandler(ctx, s.Limiter, origHandler)
	}

	lg.Info("listening at %v", s.Std.Addr)

	chServe := make(chan error, 1)
	go func() {
		chServe <- s.Std.ListenAndServeTLS("", "")
	}()

	select {
	case err := <-chServe:
		return err

	case <-ctx.Done():
		lg.Info("received exit signal, exiting")

		ctx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
		defer cancel()

		err := s.Std.Shutdown(ctx)
		if err != nil {
			lg.Warn("error shutting down: %v", err)
		}

		return nil
	}
}

func LimitHandler(ctx context.Context, limiter *sem.Sem, handler http.Handler) http.Handler {
	if limiter == nil {
		return handler
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		err := limiter.Acquire(ctx)
		if err != nil { // ctx.Done channel closed, i.e. exiting
			return
		}
		defer limiter.Release()

		handler.ServeHTTP(rw, r)
	})
}
