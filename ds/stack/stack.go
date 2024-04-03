/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package stack

import "github.com/burningxflame/gx/ds/ringbuf"

type Stack[E any] struct {
	r *ringbuf.RingBuf[E]
}

// Create a new Stack. An optional parameter may be provided to specify initial capacity.
func New[E any](initCap ...int) *Stack[E] {
	return &Stack[E]{
		r: ringbuf.New[E](initCap...),
	}
}

func (s *Stack[E]) Push(e E) {
	s.r.PushFront(e)
}

func (s *Stack[E]) Peek() (E, bool) {
	return s.r.PeekFront()
}

func (s *Stack[E]) Pop() (E, bool) {
	return s.r.PopFront()
}

func (s *Stack[E]) Len() int {
	return s.r.Len()
}

// Call fn sequentially for each element in the Stack. If fn returns false, stop iteration.
// The iteration order is LIFO.
func (s *Stack[E]) Range(fn func(E) bool) {
	s.r.Range(fn)
}

// Call fn sequentially for each element in the Stack.
// The iteration order is LIFO.
func (s *Stack[E]) ForEach(fn func(E)) {
	s.Range(func(e E) bool {
		fn(e)
		return true
	})
}

func (s *Stack[E]) Add(e E) {
	s.Push(e)
}
