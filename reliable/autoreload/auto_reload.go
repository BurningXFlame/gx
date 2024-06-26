/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package autoreload

import (
	"context"
	"errors"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/backoff"
	"github.com/burningxflame/gx/reliable/guard"
)

type Conf[C comparable] struct {
	// Path of the config file
	Path string
	// Used to load config file on file write event. The returned C is the loaded config.
	// If C is the same as the last, it's ignored.
	Load func(path string) (C, error)
	// Process is the func to be reloaded. C is the config loaded by Load.
	// Process should return ASAP when ctx.Done channel is closed.
	Process func(ctx context.Context, c C)
	// Backoff strategy determines how long to wait between retries.
	Bf backoff.Conf
	// Used to tag log messages
	Tag string
	// A TagLogger used to log messages
	Log log.TagLogger
}

// Watch a config file, and re-run a func on config changes.
// Will retry watching if the conf file is removed.
func WithAutoReload[C comparable](ctx context.Context, cf Conf[C]) {
	if cf.Log == nil {
		cf.Log = log.WithTag("")
	}
	lg := cf.Log.WithTag("autoReload " + cf.Tag)
	lg.Info("starting")

	ch := make(chan C, 1)
	done := make(chan struct{})

	go func() {
		defer close(done)

		guard.WithGuard(ctx, guard.Conf{
			Tag: "watch " + cf.Tag,
			Fn: func(ctx context.Context) error {
				return watch(ctx, cf.Path, cf.Load, ch, lg)
			},
			Bf:                 cf.Bf,
			AlsoRetryOnSuccess: true,
		})
	}()

	reload(ctx, cf.Process, ch, lg)

	<-done
}

func watch[C comparable](
	ctx context.Context,
	path string,
	load func(string) (C, error),
	ch chan<- C,
	lg log.TagLogger,
) error {
	val, err := load(path)
	if err != nil {
		return err
	}

	ch <- val

	_load := func(path string) {
		val, err := load(path)
		if err != nil {
			lg.Error("error loading conf: %v", err)
			return
		}

		ch <- val
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(path)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				_load(path)
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				return errRemove
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			return err
		}
	}
}

var errRemove = errors.New("file was removed")

func reload[C comparable](
	ctx context.Context,
	fn func(context.Context, C),
	ch <-chan C,
	lg log.TagLogger,
) {
	var val C
	var ctxFn context.Context
	var cancelFn context.CancelFunc = func() {}
	var wgFn sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			lg.Info("received exit signal, exiting")
			cancelFn()
			wgFn.Wait()
			return

		case v := <-ch:
			if v == val {
				lg.Info("ignore because of no change")
				continue
			}

			val = v

			lg.Info("reloading on change")

			cancelFn()
			wgFn.Wait()

			ctxFn, cancelFn = context.WithCancel(ctx)
			wgFn.Add(1)
			go func() {
				defer wgFn.Done()
				fn(ctxFn, val)
			}()
		}
	}
}
