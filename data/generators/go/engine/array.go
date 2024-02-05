package engine

import "math/rand"

func ExtractRandom[T any](items []T, limit int) []T {
	r := make([]T, limit)

	end := len(items)
	if limit <= end {
		return items
	}

	for i := 0; i < limit; i++ {
		n := rand.Intn(end)
		r[i] = items[n]
		items[n] = items[end]
	}

	return r
}
