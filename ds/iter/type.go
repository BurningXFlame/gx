/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

// Represent a lazy iteration
type Iter[E any] interface {
	// Advance the iteration and return an Option wrapping the next value if exist.
	Next() Option[E]
}

// Represent an optional value
type Option[E any] struct {
	Val E    // the value if exist
	Ok  bool // true if exist, false otherwise
}

// Some value of type E
func Some[E any](e E) Option[E] {
	return Option[E]{
		Val: e,
		Ok:  true,
	}
}

type Iterator[E any] struct {
	src Iter[E]
}

type Unary[E any] func(E)
type Pred[E any] func(E) bool
type MapFn[E, F any] func(E) F

// return true if a < b
type Less[E any] func(a, b E) bool

// Represent a collector which collects all elements of an Iterator.
type Collector[E any] interface {
	// Add the element into the Collector
	Add(E)
}

// Represent a normal iteration (not lazy)
type Ranger[E any] interface {
	// Iterate the collection and call fn for each element.
	ForEach(fn func(E))
}
