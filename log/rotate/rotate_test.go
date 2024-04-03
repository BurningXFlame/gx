/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package rotate

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	msg  = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-\n")
	size = len(msg)
)

const filename = "dummy.log"

func TestNewFile(t *testing.T) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), filename)

	r, err := New(Conf{
		FilePath: pa,
		NBak:     2,
	})
	as.Nil(err)
	defer r.Close()

	n, err := r.Write(msg)
	as.Nil(err)
	as.Equal(size, n)

	checkTxt(as, pa, 1)
}

func checkTxt(as *require.Assertions, pa string, n int) {
	content, err := os.ReadFile(pa)
	as.Nil(err)

	expect := bytes.Repeat(msg, n)
	as.Equal(expect, content)
}

func TestFileExist(t *testing.T) {
	as := require.New(t)

	pa := filepath.Join(t.TempDir(), filename)
	err := os.WriteFile(pa, msg, 0600)
	as.Nil(err)

	r, err := New(Conf{
		FilePath: pa,
		NBak:     2,
	})
	as.Nil(err)
	defer r.Close()

	n, err := r.Write(msg)
	as.Nil(err)
	as.Equal(size, n)

	checkTxt(as, pa, 2)
}

func TestRotate(t *testing.T) {
	as := require.New(t)

	dir := t.TempDir()
	pa := filepath.Join(dir, filename)

	ca := 100
	nBaks := 2

	r, err := New(Conf{
		FilePath: pa,
		FileSize: int64(size * ca),
		NBak:     nBaks,
	})
	as.Nil(err)

	for i := 0; i < ca*nBaks*2+1; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	err = r.Close()
	as.Nil(err)

	checkTxt(as, pa, 1)
	checkBaks(as, dir, nBaks, ca, false)
}

func checkBaks(as *require.Assertions, dir string, nBaks int, n int, noCompress bool) {
	pattern := filepath.Join(dir, "*-*")
	pas, err := filepath.Glob(pattern)
	as.Nil(err)

	as.Equal(nBaks, len(pas))

	for _, pa := range pas {
		if noCompress {
			checkTxt(as, pa, n)
		} else {
			checkGz(as, pa, n)
		}
	}
}

func checkGz(as *require.Assertions, pa string, n int) {
	in, err := os.Open(pa)
	as.Nil(err)
	defer in.Close()

	zr, err := gzip.NewReader(in)
	as.Nil(err)
	defer zr.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, zr)
	as.Nil(err)

	expect := bytes.Repeat(msg, n)
	as.Equal(expect, out.Bytes())
}

func TestNoCompress(t *testing.T) {
	as := require.New(t)

	dir := t.TempDir()
	pa := filepath.Join(dir, filename)

	ca := 100
	nBaks := 2

	r, err := New(Conf{
		FilePath:   pa,
		FileSize:   int64(size * ca),
		NBak:       nBaks,
		NoCompress: true,
	})
	as.Nil(err)

	for i := 0; i < ca*nBaks*2+1; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	err = r.Close()
	as.Nil(err)

	checkTxt(as, pa, 1)
	checkBaks(as, dir, nBaks, ca, true)
}

func TestPerm(t *testing.T) {
	tcs := []permTc{
		{0, 0600, 0400},
		{0660, 0660, 0440},
		{0777, 0777, 0444},
	}

	for ti, tc := range tcs {
		t.Run(strconv.Itoa(ti), func(t *testing.T) {
			testPerm(t, tc)
		})
	}
}

type permTc struct {
	perm fs.FileMode

	expect    fs.FileMode
	expectBak fs.FileMode
}

func testPerm(t *testing.T, tc permTc) {
	as := require.New(t)

	dir := t.TempDir()
	pa := filepath.Join(dir, filename)

	ca := 100
	nBaks := 2

	r, err := New(Conf{
		FilePath: pa,
		FileSize: int64(size * ca),
		NBak:     nBaks,
	})
	as.Nil(err)

	for i := 0; i < ca*nBaks*2+1; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	err = r.Close()
	as.Nil(err)

	checkPerm(as, pa, 0600)

	pattern := filepath.Join(dir, "*-*")
	pas, err := filepath.Glob(pattern)
	as.Nil(err)
	as.Equal(nBaks, len(pas))
	for _, pa := range pas {
		checkPerm(as, pa, 0400)
	}
}

func checkPerm(as *require.Assertions, pa string, expected fs.FileMode) {
	info, err := os.Stat(pa)
	as.Nil(err)
	as.Equal(expected, info.Mode().Perm())
}

func TestNoBak(t *testing.T) {
	as := require.New(t)

	dir := t.TempDir()
	pa := filepath.Join(dir, filename)

	ca := 100
	nBaks := 2

	r, err := New(Conf{
		FilePath: pa,
		FileSize: int64(size * ca),
	})
	as.Nil(err)

	for i := 0; i < ca*nBaks*2+1; i++ {
		n, err := r.Write(msg)
		as.Nil(err)
		as.Equal(size, n)
	}

	err = r.Close()
	as.Nil(err)

	checkTxt(as, pa, 1)
	checkBaks(as, dir, 0, ca, false)
}
