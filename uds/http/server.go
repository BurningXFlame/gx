/*
GX (https://github.com/BurningXFlame/gx).
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

// HTTP server over UDS
type Server struct {
	Std     http.Server // http.Server in std lib
	Tag     string      // Used to tag log messages
	UdsAddr string      // The UDS address to listen
	Perm    fs.FileMode // File permission of the UdsAddr
	// If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
	ShutdownTimeout time.Duration
	Log             log.TagLogger
}

// Start an HTTP server over UDS
func (s *Server) Serve(ctx context.Context) error {
	if s.Log == nil {
		s.Log = log.WithTag("")
	}
	log := s.Log.WithTag(s.Tag)

	// Clean up in case the process was killed forcibly last time.
	_ = syscall.Unlink(s.UdsAddr)

	ln, err := net.Listen("unix", s.UdsAddr)
	if err != nil {
		log.Error("error listening: %v", err)
		return err
	}
	defer ln.Close()

	log.Info("listening at " + s.UdsAddr)

	if s.Perm > 0 {
		err := os.Chmod(s.UdsAddr, s.Perm)
		if err != nil {
			log.Warn("error chmoding uds file: %v", err)
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
		log.Info("received exit signal, exiting")

		ctx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
		defer cancel()

		err := s.Std.Shutdown(ctx)
		if err != nil {
			log.Warn("error shutting down: %v", err)
		}

		return nil
	}
}
