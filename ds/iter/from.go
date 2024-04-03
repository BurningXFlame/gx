/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

// Create an Iterator consisting of all the elements of the first Iter, followed by all the elements of the second Iter, and so on.
func From[E any](its ...Iter[E]) *Iterator[E] {
	return &Iterator[E]{
		src: &chainIter[E]{
			lsrc: its,
		},
	}
}

// Create an Iterator from a slice
func FromSlice[S ~[]E, E any](s S) *Iterator[E] {
	return &Iterator[E]{
		src: &sliceIter[E]{
			src: s,
		},
	}
}

type sliceIter[E any] struct {
	src []E
	end bool
}

func (it *sliceIter[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	if len(it.src) == 0 {
		it.end = true
		it.src = nil
		return op0
	}

	val := it.src[0]
	it.src = it.src[1:]

	return Some(val)
}

type KV[K comparable, V any] struct {
	Key K
	Val V
}

// Create an Iterator from a map
func FromMap[M ~map[K]V, K comparable, V any](m M) *Iterator[KV[K, V]] {
	l := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		l = append(l, KV[K, V]{Key: k, Val: v})
	}

	return FromSlice(l)
}

// Create an Iterator from a channel
func FromChan[C ~<-chan E, E any](ch C) *Iterator[E] {
	return &Iterator[E]{
		src: &chanIter[E]{
			src: ch,
		},
	}
}

type chanIter[E any] struct {
	src <-chan E
	end bool
}

func (it *chanIter[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	val, ok := <-it.src
	if !ok {
		it.end = true
		it.src = nil
		return op0
	}

	return Some(val)
}

// Create an Iterator from a Ranger
func FromRanger[E any](r Ranger[E]) *Iterator[E] {
	var l []E

	r.ForEach(func(e E) {
		l = append(l, e)
	})

	return FromSlice(l)
}
