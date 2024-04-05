/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package set

type Set[E comparable] struct {
	m map[E]struct{}
}

// Create a new Set. An optional parameter may be provided to specify initial capacity.
func New[E comparable](initCap ...int) *Set[E] {
	var ca int
	if len(initCap) > 0 {
		ca = initCap[0]
	}
	if ca < 1 {
		return &Set[E]{
			m: make(map[E]struct{}),
		}
	}

	return &Set[E]{
		m: make(map[E]struct{}, ca),
	}
}

func (s *Set[E]) Add(e E) {
	s.m[e] = struct{}{}
}

func (s *Set[E]) Delete(e E) {
	delete(s.m, e)
}

func (s *Set[E]) Len() int {
	return len(s.m)
}

func (s *Set[E]) Contain(x E) bool {
	_, ok := s.m[x]
	return ok
}

func (s *Set[E]) Equal(x *Set[E]) bool {
	if s.Len() != x.Len() {
		return false
	}

	for e := range s.m {
		if !x.Contain(e) {
			return false
		}
	}

	return true
}

// Call fn sequentially for each element in the Set. If fn returns false, stop iteration.
// The iteration order is unspecified.
func (s *Set[E]) Range(fn func(E) bool) {
	for e := range s.m {
		if !fn(e) {
			return
		}
	}
}

// Call fn sequentially for each element in the Set.
// The iteration order is unspecified.
func (s *Set[E]) ForEach(fn func(E)) {
	s.Range(func(e E) bool {
		fn(e)
		return true
	})
}

// s ∪ x
func (s *Set[E]) Union(x *Set[E]) *Set[E] {
	ca := max(s.Len(), x.Len())
	rs := New[E](ca)

	for e := range s.m {
		rs.Add(e)
	}

	for e := range x.m {
		rs.Add(e)
	}

	return rs
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// s ∩ x
func (s *Set[E]) Intersect(x *Set[E]) *Set[E] {
	rs := New[E]()

	if s.Len() > x.Len() {
		s, x = x, s
	}

	for e := range s.m {
		if x.Contain(e) {
			rs.Add(e)
		}
	}

	return rs
}

// s - x
func (s *Set[E]) Diff(x *Set[E]) *Set[E] {
	rs := New[E]()

	for e := range s.m {
		if !x.Contain(e) {
			rs.Add(e)
		}
	}

	return rs
}
