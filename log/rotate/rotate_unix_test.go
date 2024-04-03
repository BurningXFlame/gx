/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package rotate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReCreate(t *testing.T) {
	as := require.New(t)

	dir := t.TempDir()
	path := filepath.Join(dir, filename)

	ca := 100

	r, err := New(Conf{
		FilePath: path,
		FileSize: int64(size * ca),
		NBak:     2,
	})
	as.Nil(err)
	defer r.Close()

	for i := 0; i < ca*2/3; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	err = os.Remove(path)
	as.Nil(err)

	for i := 0; i < ca/3; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	checkTxt(as, path, ca/3)
}
