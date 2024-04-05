/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package d2

import "sync"

func c5() int32 {
	var total int32

	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i + 1
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i, e := range l {
		go _c3(&total, int32(e), &wg)
		_ = i
	}
	wg.Wait()

	return total
}
