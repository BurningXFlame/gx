/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package timeouts

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithTimeoutIO(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100

	fn := func(timeout time.Duration) func(int) (string, error) {
		return func(i int) (string, error) {
			time.Sleep(timeout)
			return strconv.Itoa(i), nil
		}
	}

	tcs := []struct {
		fn      func(int) (string, error)
		timeout time.Duration
		ok      bool
	}{
		{fn(0), 0, true},
		{fn(0), timeout, true},
		{fn(timeout), 0, true},
		{fn(timeout * 2), timeout, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			fn := WithTimeoutIO(tc.timeout, tc.fn)
			v, err := fn(68)
			if tc.ok {
				as.Nil(err)
				as.Equal(v, "68")
			} else {
				as.ErrorIs(err, ErrTimeout)
			}
		})
	}
}

func TestWithTimeoutI(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100

	fn := func(timeout time.Duration) func(int) error {
		return func(int) error {
			time.Sleep(timeout)
			return nil
		}
	}

	tcs := []struct {
		fn      func(int) error
		timeout time.Duration
		ok      bool
	}{
		{fn(0), 0, true},
		{fn(0), timeout, true},
		{fn(timeout), 0, true},
		{fn(timeout * 2), timeout, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			fn := WithTimeoutI(tc.timeout, tc.fn)
			err := fn(68)
			if tc.ok {
				as.Nil(err)
			} else {
				as.ErrorIs(err, ErrTimeout)
			}
		})
	}
}

func TestWithTimeoutO(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100
	msg := "rqQFux"

	fn := func(timeout time.Duration) func() (string, error) {
		return func() (string, error) {
			time.Sleep(timeout)
			return msg, nil
		}
	}

	tcs := []struct {
		fn      func() (string, error)
		timeout time.Duration
		ok      bool
	}{
		{fn(0), 0, true},
		{fn(0), timeout, true},
		{fn(timeout), 0, true},
		{fn(timeout * 2), timeout, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			fn := WithTimeoutO(tc.timeout, tc.fn)
			v, err := fn()
			if tc.ok {
				as.Nil(err)
				as.Equal(v, msg)
			} else {
				as.ErrorIs(err, ErrTimeout)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	as := require.New(t)

	timeout := time.Second / 100

	fn := func(timeout time.Duration) func() error {
		return func() error {
			time.Sleep(timeout)
			return nil
		}
	}

	tcs := []struct {
		fn      func() error
		timeout time.Duration
		ok      bool
	}{
		{fn(0), 0, true},
		{fn(0), timeout, true},
		{fn(timeout), 0, true},
		{fn(timeout * 2), timeout, false},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			fn := WithTimeout(tc.timeout, tc.fn)
			err := fn()
			if tc.ok {
				as.Nil(err)
			} else {
				as.ErrorIs(err, ErrTimeout)
			}
		})
	}
}
