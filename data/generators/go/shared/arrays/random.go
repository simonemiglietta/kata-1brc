package arrays

import (
	"fmt"
	"math/rand"
)

func ExtractRandElements[T any](source []T, elements int) ([]T, error) {
	if elements < 0 {
		return nil, fmt.Errorf("elements cannot be negative")
	}

	sourceLen := len(source)
	if elements > sourceLen {
		return nil, fmt.Errorf("not enough elements in source")
	}

	sourceLocal := make([]T, sourceLen)
	copy(sourceLocal, source)
	result := make([]T, elements)

	for i := 0; i < elements; i++ {
		// fetches one item
		extractedIndex := rand.Intn(sourceLen)
		result[i] = sourceLocal[extractedIndex]

		// shortens the source moving the last element in place of fetched one
		sourceLen--
		sourceLocal[extractedIndex] = sourceLocal[sourceLen]
	}

	return result, nil
}
