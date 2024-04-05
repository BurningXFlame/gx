/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

import "sort"

// Call fn for each element of the Iterator.
// This is a draining operation.
func (it *Iterator[E]) ForEach(fn Unary[E]) {
	for op := it.src.Next(); op.Ok; op = it.src.Next() {
		fn(op.Val)
	}
}

// Return an Option wrapping the min element of the Iterator if exist.
// This is a draining operation.
func (it *Iterator[E]) Min(less Less[E]) Option[E] {
	var min Option[E]

	it.ForEach(func(e E) {
		if !min.Ok {
			min.Ok = true
			min.Val = e
		}

		if less(e, min.Val) {
			min.Val = e
		}
	})

	return min
}

// Return an Option wrapping the max element of the Iterator if exist.
// This is a draining operation.
func (it *Iterator[E]) Max(less Less[E]) Option[E] {
	return it.Min(func(a, b E) bool {
		return !less(a, b)
	})
}

// Return a sorted slice of all element of the Iterator.
// This is a draining operation.
func (it *Iterator[E]) Sort(less Less[E]) []E {
	l := it.ToSlice()

	sort.Slice(l, func(i, j int) bool {
		return less(l[i], l[j])
	})

	return l
}

// Return the result of applying fn to ini and the first element of the Iterator,
// then applying fn to that result and the second element, and so on.
// If the Iterator is empty, return ini and fn is not called.
// This is a draining operation.
func (it *Iterator[E]) Reduce(ini E, fn func(acc E, e E) E) E {
	return Reduce(it, ini, fn)
}

// Return the result of applying fn to ini and the first element of the Iterator,
// then applying fn to that result and the second element, and so on.
// If the Iterator is empty, return ini and fn is not called.
// This is a draining operation.
func Reduce[A any, E any](it *Iterator[E], ini A, fn func(acc A, e E) A) A {
	acc := ini

	it.ForEach(func(e E) {
		acc = fn(acc, e)
	})

	return acc
}

// Call fn sequentially for each element of the Iterator. If fn returns false, stop iteration.
// This is a draining operation.
func (it *Iterator[E]) Range(fn Pred[E]) {
	for op := it.src.Next(); op.Ok; op = it.src.Next() {
		if !fn(op.Val) {
			return
		}
	}
}

// Return true if fn(e) is true for any element of the Iterator.
// If the Iterator is empty, return false.
// This is a draining operation.
func (it *Iterator[E]) Any(fn Pred[E]) bool {
	ok := false

	it.Range(func(e E) bool {
		if fn(e) {
			ok = true
			return false
		}

		return true
	})

	return ok
}

// Return true if fn(e) is true for every element of the Iterator.
// If the Iterator is empty, return true.
// This is a draining operation.
func (it *Iterator[E]) Every(fn Pred[E]) bool {
	return !it.Any(func(e E) bool {
		return !fn(e)
	})
}
