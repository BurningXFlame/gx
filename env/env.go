/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package env

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func Require(key string) string {
	va := os.Getenv(key)
	if len(va) == 0 {
		log.Printf("no env var %v, exiting", key)
		os.Exit(1)
	}

	return va
}

func Get(key string, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}

	return v
}

func GetInt(key string, min int) int {
	v := os.Getenv(key)
	if len(v) == 0 {
		return min
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return min
	}

	if n < min {
		return min
	}

	return n
}

func GetDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}

	dur, err := time.ParseDuration(v)
	if err != nil {
		return def
	}

	return dur
}

func GetPort(key string, def uint16) uint16 {
	v := os.Getenv(key)
	if len(v) == 0 {
		return def
	}

	n, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		return def
	}

	// can't be a privileged port
	if n < 1024 {
		return def
	}

	return uint16(n)
}

func JoinHostPort(host string, port uint16) string {
	return net.JoinHostPort(host, strconv.Itoa(int(port)))
}
