/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package iter

// Returns an Iterator consisting of those elements of the Iterator for which fn(e) returns true.
func (it *Iterator[E]) Filter(fn Pred[E]) *Iterator[E] {
	it.src = &filterIter[E]{
		src: it.src,
		fn:  fn,
	}
	return it
}

type filterIter[E any] struct {
	src Iter[E]
	fn  Pred[E]
	end bool
}

func (it *filterIter[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	for op := it.src.Next(); op.Ok; op = it.src.Next() {
		if it.fn(op.Val) {
			return op
		}
	}

	// nothing left
	it.end = true // avoid unnecessary reads thereafter
	it.src = nil  // minimize memory footprint as soon as possible

	return op0
}

// Return an Iterator consisting of the results of applying fn to every element of the Iterator.
func (it *Iterator[E]) Map(fn MapFn[E, E]) *Iterator[E] {
	it.src = &mapIter[E, E]{
		src: it.src,
		fn:  fn,
	}
	return it
}

// Return an Iterator consisting of the results of applying fn to every element of the Iterator.
func Map[E, F any](it *Iterator[E], fn MapFn[E, F]) *Iterator[F] {
	return &Iterator[F]{
		src: &mapIter[E, F]{
			src: it.src,
			fn:  fn,
		},
	}
}

type mapIter[E, F any] struct {
	src Iter[E]
	fn  MapFn[E, F]
	end bool
}

func (it *mapIter[E, F]) Next() Option[F] {
	var op0 Option[F]

	if it.end {
		return op0
	}

	op := it.src.Next()
	if !op.Ok {
		it.end = true
		it.src = nil
		return op0
	}

	return Option[F]{
		Val: it.fn(op.Val),
		Ok:  true,
	}
}

// Return an Iterator consisting of the first n elements of the Iterator, or all elements if there are fewer than n.
func (it *Iterator[E]) Take(n int) *Iterator[E] {
	it.src = &take[E]{
		src: it.src,
		n:   n,
	}
	return it
}

type take[E any] struct {
	src Iter[E]
	n   int
	end bool
}

func (it *take[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	if it.n < 1 {
		it.end = true
		it.src = nil
		return op0
	}

	op := it.src.Next()
	if !op.Ok {
		it.end = true
		it.src = nil
		return op0
	}

	it.n--
	return op
}

// Returns an Iterator consisting of all but the first n elements of the Iterator.
func (it *Iterator[E]) Drop(n int) *Iterator[E] {
	it.src = &drop[E]{
		src: it.src,
		n:   n,
	}
	return it
}

type drop[E any] struct {
	src Iter[E]
	n   int
	end bool
}

func (it *drop[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	for ; it.n > 0; it.n-- {
		op := it.src.Next()
		if !op.Ok {
			it.end = true
			it.src = nil
			return op0
		}
	}

	op := it.src.Next()
	if !op.Ok {
		it.end = true
		it.src = nil
	}

	return op
}

// Return an Iterator consisting of those elements of the Iterator as long as fn(e) returns true. Once fn(e) returns false, the rest of the elements are ignored.
func (it *Iterator[E]) TakeWhile(fn Pred[E]) *Iterator[E] {
	it.src = &takeWhile[E]{
		src: it.src,
		fn:  fn,
	}
	return it
}

type takeWhile[E any] struct {
	src     Iter[E]
	fn      Pred[E]
	end     bool
	endPred bool
}

func (it *takeWhile[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	op := it.src.Next()
	if !op.Ok || !it.fn(op.Val) {
		it.end = true
		it.src = nil
		return op0
	}

	return op
}

// Return an Iterator consisting of those elements of the Iterator starting from the first element for which fn(e) returns false.
func (it *Iterator[E]) DropWhile(fn Pred[E]) *Iterator[E] {
	it.src = &dropWhile[E]{
		src: it.src,
		fn:  fn,
	}
	return it
}

type dropWhile[E any] struct {
	src     Iter[E]
	fn      Pred[E]
	end     bool
	endDrop bool
}

func (it *dropWhile[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	op := it.src.Next()

	if !it.endDrop {
		for ; op.Ok; op = it.src.Next() {
			if !it.fn(op.Val) {
				it.endDrop = true
				break
			}
		}
	}

	if !op.Ok {
		it.end = true
		it.src = nil
		return op0
	}

	return op
}

// Return an Iterator consisting of all the elements of the first Iterator, followed by all the elements of the second Iterator, and so on.
func Chain[E any](its ...*Iterator[E]) *Iterator[E] {
	lsrc := make([]Iter[E], len(its))
	for i, it := range its {
		lsrc[i] = it.src
	}

	return &Iterator[E]{
		src: &chainIter[E]{
			lsrc: lsrc,
		},
	}
}

type chainIter[E any] struct {
	lsrc []Iter[E]
	end  bool
}

func (it *chainIter[E]) Next() Option[E] {
	var op0 Option[E]

	if it.end {
		return op0
	}

	for len(it.lsrc) > 0 {
		op := it.lsrc[0].Next()
		if !op.Ok {
			it.lsrc = it.lsrc[1:]
			continue
		}

		return op
	}

	it.end = true
	it.lsrc = nil

	return op0
}
