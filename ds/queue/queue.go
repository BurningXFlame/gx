/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package queue

import "github.com/burningxflame/gx/ds/ringbuf"

type Queue[E any] struct {
	r *ringbuf.RingBuf[E]
}

// Create a new Queue. An optional parameter may be provided to specify initial capacity.
func New[E any](initCap ...int) *Queue[E] {
	return &Queue[E]{
		r: ringbuf.New[E](initCap...),
	}
}

func (q *Queue[E]) Enq(e E) {
	q.r.PushBack(e)
}

func (q *Queue[E]) Peek() (E, bool) {
	return q.r.PeekFront()
}

func (q *Queue[E]) Deq() (E, bool) {
	return q.r.PopFront()
}

func (q *Queue[E]) Len() int {
	return q.r.Len()
}

// Call fn sequentially for each element in the Queue. If fn returns false, stop iteration.
// The iteration order is FIFO.
func (q *Queue[E]) Range(fn func(E) bool) {
	q.r.Range(fn)
}

// Call fn sequentially for each element in the Queue.
// The iteration order is FIFO.
func (q *Queue[E]) ForEach(fn func(E)) {
	q.Range(func(e E) bool {
		fn(e)
		return true
	})
}

func (q *Queue[E]) Add(e E) {
	q.Enq(e)
}
