/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package fakeconn

import (
	"fmt"
	"io"
	"net"
	"time"
)

func New() net.Conn {
	r, w := io.Pipe()
	return &fakeConn{
		r: r,
		w: w,
	}
}

type fakeConn struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (c *fakeConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *fakeConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

func (c *fakeConn) Close() error {
	err := c.w.Close()
	err2 := c.r.Close()
	return fmt.Errorf("%w, %w", err, err2)
}

func (c *fakeConn) LocalAddr() net.Addr                { return nil }
func (c *fakeConn) RemoteAddr() net.Addr               { return nil }
func (c *fakeConn) SetDeadline(_ time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(_ time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(_ time.Time) error { return nil }
