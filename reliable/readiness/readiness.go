/*
GX (github.com/burningxflame/gx).
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

// TCP Server for readiness check (aka, health check).
// Only for connectivity check.
// For security purpose, no sending data nor receiving data.
type Server struct {
	// The address to listen
	Addr string
	// Used to tag log messages. Default to "readiness".
	Tag string
	// A TagLogger used to log messages
	Log log.TagLogger
}

// Start the Server.
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
		<-ctx.Done()
		ln.Close()
		log.Info("received exit signal, exiting")
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := ln.Accept()
		if err != nil && errors.Is(err, net.ErrClosed) { // ln closed, i.e. exiting
			return nil
		}
		if err != nil {
			log.Error("error accepting incoming conn: %v", err)
			continue
		}

		_ = conn.Close()
	}
}
