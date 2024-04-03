/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package readiness

import (
	"context"
	"errors"
	"net"

	"github.com/burningxflame/gx/log/log"
)

type Server struct {
	Addr string // The address to listen
	Tag  string // Used to tag log messages. Defaults to "readiness".
	Log  log.TagLogger
}

// Start a TCP server for readiness check (aka, health check). Only for connectivity check. For security purpose, no sending data nor receiving data.
func (s *Server) Serve(ctx context.Context) error {
	if len(s.Tag) == 0 {
		s.Tag = "readiness"
	}

	if s.Log == nil {
		s.Log = log.WithTag("")
	}
	log := s.Log.WithTag(s.Tag)

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Info("listening at %v", ln.Addr())

	go func() {
		defer ln.Close()

		<-ctx.Done()
		log.Info("received exit signal, exiting")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) { // listener is closed
				return nil
			}

			log.Error("error accepting incoming conn: %v", err)
			continue
		}

		_ = conn.Close()
	}
}
