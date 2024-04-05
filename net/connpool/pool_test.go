/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package connpool

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		cf Conf
		ok bool
	}{
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk}, true},

		{Conf{Init: 0, Cap: 16, New: newFakeConn, Ping: pingOk}, true},
		{Conf{Init: -1, Cap: 16, New: newFakeConn, Ping: pingOk}, true},
		{Conf{Init: 0, Cap: 1, New: newFakeConn, Ping: pingOk}, true},
		{Conf{Init: 0, Cap: 0, New: newFakeConn, Ping: pingOk}, true},

		{Conf{Init: 2, Cap: 1, New: newFakeConn, Ping: pingOk}, false},
		{Conf{Init: 8, Cap: 16, New: nil, Ping: pingOk}, false},
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: nil}, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			p, err := New(tc.cf)
			if !tc.ok {
				as.ErrorIs(err, errInvalidConf)
				return
			}

			as.Nil(err)
			time.Sleep(time.Second / 10)
			tc.cf.adjust()
			as.Equal(tc.cf.Init, len(p.ch))
		})
	}
}

func TestGet(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		cf  Conf
		len int
	}{
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk}, 7},
		{Conf{Init: 0, Cap: 16, New: newFakeConn, Ping: pingOk}, 0},
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingErr}, 0},
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingHalfErr()}, 6},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			p, err := New(tc.cf)
			as.Nil(err)

			time.Sleep(time.Second / 10)

			conn, err := p.Get()
			as.NotNil(conn)
			as.Nil(err)

			as.Equal(tc.len, len(p.ch))
		})
	}
}

func TestPut(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		cf   Conf
		len  int
		drop bool
	}{
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk}, 9, false},
		{Conf{Init: 16, Cap: 16, New: newFakeConn, Ping: pingOk}, 16, true},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			p, err := New(tc.cf)
			as.NotNil(p)
			as.Nil(err)

			time.Sleep(time.Second / 10)

			conn := &fakeConn{}
			p.Put(conn)
			as.Equal(tc.len, len(p.ch))
			if tc.drop {
				as.True(conn.closed)
			}

			p.Put(nil)
			as.Equal(tc.len, len(p.ch))
		})
	}
}

func TestClose(t *testing.T) {
	as := require.New(t)

	p, err := New(Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk})
	as.Nil(err)

	time.Sleep(time.Second / 10)

	conn := &fakeConn{}
	p.Put(conn)

	p.Close()
	as.Equal(0, len(p.ch))
	as.True(conn.closed)
}

func TestNewTimeout(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100

	tcs := []struct {
		cf  Conf
		len int
	}{
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk, Timeout: timeout}, 8},
		{Conf{Init: 8, Cap: 16, New: newConnTimeout(timeout * 2), Ping: pingOk, Timeout: timeout}, 0},
		{Conf{Init: 8, Cap: 16, New: newConnHalfTimeout(timeout * 2), Ping: pingOk, Timeout: timeout}, 4},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			p, err := New(tc.cf)
			as.Nil(err)

			time.Sleep(time.Second / 10)
			as.Equal(tc.len, len(p.ch))
		})
	}
}

func TestPingTimeout(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100

	tcs := []struct {
		cf  Conf
		len int
	}{
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingOk, Timeout: timeout}, 7},
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingTimeout(timeout * 2), Timeout: timeout}, 0},
		{Conf{Init: 8, Cap: 16, New: newFakeConn, Ping: pingHalfTimeout(timeout * 2), Timeout: timeout}, 6},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			p, err := New(tc.cf)
			as.Nil(err)

			time.Sleep(time.Second / 10)

			conn, err := p.Get()
			as.NotNil(conn)
			as.Nil(err)

			as.Equal(tc.len, len(p.ch))
		})
	}
}

func newFakeConn() (net.Conn, error) {
	return &fakeConn{}, nil
}

type fakeConn struct {
	b      bytes.Buffer
	closed bool
}

func (c *fakeConn) Read(b []byte) (n int, err error) {
	return c.b.Read(b)
}

func (c *fakeConn) Write(b []byte) (n int, err error) {
	return c.b.Write(b)
}

func (c *fakeConn) Close() error {
	c.closed = true
	return nil
}

func (c *fakeConn) LocalAddr() net.Addr                { return nil }
func (c *fakeConn) RemoteAddr() net.Addr               { return nil }
func (c *fakeConn) SetDeadline(_ time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(_ time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(_ time.Time) error { return nil }

func pingOk(net.Conn) error {
	return nil
}

func pingErr(net.Conn) error {
	return errDummy
}

var errDummy = errors.New("dummy")

func pingHalfErr() func(net.Conn) error {
	var ok bool

	return func(net.Conn) error {
		defer func() {
			ok = !ok
		}()

		if ok {
			return nil
		}

		return errDummy
	}
}

func newConnTimeout(timeout time.Duration) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		time.Sleep(timeout)
		return &fakeConn{}, nil
	}
}

func newConnHalfTimeout(timeout time.Duration) func() (net.Conn, error) {
	var ok bool

	return func() (net.Conn, error) {
		if ok {
			ok = !ok
			return &fakeConn{}, nil
		}

		ok = !ok
		time.Sleep(timeout)
		return &fakeConn{}, nil
	}
}

func pingTimeout(timeout time.Duration) func(net.Conn) error {
	return func(net.Conn) error {
		time.Sleep(timeout)
		return nil
	}
}

func pingHalfTimeout(timeout time.Duration) func(net.Conn) error {
	var ok bool

	return func(net.Conn) error {
		if ok {
			ok = !ok
			return nil
		}

		ok = !ok
		time.Sleep(timeout)
		return nil
	}
}
