/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/backoff"
	"github.com/burningxflame/gx/reliable/guard"
	"github.com/burningxflame/gx/reliable/supervisor"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		stdlog.Fatal("usage: supervisor <config.json>")
	}

	c, err := readConf(args[0])
	if err != nil {
		stdlog.Fatalf("error reading config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	initLog(ctx, c.Log)
	defer log.Close()

	err = supervisor.Supervisor(ctx, c.Procs...)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func initLog(ctx context.Context, conf light.Conf) {
	err := light.Init(conf)
	if err == nil {
		return
	}

	stdlog.Printf("error initing log: %v. Retry initing in the background.", err)

	go guard.WithGuard(ctx, guard.Conf{
		Tag: "log",
		Fn: func(_ context.Context) error {
			return light.Init(conf)
		},
		Bf: backoff.Default(),
	})
}
