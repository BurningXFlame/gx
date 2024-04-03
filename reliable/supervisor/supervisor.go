/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package supervisor

import (
	"context"
	"errors"
	"os/exec"
	"sync"

	"github.com/burningxflame/gx/reliable/backoff"
	"github.com/burningxflame/gx/reliable/guard"
)

type Proc struct {
	Tag  string       // Used to tag log messages
	Path string       // Path of the command to run
	Args []string     // Args of the command
	Bf   backoff.Conf // Backoff strategy determines how long to wait between retries
}

// Start and guard processes until ctx.Done channel is closed.
func Supervisor(ctx context.Context, procs ...Proc) error {
	if len(procs) < 1 {
		return errNoProc
	}

	for _, proc := range procs {
		if len(proc.Path) < 1 {
			return errEmptyCmd
		}
	}

	var wg sync.WaitGroup

	for _, proc := range procs {
		proc := proc

		wg.Add(1)
		go func() {
			defer wg.Done()

			guard.WithGuard(ctx, guard.Conf{
				Tag: proc.Tag,
				Fn: func(ctx context.Context) error {
					return startChild(ctx, proc.Path, proc.Args)
				},
				Bf:                 proc.Bf,
				AlsoRetryOnSuccess: true,
			})
		}()
	}

	wg.Wait()
	return nil
}

var (
	errNoProc   = errors.New("no proc")
	errEmptyCmd = errors.New("empty command")
)

func startChild(ctx context.Context, program string, args []string) error {
	return exec.CommandContext(ctx, program, args...).Run()
}
