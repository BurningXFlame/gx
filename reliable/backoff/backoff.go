/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package backoff

import "time"

type Conf struct {
	// Min delay
	Min time.Duration
	// Max delay
	Max time.Duration
	// Unit of increment
	Unit time.Duration
	// Strategy of increment. Linear or Exponent.
	Strategy Strategy
	// If a retry lasts longer than ResetAfter, the next delay will be reset to Min.
	ResetAfter time.Duration
}

// Return the default Backoff Conf.
func Default() Conf {
	return defConf
}

var defConf = Conf{
	Min:        time.Millisecond,
	Max:        time.Second * 30,
	Unit:       time.Second,
	Strategy:   Exponent,
	ResetAfter: time.Second * 30,
}

// Backoff strategy is usually used to determine how long to wait between retries.
type Backoff struct {
	conf Conf

	next     time.Duration
	delta    uint32
	calledAt time.Time
}

type Strategy byte

const (
	Linear Strategy = iota
	Exponent
)

// Create a Backoff.
func New(conf Conf) *Backoff {
	if conf.Strategy < Linear || conf.Strategy > Exponent {
		conf.Strategy = Exponent
	}

	return &Backoff{
		conf:  conf,
		next:  0,
		delta: 1,
	}
}

// Return the next delay.
func (b *Backoff) Next() time.Duration {
	b.resetIf()

	b.calledAt = time.Now().UTC()

	if b.next < b.conf.Min {
		b.next = b.conf.Min
		return b.next
	}

	if b.next == b.conf.Max {
		return b.next
	}

	switch b.conf.Strategy {
	case Linear:
		b.next += b.conf.Unit
		if b.next > b.conf.Max {
			b.next = b.conf.Max
		}

	case Exponent:
		b.next += b.conf.Unit * time.Duration(b.delta)

		if b.next > b.conf.Max {
			b.next = b.conf.Max
		}

		b.delta *= 2
	}

	return b.next
}

func (b *Backoff) resetIf() {
	if time.Since(b.calledAt)-b.next > b.conf.ResetAfter {
		b.next = 0
		b.delta = 1
	}
}
