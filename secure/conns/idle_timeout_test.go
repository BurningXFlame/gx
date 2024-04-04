/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conns

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/id/uuid"
)

func TestIdleTimeout(t *testing.T) {
	as := require.New(t)

	dur := time.Second / 100

	tcs := []testcase{
		{0, dur, true},
		{dur, dur * 2, false},
		{dur * 2, dur, true},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			port := 1024 + rand.Intn(math.MaxInt8)
			addr := fmt.Sprintf("127.0.0.1:%v", port)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go server(ctx, as, addr, tc)

			time.Sleep(dur)

			var wg sync.WaitGroup

			for i := 0; i < 3; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					client(as, addr, tc)
				}()
			}

			wg.Wait()
		})
	}
}

type testcase struct {
	timeout time.Duration
	delay   time.Duration
	ok      bool
}

func server(ctx context.Context, as *require.Assertions, addr string, tc testcase) {
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

			conn, err = WithIdleTimeout(conn, tc.timeout)
			as.Nil(err)

			_, err = io.Copy(conn, conn)
			if tc.ok {
				as.Nil(err)
			} else {
				as.ErrorIs(err, ErrIdleTimeout)
			}
		}()
	}
}

func client(as *require.Assertions, addr string, tc testcase) {
	conn, err := net.Dial("tcp", addr)
	as.Nil(err)
	defer conn.Close()

	time.Sleep(tc.delay)

	err = comm(as, conn)
	if tc.ok {
		as.Nil(err)
	} else {
		as.Error(err)
	}
}

func comm(as *require.Assertions, conn net.Conn) error {
	msg := uuid.New()
	buf := make([]byte, len(msg))

	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}

	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}

	as.Equal(msg, string(buf))
	return nil
}
