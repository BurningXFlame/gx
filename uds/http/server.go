/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package http

import (
	"context"
	"io/fs"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/burningxflame/gx/log/log"
)

// HTTP Server over UDS
type Server struct {
	// http.Server in std lib
	Std http.Server
	// The UDS address to listen
	UdsAddr string
	// File permission of the UdsAddr
	Perm fs.FileMode
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

	// Clean up in case the process was killed forcibly last time.
	_ = syscall.Unlink(s.UdsAddr)

	ln, err := net.Listen("unix", s.UdsAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	lg.Info("listening at " + s.UdsAddr)

	if s.Perm > 0 {
		err := os.Chmod(s.UdsAddr, s.Perm)
		if err != nil {
			lg.Warn("error chmoding uds file: %v", err)
		}
	}

	chServe := make(chan error, 1)
	go func() {
		chServe <- s.Std.Serve(ln)
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
