package utils

import (
	"hash/maphash"
)

// Thread-safe equivalent of rand.Intn that generates values in range: [0, n).
// See rand.Intn specification for more details. This implementation does not support seeds.
func CIntn(n int) int {
	if n <= 0 {
		panic("random-utils: the random int upper limit must be greater than zero")
	}

	i64 := int64(new(maphash.Hash).Sum64())
	if i64 < 0 {
		i64 = -i64
	}

	return int(i64) % n
}
