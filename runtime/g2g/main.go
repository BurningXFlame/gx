/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: g2g <path>")
		os.Exit(1)
	}

	light.InitTestLog()

	err := filepath.WalkDir(os.Args[1], walk)
	if err != nil {
		log.Error("%v", err)
	}

	log.Info("summary: number of go files: %v, # replaced: %v, # failed: %v", scanned, replaced, failed)
}

var (
	scanned  int
	replaced int
	failed   int
)

func walk(pa string, _ fs.DirEntry, err error) error {
	log := log.WithTag(pa)

	if err != nil {
		log.Warn("skipped because of err: %v", err)
		return fs.SkipDir
	}

	if filepath.Ext(pa) != ".go" {
		return nil
	}

	err = g2g.process(pa)
	defer g2g.clear()
	scanned++
	if err != nil {
		failed++
		log.Warn("error processing: %v", err)
		return nil
	}

	if g2g.replaced {
		replaced++
	}

	log.Info("processed")
	return nil
}

var g2g = &g2G{}
