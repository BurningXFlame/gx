package tcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnId(t *testing.T) {
	as := require.New(t)

	id := connId()
	id2 := connId()
	as.NotEqual(id, id2)
}

func TestConnIdCollision(t *testing.T) {
	as := require.New(t)

	n := 1_000_000
	used := make(map[string]struct{})

	for i := 0; i < n; i++ {
		id := connId()

		_, ok := used[id]
		as.False(ok)

		used[id] = struct{}{}
	}
}

func TestConnIdCtx(t *testing.T) {
	as := require.New(t)

	ctx := context.Background()
	_, ok := GetConnId(ctx)
	as.False(ok)

	id := connId()
	ctx = withConnId(ctx, id)
	id2, ok := GetConnId(ctx)
	as.True(ok)
	as.Equal(id, id2)
}
