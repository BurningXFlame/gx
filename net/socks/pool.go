/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import "sync"

const maxBufSize = 263

var bufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, maxBufSize)
		return &buf
	},
}

func getBuf() *[]byte {
	return bufPool.Get().(*[]byte)
}

func putBuf(buf *[]byte) {
	bufPool.Put(buf)
}
