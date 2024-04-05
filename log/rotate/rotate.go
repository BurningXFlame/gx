/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// Log Rotator
package rotate

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Conf struct {
	// Fullpath of log file
	FilePath string
	// Max byte size of a log file. If a file exceeds this size, the file will be rotated. Default to 10MB.
	FileSize int64
	// Max number of old log files. Older files will be removed.
	NBak int
	// Permission of log file. Default to 0600.
	Perm fs.FileMode
	// If true, rotated log files will not be compressed. Otherwise, rotated log files will be compressed with gzip.
	NoCompress bool
	// If ture, rotated log files will be renamed based on UTC time. Local time otherwise.
	Utc bool

	gzPerm         fs.FileMode
	fmtBak         string
	patternBak     string
	patternBakOrGz string
}

func (c *Conf) adjust() error {
	path, err := filepath.Abs(c.FilePath)
	if err != nil {
		return err
	}
	c.FilePath = path

	err = os.MkdirAll(filepath.Dir(c.FilePath), c.Perm|0100)
	if err != nil {
		return err
	}

	if c.FileSize < 1 {
		c.FileSize = 10 << 20
	}

	if c.NBak < 0 {
		c.NBak = 0
	}

	const minPerm = 0600
	if c.Perm < minPerm {
		c.Perm = minPerm
	}

	c.gzPerm = c.Perm & 0444 // -w -x

	// name format of backup files
	ext := filepath.Ext(c.FilePath)
	if len(ext) == 0 {
		ext = ".log"
	}
	allButExt := c.FilePath[:len(c.FilePath)-len(ext)]
	c.fmtBak = fmt.Sprintf("%v-%%v%v", allButExt, ext)

	// glob patterns of backup files
	c.patternBak = fmt.Sprintf(c.fmtBak, "*")
	c.patternBakOrGz = c.patternBak + "*"

	return nil
}

const gzExt = ".gz"

type rotate struct {
	conf     Conf
	mu       sync.Mutex
	file     *os.File
	size     int64
	chRotate chan struct{}
	chDone   chan struct{}
}

// Create a log rotator
func New(conf Conf) (io.WriteCloser, error) {
	err := conf.adjust()
	if err != nil {
		return nil, err
	}

	r := &rotate{
		conf:     conf,
		chRotate: make(chan struct{}, 1),
		chDone:   make(chan struct{}, 1),
	}

	err = r.openLogFile()
	if err != nil {
		return nil, err
	}

	r.startHandleBaks()

	return r, nil
}

// Close the log rotator. Close files, finish handling backups, etc.
// Must be called before process exit.
func (r *rotate) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stopHandleBaks()

	return r.closeLogFile()
}

func (r *rotate) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.reCreateIf()
	if err != nil {
		return 0, err
	}

	if r.size+int64(len(p)) > r.conf.FileSize {
		err := r.rotate()
		if err != nil {
			return 0, err
		}
	}

	n, err := r.file.Write(p)
	r.size += int64(n)

	return n, err
}

func (r *rotate) openLogFile() error {
	file, err := os.OpenFile(r.conf.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, r.conf.Perm)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	r.file = file
	r.size = info.Size()
	return nil
}

func (r *rotate) closeLogFile() error {
	if r.file == nil {
		return nil
	}

	err := r.file.Close()
	r.file = nil
	r.size = 0
	return err
}

func (r *rotate) rotate() error {
	err := r.file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(r.conf.FilePath, r.genBakName())
	if err != nil {
		return err
	}

	r.notifyHandleBaks()

	return r.openLogFile()
}

func (r *rotate) genBakName() string {
	const format = "20060102150405.000"

	now := time.Now()
	if r.conf.Utc {
		now = now.UTC()
	}

	ts := now.Format(format)
	ts = strings.Replace(ts, ".", "", 1)

	return fmt.Sprintf(r.conf.fmtBak, ts)
}

func (r *rotate) startHandleBaks() {
	go func() {
		defer close(r.chDone)

		for range r.chRotate {
			_ = r.delOldBaks()
			_ = r.compressBaks()
		}
	}()
}

func (r *rotate) stopHandleBaks() {
	close(r.chRotate)
	<-r.chDone
}

func (r *rotate) notifyHandleBaks() {
	select {
	case r.chRotate <- struct{}{}:
	default:
	}
}

func (r *rotate) delOldBaks() error {
	baks, err := filepath.Glob(r.conf.patternBakOrGz)
	if err != nil {
		return err
	}

	if len(baks) <= r.conf.NBak {
		return nil
	}

	sort.Slice(baks, func(i, j int) bool {
		return baks[i] < baks[j]
	})

	var errs []error
	for _, file := range baks[:len(baks)-r.conf.NBak] {
		err := os.Remove(file)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

func (r *rotate) compressBaks() error {
	if r.conf.NoCompress {
		return nil
	}

	baks, err := filepath.Glob(r.conf.patternBak)
	if err != nil {
		return err
	}

	var errs []error
	for _, bak := range baks {
		err := compress(bak, r.conf.gzPerm)
		if err != nil {
			errs = append(errs, err)
		}

		err = os.Remove(bak)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

func compress(pa string, perm fs.FileMode) error {
	in, err := os.Open(pa)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(pa+gzExt, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	defer out.Close()

	zw := gzip.NewWriter(out)
	defer zw.Close()

	_, err = io.Copy(zw, in)
	return err
}
