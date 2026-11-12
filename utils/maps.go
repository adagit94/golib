package utils

import (
	"maps"
	"slices"
)

// SwapMap returns a new map with keys being swapped for respective values and vice versa.
func SwapMap[K comparable, V comparable](m map[K]V) map[V]K {
	keys := slices.Collect(maps.Keys(m))
	m2 := make(map[V]K, len(keys))

	for _, k := range keys {
		m2[m[k]] = k
	}

	return m2
}
