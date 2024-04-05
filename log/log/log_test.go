/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package log

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	as := require.New(t)

	lg := _new("x")
	err := Set(lg, LevelInfo)
	as.Nil(err)
	the := theConf.Load().(conf)
	as.Equal(lg, the.logger)
	as.Equal(LevelInfo, the.level)

	lg = _new("y")
	err = Set(lg, LevelDebug)
	as.Nil(err)
	the = theConf.Load().(conf)
	as.Equal(lg, the.logger)
	as.Equal(LevelDebug, the.level)
}

func TestNilLogger(t *testing.T) {
	as := require.New(t)

	err := Set(nil, LevelInfo)
	as.ErrorIs(err, errNilLogger)
}

func TestInvalidLevel(t *testing.T) {
	as := require.New(t)

	err := Set(_new("x"), 5)
	as.ErrorIs(err, errInvalidLevel)
}

func TestLevel(t *testing.T) {
	as := require.New(t)

	lg := _new("x")
	err := Set(lg, LevelInfo)
	as.Nil(err)
	defer Close()

	Error("some error: %v", "dummy error")
	as.Equal("ERROR some error: dummy error", lg.inner.String())
	lg.clear()

	Warn("some warning")
	as.Equal("WARN  some warning", lg.inner.String())
	lg.clear()

	Info("some info")
	as.Equal("INFO  some info", lg.inner.String())
	lg.clear()

	Debug("some debug")
	as.Equal(0, lg.inner.Len())

	Trace("some trace")
	as.Equal(0, lg.inner.Len())
}

func TestWithTag(t *testing.T) {
	as := require.New(t)

	lg := _new("x")
	err := Set(lg, LevelInfo)
	as.Nil(err)
	defer Close()

	tl := WithTag("tag1")

	tl.Error("some error: %v", "dummy error")
	as.Equal("ERROR [tag1] some error: dummy error", lg.inner.String())
	lg.clear()

	tl.Warn("some warning")
	as.Equal("WARN  [tag1] some warning", lg.inner.String())
	lg.clear()

	tl.Info("some info")
	as.Equal("INFO  [tag1] some info", lg.inner.String())
	lg.clear()

	tl.Debug("some debug")
	as.Equal(0, lg.inner.Len())

	tl.Trace("some trace")
	as.Equal(0, lg.inner.Len())
}

func TestWithTagChain(t *testing.T) {
	as := require.New(t)

	lg := _new("x")
	err := Set(lg, LevelInfo)
	as.Nil(err)
	defer Close()

	tl := WithTag("tag1").WithTag("tag2")

	tl.Error("some error: %v", "dummy error")
	as.Equal("ERROR [tag1] [tag2] some error: dummy error", lg.inner.String())
	lg.clear()

	tl.Warn("some warning")
	as.Equal("WARN  [tag1] [tag2] some warning", lg.inner.String())
	lg.clear()

	tl.Info("some info")
	as.Equal("INFO  [tag1] [tag2] some info", lg.inner.String())
	lg.clear()

	tl.Debug("some debug")
	as.Equal(0, lg.inner.Len())

	tl.Trace("some trace")
	as.Equal(0, lg.inner.Len())
}

func TestClose(t *testing.T) {
	as := require.New(t)

	lg := _new("x")
	err := Set(lg, LevelInfo)
	defer Close()
	as.Nil(err)

	msg := "some msg"
	_, _ = lg.inner.WriteString(msg)
	as.Equal(len(msg), lg.inner.Len())

	Close()
	as.Equal(0, lg.inner.Len())
}

type dummy struct {
	name  string
	inner bytes.Buffer
}

func _new(name string) *dummy {
	return &dummy{name: name}
}

func (l *dummy) Printf(format string, v ...any) {
	_, _ = l.inner.WriteString(fmt.Sprintf(format, v...))
}

func (l *dummy) Close() error {
	l.clear()
	return nil
}

func (l *dummy) clear() {
	_, _ = l.inner.WriteTo(io.Discard)
}
