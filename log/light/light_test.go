/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package light

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/log/log"
)

const filename = "dummy.log"

func TestLevel(t *testing.T) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), filename)

	err := Init(Conf{
		Level:         log.LevelInfo,
		BufSize:       1024,
		FlushInterval: time.Second,
		Rc: RotateConf{
			FilePath: pa,
			NBak:     2,
		},
	})
	as.Nil(err)

	log.Error("some error: %v", "dummy error")
	log.Warn("some warning")
	log.Info("some info")
	log.Debug("some debug")
	log.Trace("some trace")

	log.Close()

	_content, err := os.ReadFile(pa)
	as.Nil(err)
	content := string(_content)

	as.Contains(content, "ERROR some error: dummy error")
	as.Contains(content, "WARN  some warning")
	as.Contains(content, "INFO  some info")
	as.NotContains(content, "some debug")
	as.NotContains(content, "some trace")
}

func TestWithTag(t *testing.T) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), filename)

	err := Init(Conf{
		Level:         log.LevelInfo,
		BufSize:       1024,
		FlushInterval: time.Second,
		Rc: RotateConf{
			FilePath: pa,
			NBak:     2,
		},
	})
	as.Nil(err)

	tl := log.WithTag("tag1")

	tl.Error("some error: %v", "dummy error")
	tl.Warn("some warning")
	tl.Info("some info")
	tl.Debug("some debug")
	tl.Trace("some trace")

	log.Close()

	_content, err := os.ReadFile(pa)
	as.Nil(err)
	content := string(_content)

	as.Contains(content, "ERROR [tag1] some error: dummy error")
	as.Contains(content, "WARN  [tag1] some warning")
	as.Contains(content, "INFO  [tag1] some info")
	as.NotContains(content, "some debug")
	as.NotContains(content, "some trace")
}

func TestWithTagChain(t *testing.T) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), filename)

	err := Init(Conf{
		Level:         log.LevelInfo,
		BufSize:       1024,
		FlushInterval: time.Second,
		Rc: RotateConf{
			FilePath: pa,
			NBak:     2,
		},
	})
	as.Nil(err)

	tl := log.WithTag("tag1").WithTag("tag2")

	tl.Error("some error: %v", "dummy error")
	tl.Warn("some warning")
	tl.Info("some info")
	tl.Debug("some debug")
	tl.Trace("some trace")

	log.Close()

	_content, err := os.ReadFile(pa)
	as.Nil(err)
	content := string(_content)

	as.Contains(content, "ERROR [tag1] [tag2] some error: dummy error")
	as.Contains(content, "WARN  [tag1] [tag2] some warning")
	as.Contains(content, "INFO  [tag1] [tag2] some info")
	as.NotContains(content, "some debug")
	as.NotContains(content, "some trace")
}

func TestRotate(t *testing.T) {
	as := require.New(t)

	dir := t.TempDir()
	pa := filepath.Join(dir, filename)

	err := Init(Conf{
		Level:         log.LevelInfo,
		BufSize:       1024,
		FlushInterval: time.Second,
		Rc: RotateConf{
			FilePath: pa,
			NBak:     2,
		},
	})
	as.Nil(err)

	const (
		msg    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
		expect = 20
	)

	for i := 0; i < expect; i++ {
		log.Info(msg)
	}
	err = log.Close()
	as.Nil(err)

	// current log file
	_content, err := os.ReadFile(pa)
	as.Nil(err)
	content := string(_content)
	actual := strings.Count(content, msg)

	// backup log file
	pattern := filepath.Join(dir, "*.gz")
	matches, err := filepath.Glob(pattern)
	as.Nil(err)

	for _, pa := range matches {
		var out bytes.Buffer

		func() {
			in, err := os.Open(pa)
			as.Nil(err)
			defer in.Close()

			zr, err := gzip.NewReader(in)
			as.Nil(err)
			defer zr.Close()

			_, err = io.Copy(&out, zr)
			as.Nil(err)
		}()

		actual += strings.Count(out.String(), msg)
	}

	as.Equal(expect, actual)
}
