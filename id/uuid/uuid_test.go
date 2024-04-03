/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package uuid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	as := require.New(t)

	for i := 0; i < 1000; i++ {
		id := New()
		as.Equal(strSize, len(id))
		as.True(valid(id))
	}
}

var (
	alphabet = "0123456789abcedf"
	strSize  = size * 2
)

func valid(id string) bool {
	for _, c := range id {
		if !strings.ContainsRune(alphabet, c) {
			return false
		}
	}

	return id[12] == '4' && strings.ContainsRune("89ab", rune(id[16]))
}

func TestNewCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{}, n)

	for i := 0; i < n; i++ {
		id := New()
		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}

func TestNewSec(t *testing.T) {
	as := require.New(t)

	for i := 0; i < 1000; i++ {
		id, err := NewSec()
		as.Nil(err)
		as.Equal(strSize, len(id))
		as.True(valid(id))
	}
}

func TestNewSecCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{}, n)

	for i := 0; i < n; i++ {
		id, err := NewSec()
		as.Nil(err)
		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}

func TestSecGenerator(t *testing.T) {
	as := require.New(t)

	batch := 100
	gen := SecGenerator(batch)

	for i := 0; i < 1000; i++ {
		id, err := gen()
		as.Nil(err)
		as.Equal(strSize, len(id))
		as.True(valid(id))
	}
}

func TestSecGeneratorCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{})
	gen := SecGenerator(100)

	for i := 0; i < n; i++ {
		id, err := gen()
		as.Nil(err)
		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}

func TestSecConGenerator(t *testing.T) {
	as := require.New(t)

	batch := 100
	gen := SecConGenerator(batch)

	for i := 0; i < 1000; i++ {
		id, err := gen()
		as.Nil(err)
		as.Equal(strSize, len(id))
		as.True(valid(id))
	}
}

func TestSecConGeneratorCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{})
	gen := SecConGenerator(100)

	for i := 0; i < n; i++ {
		id, err := gen()
		as.Nil(err)
		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}
