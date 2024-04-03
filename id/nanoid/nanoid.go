/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

// NanoID
package nanoid

import (
	"crypto/rand"
	"sync"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
	maxIndex = 63
	size     = 21
)

// Generate a NanoID
func New() (string, error) {
	bs := make([]byte, size)
	_, err := rand.Read(bs)
	if err != nil {
		return "", err
	}

	for i := 0; i < size; i++ {
		index := bs[i] & maxIndex
		bs[i] = alphabet[index]
	}

	return string(bs), nil
}

// Create a generator which generates NanoID in batches, for better performance.
// The batchSize specifies the number of NanoIDs to generate in each batch.
func Generator(batchSize int) func() (string, error) {
	byteSize := batchSize * size
	bs := make([]byte, byteSize)
	offset := byteSize
	id := make([]byte, size)

	return func() (string, error) {
		if offset == byteSize {
			_, err := rand.Read(bs)
			if err != nil {
				return "", err
			}

			offset = 0
		}

		for i := 0; i < size; i++ {
			index := bs[offset+i] & maxIndex
			id[i] = alphabet[index]
		}

		offset += size

		return string(id), nil
	}
}

// Concurrency-safe version of Generator
func ConGenerator(batchSize int) func() (string, error) {
	byteSize := batchSize * size
	bs := make([]byte, byteSize)
	offset := byteSize
	id := make([]byte, size)
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

		for i := 0; i < size; i++ {
			index := bs[offset+i] & maxIndex
			id[i] = alphabet[index]
		}

		offset += size

		return string(id), nil
	}
}
