/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

// Feed the Collector c with all elements of the Iterator.
func (it *Iterator[E]) To(c Collector[E]) {
	it.ForEach(func(e E) {
		c.Add(e)
	})
}

// Return a slice of all elements of the Iterator.
// This is a draining operation.
func (it *Iterator[E]) ToSlice() []E {
	var l []E

	it.ForEach(func(e E) {
		l = append(l, e)
	})

	return l
}

// Return a map of all elements of the Iterator.
// This is a draining operation.
func ToMap[K comparable, V any](it *Iterator[KV[K, V]]) map[K]V {
	m := make(map[K]V)

	it.ForEach(func(kv KV[K, V]) {
		m[kv.Key] = kv.Val
	})

	return m
}

// Return a channel of all elements of the Iterator.
// This is a draining operation.
func (it *Iterator[E]) ToChan() <-chan E {
	ch := make(chan E, 1)

	go func() {
		it.ForEach(func(e E) {
			ch <- e
		})
		close(ch)
	}()

	return ch
}
