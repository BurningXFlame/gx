/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package gid

import (
	"unsafe"
)

// Return the ID of the current goroutine.
// More precisely speaking, it's the address of the current g (i.e. goroutine). Therefore, the ID makes sense only during the life cycle of the current goroutine.
func Gid() uint {
	return uint(uintptr(g()))
}

// Pointer to the current g (i.e. goroutine)
func g() unsafe.Pointer
