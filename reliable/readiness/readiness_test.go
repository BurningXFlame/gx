/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package readiness

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/log/light"
)

func TestReadiness(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const addr = "127.0.0.1:5063"

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		s := &Server{Addr: addr}
		s.Serve(ctx)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		defer wg.Done()
		defer cancel()

		for i := 0; i < 1000; i++ {
			conn, err := net.Dial("tcp", addr)
			as.Nil(err)
			conn.Close()
		}
	}()

	wg.Wait()
}
