/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package nanoid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	as := require.New(t)

	for i := 0; i < 10000; i++ {
		id, err := New()
		as.Nil(err)

		as.Equal(size, len(id))
		for _, c := range id {
			as.True(strings.ContainsRune(alphabet, c))
		}
	}
}

func TestNewCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{}, n)

	for i := 0; i < n; i++ {
		id, err := New()
		as.Nil(err)

		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}

func TestGenerator(t *testing.T) {
	as := require.New(t)

	batch := 100
	gen := Generator(batch)

	for i := 0; i < 10000; i++ {
		id, err := gen()
		as.Nil(err)

		as.Equal(size, len(id))
		for _, c := range id {
			as.True(strings.ContainsRune(alphabet, c))
		}
	}
}

func TestGeneratorCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{})
	gen := Generator(100)

	for i := 0; i < n; i++ {
		id, err := gen()
		as.Nil(err)

		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}

func TestConGenerator(t *testing.T) {
	as := require.New(t)

	batch := 100
	gen := ConGenerator(batch)

	for i := 0; i < 10000; i++ {
		id, err := gen()
		as.Nil(err)

		as.Equal(size, len(id))
		for _, c := range id {
			as.True(strings.ContainsRune(alphabet, c))
		}
	}
}

func TestConGeneratorCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	idSet := make(map[string]struct{})
	gen := ConGenerator(100)

	for i := 0; i < n; i++ {
		id, err := gen()
		as.Nil(err)

		idSet[id] = struct{}{}
	}

	as.Equal(n, len(idSet))
}
