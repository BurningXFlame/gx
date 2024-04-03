/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package sem

import (
	"context"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	sizeHint = 2
	key      = "zAEmh"
)

func TestKSem(t *testing.T) {
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	const key2 = key + "2"

	ks := NewKSem[string](ca, sizeHint)
	s := ks.Get(key)
	s2 := ks.Get(key2)

	for i := 0; i < ca; i++ {
		as.Nil(s.Acquire(ctx))
		as.Nil(s2.Acquire(ctx))
	}

	as.False(s.TryAcquire())
	as.False(s2.TryAcquire())

	s.Release()
	s2.Release()

	as.Nil(s.Acquire(ctx))
	as.Nil(s2.Acquire(ctx))
}

func TestKSemShrink(t *testing.T) {
	as := require.New(t)

	ks := NewKSem[string](ca, sizeHint)

	for i := 1; i <= sizeHint; i++ {
		s := ks.Get(key + strconv.Itoa(i))
		as.True(s.TryAcquire())
		as.Equal(i, ks.size)
	}

	ks.Get(key + "1").Release()

	ks.Get(key + strconv.Itoa(sizeHint+1))
	as.Equal(sizeHint, ks.size)
}

func TestKSemShrinkNone(t *testing.T) {
	as := require.New(t)

	ks := NewKSem[string](ca, sizeHint)

	for i := 1; i <= sizeHint; i++ {
		s := ks.Get(key + strconv.Itoa(i))
		as.True(s.TryAcquire())
		as.Equal(i, ks.size)
	}

	ks.Get(key + strconv.Itoa(sizeHint+1))
	as.Equal(sizeHint+1, ks.size)
}

func TestKSemCapLow(t *testing.T) {
	as := require.New(t)

	ks := NewKSem[string](0, sizeHint)
	as.Equal(defCap, ks.ca)

	s := ks.Get(key)
	as.Equal(defCap, s.Available())
}

func TestKSemSizeHintLow(t *testing.T) {
	as := require.New(t)

	ks := NewKSem[string](ca, 0)
	as.Equal(math.MaxInt, ks.sizeHint)
}
