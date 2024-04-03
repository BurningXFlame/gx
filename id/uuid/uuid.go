/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync"
	"time"

	xrand "golang.org/x/exp/rand"
)

// Generate a UUID v4
func New() string {
	id := make([]byte, size)
	cr.Read(id) // It always returns len(p) and a nil error.
	v4(id)
	return hex.EncodeToString(id)
}

const size = 16

var cr = newConRand()

// Return a concurrency-safe rand generator
func newConRand() *xrand.Rand {
	src := &xrand.LockedSource{}

	var seed uint64
	err := binary.Read(rand.Reader, binary.BigEndian, &seed)
	if err != nil {
		seed = uint64(time.Now().UnixNano())
	}

	src.Seed(seed)
	return xrand.New(src)
}

func v4(id []byte) {
	if len(id) < size {
		return
	}

	id[6] = (id[6] & 0x0f) | 0x40 // ver: 0x4
	id[8] = (id[8] & 0x3f) | 0x80 // variant: 0b10
}

// Generate a secure UUID v4
func NewSec() (string, error) {
	id := make([]byte, size)
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}

	v4(id)
	return hex.EncodeToString(id), nil
}

// Create a generator which generates secure UUID v4 in batches, for better performance.
// The batchSize specifies the number of UUIDs to generate in each batch.
func SecGenerator(batchSize int) func() (string, error) {
	byteSize := batchSize * size
	bs := make([]byte, byteSize)
	offset := byteSize

	return func() (string, error) {
		if offset == byteSize {
			_, err := rand.Read(bs)
			if err != nil {
				return "", err
			}

			offset = 0
		}

		id := bs[offset : offset+size]
		v4(id)
		sid := hex.EncodeToString(id)
		offset += size
		return sid, nil
	}
}

// Concurrency-safe version of SecGenerator
func SecConGenerator(batchSize int) func() (string, error) {
	byteSize := batchSize * size
	bs := make([]byte, byteSize)
	offset := byteSize
	var mu sync.Mutex

	return func() (string, error) {
		mu.Lock()
		defer mu.Unlock()

		if offset == byteSize {
			_, err := rand.Read(bs)
			if err != nil {
				return "", err
			}

			offset = 0
		}

		id := bs[offset : offset+size]
		v4(id)
		sid := hex.EncodeToString(id)
		offset += size
		return sid, nil
	}
}
