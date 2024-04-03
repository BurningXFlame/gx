/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package main

import (
	"testing"
	"time"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/log/rotate"
	"github.com/burningxflame/gx/reliable/backoff"
	"github.com/burningxflame/gx/reliable/supervisor"
	"github.com/stretchr/testify/require"
)

func TestConf(t *testing.T) {
	as := require.New(t)

	actual, err := readConf("testdata/sample_conf.yaml")
	as.Nil(err)

	expect := conf{
		Procs: []supervisor.Proc{
			{
				Tag:  "a",
				Path: "/bin/sh",
				Args: []string{"-c", "date >> /tmp/xyz/a.txt"},
				Bf: backoff.Conf{
					Min:        time.Millisecond,
					Max:        10 * time.Second,
					Unit:       time.Second,
					Strategy:   backoff.Linear,
					ResetAfter: 10 * time.Second,
				},
			},
			{
				Tag:  "b",
				Path: "/bin/sh",
				Args: []string{"-c", "date +%T >> /tmp/xyz/b.txt"},
				Bf: backoff.Conf{
					Min:        time.Millisecond,
					Max:        10 * time.Second,
					Unit:       time.Second,
					Strategy:   backoff.Exponent,
					ResetAfter: 10 * time.Second,
				},
			},
		},
		Log: light.Conf{
			Level:         log.LevelInfo,
			BufSize:       1 << 20,
			FlushInterval: 5 * time.Second,
			Rc: rotate.Conf{
				FilePath:   "/tmp/xyz/supervisor.log",
				FileSize:   10 << 20,
				NBak:       2,
				Perm:       0600,
				NoCompress: false,
				Utc:        false,
			},
		},
	}

	as.Equal(expect, actual)
}
