/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package d2

import "sync"

func c4() int32 {
	var total int32

	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i + 1
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for _, e := range l {
		go _c3(&total, int32(e), &wg)
	}
	wg.Wait()

	return total
}
