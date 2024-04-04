/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package connpool

import (
	"errors"
	"net"
	"time"

	"github.com/burningxflame/gx/reliable/timeouts"
)

type Conf struct {
	Init int // Initial number of connections in the pool
	Cap  int // Max number of connections in the pool
	// Used to create a connection
	New func() (net.Conn, error)
	// Used to check whether a connection is still connected. Ping returns nil if still connected.
	Ping func(net.Conn) error
	// Timeout of New and Ping. Abandon the call to New or Ping after the timeout elapses.
	Timeout time.Duration
}

func (c *Conf) adjust() error {
	if c.Init < 0 {
		c.Init = 0
	}

	if c.Cap < 1 {
		c.Cap = 1
	}

	if c.Init > c.Cap || c.New == nil || c.Ping == nil {
		return errInvalidConf
	}

	if c.Timeout > 0 {
		c.New = timeouts.WithTimeoutO(c.Timeout, c.New)
		c.Ping = timeouts.WithTimeoutI(c.Timeout, c.Ping)
	}

	return nil
}

var errInvalidConf = errors.New("invaid conf")

// Concurrency-safe connection pool
type Pool struct {
	cf Conf
	ch chan net.Conn
}

// Create a connection pool
func New(cf Conf) (*Pool, error) {
	err := cf.adjust()
	if err != nil {
		return nil, err
	}

	p := &Pool{
		cf: cf,
		ch: make(chan net.Conn, cf.Cap),
	}

	go func() {
		for i := 0; i < cf.Init; i++ {
			conn, err := cf.New()
			if err != nil {
				continue
			}

			p.ch <- conn
		}
	}()

	return p, nil
}

// Get a connection from the pool.
// If the connection is not connected any more, will drop it and get another, until a connected connection is found.
// If no connection in the pool is connected, will create a new one.
func (p *Pool) Get() (net.Conn, error) {
	for {
		select {
		case conn := <-p.ch:
			if p.cf.Ping(conn) != nil {
				_ = conn.Close()
				continue
			}
			return conn, nil

		default:
			return p.cf.New()
		}
	}
}

// Put a connection back to the pool.
// Will close and drop the connection if the pool is full.
// Nil connection is ignored.
func (p *Pool) Put(conn net.Conn) {
	if conn == nil {
		return
	}

	select {
	case p.ch <- conn:
	default:
		_ = conn.Close()
	}

	return
}

// Close all connections in the pool
func (p *Pool) Close() {
	for {
		select {
		case conn := <-p.ch:
			_ = conn.Close()
		default:
			return
		}
	}
}
