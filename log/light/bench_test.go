/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package light

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/burningxflame/gx/log/log"
)

const (
	msg           = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
	bufSize       = 2 << 20
	flushInterval = 5 * time.Second
)

var msgLong = strings.Repeat(msg, 5)

func BenchmarkShort(b *testing.B) {
	b.Run("short", func(b *testing.B) {
		benchmark(b, msg)
	})
	b.Run("long", func(b *testing.B) {
		benchmark(b, msgLong)
	})
}

func benchmark(b *testing.B, msg string) {
	pa := filepath.Join(b.TempDir(), "light.log")

	err := Init(Conf{
		Rc: RotateConf{
			FilePath: pa,
			FileSize: 20 << 20,
			NBak:     10,
			Perm:     0600,
		},
		BufSize:       bufSize,
		FlushInterval: flushInterval,
		Level:         log.LevelInfo,
	})
	if err != nil {
		b.FailNow()
	}
	defer log.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info(msg)
		}
	})

	b.StopTimer()
}
