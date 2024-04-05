/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package conbuf

import "io"

type WriteFlusher interface {
	io.Writer

	Flush() error
}
