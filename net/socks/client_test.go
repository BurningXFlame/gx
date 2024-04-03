/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import (
	"io"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/armon/go-socks5"
	"github.com/burningxflame/gx/id/uuid"
	"github.com/stretchr/testify/require"
)

const (
	socksAddr      = "127.0.0.1:1080"
	destAddr       = "127.0.0.1:6666"
	destAddrV6     = "[::1]:6666"
	destAddrDomain = "localhost:6666"
)

func TestClientHandShake(t *testing.T) {
	as := require.New(t)

	go socksServer(as)
	go echoServer(as, destAddr)
	go echoServer(as, destAddrV6)
	time.Sleep(time.Second / 10)

	tcs := []string{
		destAddr,
		destAddrV6,
		destAddrDomain,
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			conn, err := net.Dial("tcp", socksAddr)
			as.Nil(err)
			defer conn.Close()

			err = ClientHandshake(conn, tc)
			as.Nil(err)

			err = comm(as, conn)
			as.Nil(err)
		})
	}
}

func socksServer(as *require.Assertions) {
	server, err := socks5.New(&socks5.Config{})
	as.Nil(err)

	err = server.ListenAndServe("tcp", socksAddr)
	as.Nil(err)
}

func echoServer(as *require.Assertions, addr string) {
	ln, err := net.Listen("tcp", addr)
	as.Nil(err)

	for {
		conn, err := ln.Accept()
		as.Nil(err)

		go func() {
			defer conn.Close()

			_, err := io.Copy(conn, conn)
			as.Nil(err)
		}()
	}
}

func comm(as *require.Assertions, conn net.Conn) error {
	for i := 0; i < 3; i++ {
		msg := uuid.New()
		buf := make([]byte, len(msg))

		_, err := conn.Write([]byte(msg))
		as.Nil(err)

		_, err = io.ReadFull(conn, buf)
		as.Nil(err)

		as.Equal(msg, string(buf))
	}

	return nil
}
