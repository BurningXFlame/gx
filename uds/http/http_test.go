/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package http

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/id/uuid"
	"github.com/burningxflame/gx/log/light"
)

const (
	tag     = "echoServer"
	udsName = tag + ".sock"
)

func TestServerClient(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	udsAddr := filepath.Join(t.TempDir(), udsName)

	chServe := make(chan error, 1)

	go func() {
		s := &Server{
			Std: http.Server{
				Handler: http.HandlerFunc(handleEcho),
			},
			Tag:     tag,
			UdsAddr: udsAddr,
		}
		chServe <- s.Serve(ctx)
	}()

	time.Sleep(time.Millisecond * 10)

	req(as, udsAddr)

	cancel()
	as.Nil(<-chServe)
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func req(as *require.Assertions, udsAddr string) {
	c := NewClient(udsAddr)
	c.Timeout = time.Second * 3

	for i := 0; i < 3; i++ {
		msg := uuid.New()

		resp, err := c.Post("http://unix/echo", "raw", strings.NewReader(msg))
		as.Nil(err)

		body, err := io.ReadAll(resp.Body)
		as.Nil(err)
		resp.Body.Close()

		as.Equal(msg, string(body))
	}
}

func TestServerNormalExit(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	udsAddr := filepath.Join(t.TempDir(), udsName)

	s := &Server{
		Std: http.Server{
			Handler: http.HandlerFunc(handleEcho),
		},
		Tag:     tag,
		UdsAddr: udsAddr,
	}
	err := s.Serve(ctx)
	as.Nil(err)
}

func TestServerError(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	dir := fmt.Sprintf("%v-%v", t.TempDir(), time.Now().UnixNano())
	udsAddr := filepath.Join(dir, udsName)

	s := &Server{
		Std: http.Server{
			Handler: http.HandlerFunc(handleEcho),
		},
		Tag:     tag,
		UdsAddr: udsAddr,
	}
	err := s.Serve(ctx)
	as.Error(err)
}

func TestServerCleanup(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
	defer cancel()

	udsAddr := filepath.Join(t.TempDir(), udsName)
	f, err := os.Create(udsAddr)
	as.Nil(err)
	f.Close()

	s := &Server{
		Std: http.Server{
			Handler: http.HandlerFunc(handleEcho),
		},
		Tag:     tag,
		UdsAddr: udsAddr,
	}
	err = s.Serve(ctx)
	as.Nil(err)
}

func TestServerPerm(t *testing.T) {
	light.InitTestLog()
	as := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	udsAddr := filepath.Join(t.TempDir(), udsName)
	perm := fs.FileMode(0600)

	chServe := make(chan error, 1)

	go func() {
		s := &Server{
			Std: http.Server{
				Handler: http.HandlerFunc(handleEcho),
			},
			Tag:     tag,
			UdsAddr: udsAddr,
			Perm:    perm,
		}
		chServe <- s.Serve(ctx)
	}()

	time.Sleep(time.Millisecond * 10)

	info, err := os.Stat(udsAddr)
	as.Nil(err)
	as.Equal(perm, info.Mode().Perm())

	cancel()
	as.Nil(<-chServe)
}
