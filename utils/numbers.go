package utils

import (
	"crypto/rand"
	"math/big"
	t "github.com/adagit94/golib/types"
)

// Clamp clamps any Number value passed as v to range defined by min and max bounds.
func Clamp[T t.Number](v T, min T, max T) T {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

// Simplified wrapper around rand.Int(rand.Reader, big.NewInt(max)).
func RandInt(max int64) (*big.Int, error) {
	return rand.Int(rand.Reader, big.NewInt(max))
}

// Simplified wrapper around rand.Int(rand.Reader, big.NewInt(max)). It panics in case of returned non-nil error.
func RandIntPanic(max int64) *big.Int {
	n, err := RandInt(max)

	if err != nil {
		panic(err)
	}

	return n
}
