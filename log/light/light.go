/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// Light is an all-in-one logger.
// Light Logger = Go Std Log Wrapper + Concurrent Buffer Writer (with Auto Flusher) + Log Rotator.
package light

import (
	"context"
	"fmt"
	"io"
	std "log"
	"os"
	"time"

	"github.com/burningxflame/gx/log/conbuf"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/log/rotate"
)

type RotateConf = rotate.Conf

type Conf struct {
	// Log level. Default to LevelError.
	Level log.Level
	// Log format flag. Refer to go std log. Default to LstdFlags | Lmicroseconds | Lmsgprefix.
	Format int
	// Buffer Size in bytes. Default to 1M.
	BufSize int
	// Auto-flush interval. Default to 5s.
	FlushInterval time.Duration
	// Log-rotating config
	Rc RotateConf
}

func (c *Conf) adjust() {
	if c.BufSize < 1 {
		c.BufSize = 1 << 20
	}

	if c.FlushInterval < time.Second {
		c.FlushInterval = time.Second * 5
	}
}

// Create a Light logger, and register it as the global logger.
func Init(conf Conf) error {
	conf.adjust()

	rw, err := rotate.New(conf.Rc)
	if err != nil {
		return err
	}

	bw := conbuf.NewWriter(rw, conf.BufSize)

	ctx, cancel := context.WithCancel(context.Background())
	fw := conbuf.WithAutoFlush(ctx, bw, conf.FlushInterval, nil)

	w := &writer{
		Writer: fw,
		closers: []func() error{
			func() error {
				cancel()
				return nil
			},
			bw.Flush,
			rw.Close,
		},
	}

	lg := newStdWrapper(w, conf.Format)
	return log.Set(lg, conf.Level)
}

func newStdWrapper(w io.WriteCloser, format int) log.Logger {
	const defFmt = std.LstdFlags | std.Lmicroseconds | std.Lmsgprefix

	if format == 0 {
		format = defFmt
	}

	return &stdWrapper{
		inner: std.New(w, "", format),
	}
}

type stdWrapper struct {
	inner *std.Logger
}

func (l *stdWrapper) Printf(format string, v ...any) {
	l.inner.Printf(format, v...)
}

func (l *stdWrapper) Close() error {
	return l.inner.Writer().(io.WriteCloser).Close()
}

type writer struct {
	io.Writer
	closers []func() error
}

func (w *writer) Close() error {
	var errs []error

	for _, fn := range w.closers {
		err := fn()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

// For test purpose, create a simple logger, and register it as the global logger.
// All log messages will be written to stdout, and log level is debug.
func InitTestLog() error {
	lg := newStdWrapper(&noClose{os.Stdout}, 0)
	return log.Set(lg, log.LevelDebug)
}

type noClose struct {
	io.Writer
}

func (w *noClose) Close() error {
	return nil
}
