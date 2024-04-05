/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package rotate

import (
	"errors"
	"os"
)

// re-create log file if deleted
func (r *rotate) reCreateIf() error {
	_, err := os.Stat(r.conf.FilePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return r.openLogFile()
	}

	return nil
}
