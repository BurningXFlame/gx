/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package backoff

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBackoff(t *testing.T) {
	as := require.New(t)

	const unit = time.Second

	tcs := []struct {
		min        int
		max        int
		unit       int
		strategy   Strategy
		resetAfter int
		expect     []int
	}{
		{1, 10, 1, Linear, 10, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{3, 10, 2, Linear, 10, []int{3, 5, 7, 9, 10}},
		{1, 10, 1, Exponent, 10, []int{1, 2, 4, 8, 10}},
		{3, 10, 2, Exponent, 10, []int{3, 5, 9, 10}},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			b := New(Conf{
				Min:        unit * time.Duration(tc.min),
				Max:        unit * time.Duration(tc.max),
				Unit:       unit * time.Duration(tc.unit),
				Strategy:   tc.strategy,
				ResetAfter: unit * time.Duration(tc.resetAfter),
			})

			for i := 0; i < len(tc.expect); i++ {
				as.Equal(b.Next(), unit*time.Duration(tc.expect[i]))
			}

			as.Equal(b.Next(), unit*time.Duration(tc.max))
		})
	}
}

func TestReset(t *testing.T) {
	as := require.New(t)

	const unit = time.Millisecond
	b := New(Conf{
		Min:        unit,
		Max:        unit * 30,
		Unit:       unit,
		Strategy:   Exponent,
		ResetAfter: unit * 30,
	})
	min := unit
	b.Next()
	as.Greater(b.Next(), min)

	time.Sleep(b.Next() + b.conf.ResetAfter)
	as.Equal(min, b.Next())
}
