/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package gls

import (
	"context"
	"sync"

	"github.com/burningxflame/gx/runtime/gid"
)

// Use function `Go` instead of the keyword `go` to spawn a goroutine. The spawned goroutine inherits local storage from its parent goroutine.
func Go(fn func()) {
	ls := getStore()

	go func() {
		setStore(ls)
		defer dropStore()

		fn()
	}()
}

// Associate the value with the key in the current goroutine local storage.
func Put(key any, value any) {
	ls := getStore()
	ls = context.WithValue(ls, key, value)
	setStore(ls)
}

// Return the value corresponding to the key in the current goroutine local storage.
func Get(key any) any {
	ls := getStore()
	return ls.Value(key)
}

// the map which holds all goroutine local storages
var mStore sync.Map

// local storage
type store = context.Context

// zero value of local storage
var s0 = context.Background()

// Return the current goroutine local storage.
func getStore() store {
	id := gid.Gid()
	ls, _ := mStore.LoadOrStore(id, s0)
	return ls.(store)
}

// Set the current goroutine local storage.
func setStore(ls store) {
	id := gid.Gid()
	mStore.Store(id, ls)
}

// Drop the current goroutine local storage.
func dropStore() {
	id := gid.Gid()
	mStore.Delete(id)
}
