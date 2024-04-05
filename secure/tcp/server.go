package tcp

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/timeouts"
	"github.com/burningxflame/gx/secure/conns"
	"github.com/burningxflame/gx/sync/sem"
)

// The Secure TCP Server has builtin abilities to defend against DDoS attacks, to defend against slow attacks, to close idle connections, etc.
type Server struct {
	// The address to listen
	Addr string
	// Connection handler is where you communicate with a client, i.e. receive/send data from/to a client.
	// Connection handler should return ASAP when ctx.Done channel is closed, which usually means an exit signal is sent.
	ConnHandler func(ctx context.Context, conn net.Conn) error
	// If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
	// Default to no timeout.
	ShutdownTimeout time.Duration
	// Used to limit max number of concurrent connections.
	// Default to no limit.
	ConnLimiter *sem.Sem
	// If no data is sent from a connection in the specified duration, close the connection.
	// Default to no timeout.
	IdleTimeout time.Duration
	// Used for TLS handshake. If not provided, no TLS handshake.
	TlsConfig *tls.Config
	// If TLS handshake does not finish in the specified duration, close the connection.
	// Default to no timeout.
	TlsHandshakeTimeout time.Duration
	// If true, the Context argument of ConnHandler contains the identity of the TLS peer.
	// Call GetTlsPeer(ctx) to get peer identity.
	// And of course the peer should send a certificate, i.e. TlsConfig.ClientAuth should be RequireAnyClientCert or RequireAndVerifyClientCert.
	// Peer identity is a set of Common Name and SAN DNS Names of certificate holder, i.e. cert.Subject.CommonName and cert.DNSNames.
	CtxTlsPeer bool
	// If true, the Context argument of ConnHandler contains the connection id.
	// Call GetConnId(ctx) to get connection id.
	CtxConnId bool
	// Used to tag log messages
	Tag string
	// A TagLogger used to log messages
	Log log.TagLogger

	wg sync.WaitGroup
}

// Start the Server
func (s *Server) Serve(ctx context.Context) error {
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
		conn, err := ln.Accept()
		if err != nil && errors.Is(err, net.ErrClosed) { // ln closed, i.e. exiting
			break
		}
		if err != nil {
			log.Error("error accepting incoming conn: %v", err)
			continue
		}

		s.handleConn(ctx, conn, log)
	}

	return timeouts.WithTimeout(s.ShutdownTimeout, func() error {
		s.wg.Wait()
		return nil
	})()
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn, log log.TagLogger) {
	id := connId()
	log = log.WithTag(id)
	log.Info("incoming conn [%v]", conn.RemoteAddr())

	if s.ConnLimiter != nil {
		err := s.ConnLimiter.Acquire(ctx)
		if err != nil { // ctx.Done channel closed, i.e. exiting
			_ = conn.Close()
			return
		}
	}

	s.wg.Add(1)
	go func() {
		defer func() {
			conn.Close()
			s.wg.Done()
			if s.ConnLimiter != nil {
				s.ConnLimiter.Release()
			}
		}()

		if s.IdleTimeout > 0 {
			c, err := conns.WithIdleTimeout(conn, s.IdleTimeout)
			if err != nil {
				log.Error("error enabling idle timeout: %v", err)
			} else {
				conn = c
			}
		}

		if s.TlsConfig != nil {
			_ctx, _conn, err := s.tlsHandshake(ctx, conn)
			if err != nil {
				log.Error("error TLS handshake: %v", err)
				return
			}

			conn = _conn
			ctx = _ctx
		}

		if s.CtxConnId {
			ctx = withConnId(ctx, id)
		}

		err := s.ConnHandler(ctx, conn)
		if err != nil && errors.Is(err, conns.ErrIdleTimeout) {
			log.Warn("%v", err)
			return
		}
		if err != nil {
			log.Error("%v", err)
			return
		}
		log.Info("done processing")
	}()
}

func (s *Server) tlsHandshake(ctx context.Context, conn net.Conn) (context.Context, net.Conn, error) {
	tlsConn := tls.Server(conn, s.TlsConfig)
	if s.TlsHandshakeTimeout > 0 {
		ctx, cancel := context.WithTimeout(ctx, s.TlsHandshakeTimeout)
		defer cancel()

		err := tlsConn.HandshakeContext(ctx)
		if err != nil {
			return nil, nil, err
		}
	} else {
		err := tlsConn.Handshake()
		if err != nil {
			return nil, nil, err
		}
	}

	if s.CtxTlsPeer {
		ctx = withTlsPeer(ctx, tlsConn)
	}

	return ctx, tlsConn, nil
}
