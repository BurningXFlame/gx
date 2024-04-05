/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package heap

import "container/heap"

type Heap[E any] struct {
	h *hp[E]
}

// return true if a < b
type Less[E any] func(a, b E) bool

// Create a new Heap. An optional parameter may be provided to specify initial capacity.
func New[E any](less Less[E], initCap ...int) *Heap[E] {
	var ca int
	if len(initCap) > 0 {
		ca = initCap[0]
	}
	if ca < 1 {
		return &Heap[E]{
			&hp[E]{
				less: less,
			},
		}
	}

	return &Heap[E]{
		&hp[E]{
			xs:   make([]E, 0, ca),
			less: less,
		},
	}
}

func (h *Heap[E]) Push(e E) {
	heap.Push(h.h, e)
}

func (h *Heap[E]) Pop() (E, bool) {
	if h.h.Len() == 0 {
		var z E
		return z, false
	}

	return heap.Pop(h.h).(E), true
}

func (h *Heap[E]) Len() int {
	return h.h.Len()
}

// Call fn sequentially for each element in the Heap. If fn returns false, stop iteration.
// The iteration order is not specified.
func (h *Heap[E]) Range(fn func(E) bool) {
	for _, e := range h.h.xs {
		if !fn(e) {
			return
		}
	}
}

// Call fn sequentially for each element in the Heap.
// The iteration order is not specified.
func (h *Heap[E]) ForEach(fn func(E)) {
	h.Range(func(e E) bool {
		fn(e)
		return true
	})
}

func (h *Heap[E]) Add(e E) {
	h.Push(e)
}

type hp[E any] struct {
	xs   []E
	less Less[E]
}

func (h *hp[E]) Len() int {
	return len(h.xs)
}

func (h *hp[E]) Less(i, j int) bool {
	return h.less(h.xs[i], h.xs[j])
}

func (h *hp[E]) Swap(i, j int) {
	h.xs[i], h.xs[j] = h.xs[j], h.xs[i]
}

func (h *hp[E]) Push(x any) {
	h.xs = append(h.xs, x.(E))
}

func (h *hp[E]) Pop() any {
	old := h.xs
	i := len(old) - 1
	x := old[i]
	h.xs = old[:i]
	return x
}
