package utils

import (
	"slices"
)

// MapSlice maps S passed as argument with elements of type E to different type EE returned by the passed function. len(E) == len(EE) and new slice with underlaying array is allocated without mutation of original one.
func MapSlice[S ~[]E, E any, EE any](s S, f func(e E, i int) EE) []EE {
	ss := make([]EE, len(s))

	for i, e := range s {
		ss[i] = f(e, i)
	}

	return ss
}

// DifferenceSlice returns new slice with elements from s1 that aren't found in s2. It uses slices.Contains behind the scenes to compare element values.
func DifferenceSlice[S ~[]E, E comparable](s1 S, s2 S) S {
	s := make(S, 0)

	for _, v := range s1 {
		if !slices.Contains(s2, v) {
			s = append(s, v)
		}
	}

	return s
}

// FlatSlice flattens elements of nested slices of 2-dim slice into new single dimensional one (1 level deep).
func FlatSlice[S ~[][]E, E any](s S) []E {
	ss := make([]E, 0)

	for _, e := range s {
		for _, ee := range e {
			ss = append(ss, ee)
		}
	}

	return ss
}

// ExpandRange expands passed range (from, to) and returns slice with intermediary values included. It's inclusive of destination (to) value, works for both positive and negative ranges and in both ascending and descending order. E.g. [1, 3] returns [1, 2, 3]; [3, 1] returns [3, 2, 1]; [-1, -3] returns [-1, -2, -3]; [-3, -1] returns [-3, -2, -1]; [-1, 1] returns [-1, 0, 1]; [1, -1] returns [1, 0, -1].
func ExpandRange(from int, to int) []int {
	interval := to - from

	if interval < 0 {
		interval *= -1
	}

	// +1 to include final value also.
	interval += 1

	r := make([]int, 0, interval)

	if from < to {
		for n := from; n <= to; n++ {
			r = append(r, n)
		}
	} else {
		for n := from; n >= to; n-- {
			r = append(r, n)
		}
	}

	return r
}

// Replace elements from a to b (exclusive) with passed compatible value e and return modified slice.
func ReplaceRepeatedSlice[S ~[]E, E any](s S, e E, a, b int) S {
	for i := a; i < b; i++ {
		s[i] = e
	}

	return s
}

// Replace every element in s with passed compatible value e and return modified slice.
func ReplaceRepeatedEverySlice[S ~[]E, E any](s S, e E) S {
	for i := range len(s) {
		s[i] = e
	}

	return s
}

// Creates and fills new slice based on passed capacity with element value.
func NewFilledSlice[E any](element E, capacity int) []E {
	s := make([]E, 0, capacity)

	for range cap(s) {
		s = append(s, element)
	}

	return s
}

// Reduce a slice to the single value of same type as it's elements.
func ReduceSlice[S ~[]E, E any](s S, f func(acc, el E, i int) E) E {
	v := s[0]

	for i := 1; i < len(s); i++ {
		v = f(v, s[i], i)
	}

	return v
}

// Reduce a slice to the single value of potentially different type than it's elements. With initial value provided.
func ReduceSliceToValue[S ~[]E, E any, V any](s S, v V, f func(acc V, el E, i int) V) V {
	for i, e := range s {
		v = f(v, e, i)
	}

	return v
}
