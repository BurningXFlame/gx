/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package d1

import (
	"sync"
	"sync/atomic"
)

func c2() int32 {
	var total int32

	var wg sync.WaitGroup
	wg.Add(n)
	for e := 1; e <= n; e++ {
		go func(e int) {
			defer wg.Done()
			atomic.AddInt32(&total, int32(e))
		}(e)
	}
	wg.Wait()

	return total
}
