package arrays

import (
	"maps"
	"slices"
)

func ToUnique[T comparable](a []T) []T {
	m := make(map[T]bool)

	for _, v := range a {
		m[v] = true
	}

	return slices.Collect(maps.Keys(m))
}

func HasOnlyUnique[T comparable](a []T) bool {
	m := make(map[T]bool)

	for _, v := range a {
		_, ok := m[v]

		if ok {
			return false
		}
	}

	return true
}
