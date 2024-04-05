/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAutoFlush(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const interval = time.Millisecond * 100

	var w bytes.Buffer
	bw := NewWriter(&w, msgLen*10)
	fw := WithAutoFlush(ctx, bw, interval, nil)

	n, err := fw.Write(msg)
	as.Nil(err)
	as.Equal(len(msg), n)

	time.Sleep(interval * 3 / 2)
	as.Equal(msg, w.Bytes())
}
