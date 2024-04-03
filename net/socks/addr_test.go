/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import (
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseAddr(t *testing.T) {
	as := require.New(t)

	tcs := []struct {
		in  string
		out addr
	}{
		{
			"192.168.1.81:3618",
			addr{
				host: &ipv4{[4]byte{192, 168, 1, 81}},
				port: 3618,
			},
		},
		{
			"[fc00::2:12]:3588",
			addr{
				host: &ipv6{[16]byte{0xfc, 0, 13: 0x02, 15: 0x12}},
				port: 3588,
			},
		},
		{
			"3k0nyirdgb5ge.f382g.gr:8768",
			addr{
				host: &domain{"3k0nyirdgb5ge.f382g.gr"},
				port: 8768,
			},
		},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			ad, err := parseAddr(tc.in)
			as.Nil(err)
			as.Equal(tc.out, ad)
		})
	}
}

func TestParseAddrInvalid(t *testing.T) {
	as := require.New(t)

	tcs := []string{
		"192.168.1.81",                                 // no port
		"fc00::2:12:3588",                              // invalid format
		"3k0nyirdgb5ge.f382g.gr:65536",                 // port out of range
		"192.168.1.81:abc",                             // invalid port
		strings.Repeat("x", domainMaxSize+1) + ":3618", // domain name oversize
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			_, err := parseAddr(tc)
			as.Error(err)
		})
	}
}

func TestAddrSendRecv(t *testing.T) {
	as := require.New(t)

	tcs := []string{
		"192.168.1.81:3618",
		"[fc00::2:12]:3588",
		"3k0nyirdgb5ge.f382g.gr:8768",
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			ad, err := parseAddr(tc)
			as.Nil(err)

			conn := &fakeConn{}

			buf := getBuf()
			defer putBuf(buf)

			err = ad.send(conn, *buf)
			as.Nil(err)

			var ad2 addr

			buf2 := getBuf()
			defer putBuf(buf2)

			err = (&ad2).recv(conn, *buf2)
			as.Nil(err)

			as.Equal(ad, ad2)
		})
	}
}

func TestDiscardAddr(t *testing.T) {
	as := require.New(t)

	tcs := []string{
		"192.168.1.81:3618",
		"[fc00::2:12]:3588",
		"3k0nyirdgb5ge.f382g.gr:8768",
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			ad, err := parseAddr(tc)
			as.Nil(err)

			conn := &fakeConn{}

			buf := getBuf()
			defer putBuf(buf)

			err = ad.send(conn, *buf)
			as.Nil(err)

			as.Greater(conn.b.Len(), 0)

			buf2 := getBuf()
			defer putBuf(buf2)

			discardAddr(conn, *buf2)
			as.Nil(err)

			as.Equal(conn.b.Len(), 0)
		})
	}
}

type fakeConn struct {
	b bytes.Buffer
}

func (c *fakeConn) Read(b []byte) (n int, err error) {
	return c.b.Read(b)
}

func (c *fakeConn) Write(b []byte) (n int, err error) {
	return c.b.Write(b)
}

func (c *fakeConn) Close() error                       { return nil }
func (c *fakeConn) LocalAddr() net.Addr                { return nil }
func (c *fakeConn) RemoteAddr() net.Addr               { return nil }
func (c *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (c *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *fakeConn) SetWriteDeadline(t time.Time) error { return nil }
