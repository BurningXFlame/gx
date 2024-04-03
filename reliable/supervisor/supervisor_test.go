/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package supervisor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/backoff"
)

func TestSupervisor(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), "a")
	pa2 := filepath.Join(t.TempDir(), "b")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	Supervisor(ctx,
		Proc{
			Tag:  "a",
			Path: "/bin/sh",
			Args: []string{"-c", "date >> " + pa},
			Bf: backoff.Conf{
				Min:        time.Millisecond,
				Max:        time.Millisecond * 100,
				Unit:       time.Millisecond * 10,
				Strategy:   backoff.Linear,
				ResetAfter: time.Millisecond * 100,
			},
		},
		Proc{
			Tag:  "b",
			Path: "/bin/sh",
			Args: []string{"-c", "date +%T >> " + pa2},
			Bf: backoff.Conf{
				Min:        time.Millisecond,
				Max:        time.Millisecond * 100,
				Unit:       time.Millisecond * 10,
				Strategy:   backoff.Linear,
				ResetAfter: time.Millisecond * 100,
			},
		},
	)

	_content, err := os.ReadFile(pa)
	as.Nil(err)
	content := string(_content)
	log.Debug("a: %v", content)
	as.Greater(strings.Count(content, "\n"), 1)

	_content, err = os.ReadFile(pa2)
	as.Nil(err)
	content = string(_content)
	log.Debug("b: %v", content)
	as.Greater(strings.Count(content, "\n"), 1)
}
