/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import (
	"context"
	"errors"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go socksServer(ctx, as)
	go server(ctx, as, destAddr)
	go server(ctx, as, destAddrV6)
	time.Sleep(time.Second / 10)

	tcs := []struct {
		addr string
	}{
		{destAddr},
		{destAddrV6},
		{destAddrDomain},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			conn, err := net.Dial("tcp", socksAddr)
			as.Nil(err)
			defer conn.Close()

			err = ClientHandshake(conn, tc.addr)
			as.Nil(err)

			comm(as, conn)
		})
	}
}

func socksServer(_ context.Context, as *require.Assertions) {
	server, err := socks5.New(&socks5.Config{})
	as.Nil(err)

	err = server.ListenAndServe("tcp", socksAddr)
	as.Nil(err)
}

func server(ctx context.Context, as *require.Assertions, addr string) {
	ln, err := net.Listen("tcp", addr)
	as.Nil(err)
	defer ln.Close()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil && errors.Is(err, net.ErrClosed) {
			return
		}
		as.Nil(err)

		go func() {
			defer conn.Close()

			_, err := io.Copy(conn, conn)
			as.Nil(err)
		}()
	}
}

func comm(as *require.Assertions, conn net.Conn) {
	for i := 0; i < 3; i++ {
		msg := uuid.New()
		buf := make([]byte, len(msg))

		_, err := conn.Write([]byte(msg))
		as.Nil(err)

		_, err = io.ReadFull(conn, buf)
		as.Nil(err)

		as.Equal(msg, string(buf))
	}
}
