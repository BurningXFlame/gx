/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

/*
A Timeout Decorator abandons the call to the wrapped function if the call does not finish in a specified duration.
There're 4 variants, each for a type of functions. All functions can be classfied into these four.
*/
package timeout

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("timeout")

// Timeout Decorator for functions with input and output parameters, i.e func(I) (O, error)
func WithTimeoutIO[I, O any](timeout time.Duration, fn func(I) (O, error)) func(I) (O, error) {
	if timeout <= 0 {
		return fn
	}

	return func(in I) (O, error) {
		ch := make(chan result[O], 1)
		go func() {
			val, err := fn(in)
			ch <- result[O]{val, err}
		}()

		select {
		case <-time.After(timeout):
			var out O
			return out, ErrTimeout
		case r := <-ch:
			return r.val, r.err
		}
	}
}

type result[T any] struct {
	val T
	err error
}

// Timeout Decorator for functions with input parameters only, i.e func(I) error
func WithTimeoutI[I any](timeout time.Duration, fn func(I) error) func(I) error {
	fn1 := func(in I) (none, error) {
		err := fn(in)
		return none{}, err
	}

	fn2 := WithTimeoutIO(timeout, fn1)

	return func(in I) error {
		_, err := fn2(in)
		return err
	}
}

type none = struct{}

// Timeout Decorator for functions with output parameters only, i.e func() (O, error)
func WithTimeoutO[O any](timeout time.Duration, fn func() (O, error)) func() (O, error) {
	fn1 := func(none) (O, error) {
		return fn()
	}

	fn2 := WithTimeoutIO(timeout, fn1)

	return func() (O, error) {
		return fn2(none{})
	}
}

// Timeout Decorator for functions with neither input nor output parameters, i.e func() error
func WithTimeout(timeout time.Duration, fn func() error) func() error {
	fn1 := func(none) (none, error) {
		return none{}, fn()
	}

	fn2 := WithTimeoutIO(timeout, fn1)

	return func() error {
		_, err := fn2(none{})
		return err
	}
}
