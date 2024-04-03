/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package sem

import (
	"math"
	"sync"
)

// Keyed-Semaphores. Commonly used for limiting max concurrency per key, e.g. max concurrent connections per client.
type KSem[K comparable] struct {
	ca       int
	sizeHint int
	sems     sync.Map
	mu       sync.Mutex
	size     int
}

// Create Keyed-Semaphores.
// The ca specifies the capacity of every semaphore.
// If the number of semaphores exceeds sizeHint, will try to shrink, i.e. remove semaphores who have no permits taken.
func NewKSem[K comparable](ca int, sizeHint int) *KSem[K] {
	if ca < 1 {
		ca = defCap
	}

	if sizeHint < 1 {
		sizeHint = math.MaxInt
	}

	return &KSem[K]{
		ca:       ca,
		sizeHint: sizeHint,
	}
}

// Get the semaphore of the key, create if not exist.
func (ks *KSem[K]) Get(key K) *Sem {
	s, ok := ks.sems.Load(key)
	if ok {
		return s.(*Sem)
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()

	// may have been created
	s, ok = ks.sems.Load(key)
	if ok {
		return s.(*Sem)
	}

	if ks.size >= ks.sizeHint {
		ks.shrink()
	}

	s2 := New(ks.ca)

	ks.sems.Store(key, s2)
	ks.size++

	return s2
}

func (ks *KSem[K]) shrink() {
	ks.sems.Range(func(key, value any) bool {
		s := value.(*Sem)
		if s.Available() == ks.ca {
			ks.sems.Delete(key)
			ks.size--
		}
		return true
	})
}
