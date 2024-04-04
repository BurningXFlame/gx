package conns

import (
	"errors"
	"net"
	"os"
	"time"
)

var ErrIdleTimeout = errors.New("idle timeout")

type connIdleTimeout struct {
	net.Conn
	timeout time.Duration
}

func WithIdleTimeout(conn net.Conn, timeout time.Duration) (net.Conn, error) {
	if timeout <= 0 {
		return conn, nil
	}

	err := conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return conn, err
	}

	return &connIdleTimeout{
		conn,
		timeout,
	}, nil
}

func (c *connIdleTimeout) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			err = ErrIdleTimeout
		}
		return n, err
	}

	err = c.SetReadDeadline(time.Now().Add(c.timeout))
	return n, err
}
