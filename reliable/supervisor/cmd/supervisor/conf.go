/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/backoff"
	"github.com/burningxflame/gx/reliable/supervisor"
)

type conf struct {
	Procs []supervisor.Proc
	Log   light.Conf
}

func readConf(pa string) (conf, error) {
	var cf conf
	var tmp confJson

	f, err := os.Open(pa)
	if err != nil {
		return cf, err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return conf{}, err
	}

	err = yaml.Unmarshal(bs, &tmp)
	if err != nil {
		return cf, err
	}

	for _, p := range tmp.Procs {
		cf.Procs = append(cf.Procs, supervisor.Proc{
			Tag:  p.Tag,
			Path: p.Path,
			Args: p.Args,
			Bf: backoff.Conf{
				Min:        time.Millisecond,
				Max:        p.Bf.Max.Duration,
				Unit:       p.Bf.Unit.Duration,
				Strategy:   p.Bf.Strategy.Strategy,
				ResetAfter: p.Bf.ResetAfter.Duration,
			},
		})
	}

	cf.Log = light.Conf{
		Level:         log.LevelInfo,
		BufSize:       int(tmp.Log.BufSize) << 10,
		FlushInterval: tmp.Log.FlushInterval.Duration,
		Rc: light.RotateConf{
			FilePath:   tmp.Log.FilePath,
			FileSize:   int64(tmp.Log.FileSize) << 20,
			NBak:       int(tmp.Log.NBak),
			Perm:       tmp.Log.Perm.FileMode,
			NoCompress: tmp.Log.NoCompress,
			Utc:        tmp.Log.Utc,
		}}
	return cf, nil
}

type confJson struct {
	Procs []struct {
		Tag  string
		Path string
		Args []string
		Bf   struct {
			Max        duration
			Unit       duration
			Strategy   strategy
			ResetAfter duration
		}
	}
	Log struct {
		FilePath      string
		FileSize      uint16
		NBak          byte
		Perm          perm
		NoCompress    bool
		Utc           bool
		BufSize       uint16
		FlushInterval duration
	}
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var tmp byte

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	d.Duration = time.Duration(tmp) * time.Second
	return nil
}

type strategy struct {
	backoff.Strategy
}

func (s *strategy) UnmarshalJSON(b []byte) error {
	var tmp string

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch tmp {
	case "l":
		s.Strategy = backoff.Linear
	case "e":
		s.Strategy = backoff.Exponent
	default:
		return fmt.Errorf("invalid strategy %v", tmp)
	}

	return nil
}

type perm struct {
	fs.FileMode
}

func (d *perm) UnmarshalJSON(b []byte) error {
	var tmp string

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	v, err := strconv.ParseUint(tmp, 8, 32)
	if err != nil {
		return err
	}

	d.FileMode = fs.FileMode(v)
	return nil
}
