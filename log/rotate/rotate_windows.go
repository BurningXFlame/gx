/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package rotate

// Windows prevents files from being deleted if in use. So, nothing to do.
func (r *rotate) reCreateIf() error {
	return nil
}
