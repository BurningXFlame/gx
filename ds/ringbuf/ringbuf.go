/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// Auto-Scalable ringbuf (Ring Buffer) and deque (Double Ended Queue).
package ringbuf

type RingBuf[E any] struct {
	buf  []E
	max  int
	len  int
	head int
	tail int
}

// Create a new RingBuf. An optional parameter may be provided to specify initial capacity.
func New[E any](initCap ...int) *RingBuf[E] {
	var ca int
	if len(initCap) > 0 {
		ca = initCap[0]
	}
	if ca < 1 {
		ca = 1
	}

	return &RingBuf[E]{
		buf: make([]E, ca),
		max: ca - 1,
	}
}

func (r *RingBuf[E]) PushBack(e E) {
	r.scaleIf()

	r.buf[r.tail] = e
	r.tail = r.next(r.tail)
	r.len++
}

func (r *RingBuf[E]) PushFront(e E) {
	r.scaleIf()

	r.head = r.prev(r.head)
	r.buf[r.head] = e
	r.len++
}

func (r *RingBuf[E]) next(i int) int {
	if i == r.max {
		return 0
	}

	return i + 1
}

func (r *RingBuf[E]) prev(i int) int {
	if i == 0 {
		return r.max
	}

	return i - 1
}

func (r *RingBuf[E]) PeekBack() (E, bool) {
	if r.len == 0 {
		var z E
		return z, false
	}

	i := r.prev(r.tail)
	return r.buf[i], true
}

func (r *RingBuf[E]) PeekFront() (E, bool) {
	if r.len == 0 {
		var z E
		return z, false
	}

	return r.buf[r.head], true
}

func (r *RingBuf[E]) PopBack() (E, bool) {
	var z E

	if r.len == 0 {
		return z, false
	}

	i := r.prev(r.tail)
	e := r.buf[i]
	r.buf[i] = z
	r.len--
	r.tail = i

	return e, true
}

func (r *RingBuf[E]) PopFront() (E, bool) {
	var z E

	if r.len == 0 {
		return z, false
	}

	e := r.buf[r.head]
	r.buf[r.head] = z
	r.len--
	r.head = r.next(r.head)

	return e, true
}

func (r *RingBuf[E]) Len() int {
	return r.len
}

// Call fn sequentially for each element in the RingBuf. If fn returns false, stop iteration.
// The iteration order is head to tail.
func (r *RingBuf[E]) Range(fn func(E) bool) {
	if r.len < 1 {
		return
	}

	if r.tail > r.head {
		for _, e := range r.buf[r.head:r.tail] {
			if !fn(e) {
				return
			}
		}
		return
	}

	for _, e := range r.buf[r.head:] {
		if !fn(e) {
			return
		}
	}

	for _, e := range r.buf[:r.tail] {
		if !fn(e) {
			return
		}
	}
}

// Call fn sequentially for each element in the RingBuf.
// The iteration order is head to tail.
func (r *RingBuf[E]) ForEach(fn func(E)) {
	r.Range(func(e E) bool {
		fn(e)
		return true
	})
}

// Scale up if necessary
func (r *RingBuf[E]) scaleIf() {
	if r.len < cap(r.buf) {
		return
	}

	// similar to how slice grows
	const threshold = 256
	oldCap := cap(r.buf)
	var newCap int
	if oldCap < threshold {
		newCap = oldCap * 2
	} else {
		newCap = oldCap + (oldCap+3*threshold)/4
	}

	newBuf := make([]E, newCap)
	if r.tail > r.head {
		copy(newBuf, r.buf[r.head:r.tail])
	} else {
		n := copy(newBuf, r.buf[r.head:])
		copy(newBuf[n:], r.buf[:r.tail])
	}

	r.buf = newBuf
	r.max = newCap - 1
	r.head = 0
	r.tail = r.len
}

func (r *RingBuf[E]) Add(e E) {
	r.PushBack(e)
}
