/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package d2

import (
	"sync"
	"sync/atomic"
)

const n = 100

func c3() int32 {
	var total int32

	var wg sync.WaitGroup
	wg.Add(n)
	for e := 1; e <= n; e++ {
		go _c3(&total, int32(e), &wg)
	}
	wg.Wait()

	return total
}

func _c3(total *int32, e int32, wg *sync.WaitGroup) {
	defer wg.Done()
	atomic.AddInt32(total, e)
}
