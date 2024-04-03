/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package autoreload

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/reliable/backoff"
)

func TestAutoReload(t *testing.T) {
	light.InitTestLog()

	tcs := []testCase{
		{"usual", "pXt", []string{"qThwh", "hHGdQsl"}, false, []string{"pXt", "qThwh", "hHGdQsl"}},
		{"noChange", "qBL", []string{"qBL", "qBL"}, false, []string{"qBL"}},
		{"rewatch", "ZPz", []string{"bVXzU", "lHZZMul"}, true, []string{"ZPz", "bVXzU", "lHZZMul"}},
		{"rewatchNoChange", "cqV", []string{"cqV", "cqV"}, true, []string{"cqV"}},
	}

	for _, tc := range tcs {
		t.Run(tc.tag, func(t *testing.T) {
			testAutoReload(t, tc)
		})
	}
}

type testCase struct {
	tag     string
	init    string
	changes []string
	remove  bool
	expect  []string
}

func testAutoReload(t *testing.T, tc testCase) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), tc.tag)
	err := os.WriteFile(pa, []byte(tc.init), 0600)
	as.Nil(err)

	var actual []string
	process := func(ctx context.Context, conf string) {
		process(ctx, conf, &actual)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		unit := time.Millisecond
		bf := backoff.Conf{
			Min:        unit,
			Max:        unit * 10,
			Unit:       unit * 2,
			Strategy:   backoff.Exponent,
			ResetAfter: unit * 10,
		}

		WithAutoReload(ctx, Conf[string]{
			Tag:     tc.tag,
			Path:    pa,
			Load:    load,
			Process: process,
			Bf:      bf,
		})
	}()

	go func() {
		defer wg.Done()

		if tc.remove {
			time.Sleep(time.Millisecond * 10)
			err := os.Remove(pa)
			as.Nil(err)
		}

		for _, txt := range tc.changes {
			time.Sleep(time.Millisecond * 10)
			err := os.WriteFile(pa, []byte(txt), 0600)
			as.Nil(err)
		}
	}()

	wg.Wait()
	as.Equal(tc.expect, actual)
}

func load(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func process(ctx context.Context, conf string, actual *[]string) {
	*actual = append(*actual, conf)

	<-ctx.Done()
}
