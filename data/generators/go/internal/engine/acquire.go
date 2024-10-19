package engine

import (
	"bufio"
	"fmt"
	models "lvciot/generate/internal/entites"
	"lvciot/generate/shared/arrays"
	"maps"
	"os"
	"slices"
)

func RetrieveStations(f *os.File, total int) ([]models.Station, error) {
	rows := make(map[string]bool, total)

	s := bufio.NewScanner(f)

	// first two lines are comments and not actual data
	s.Scan()
	s.Scan()

	for s.Scan() {
		rows[s.Text()] = true
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	selected, err := arrays.ExtractRandElements(slices.Collect(maps.Keys(rows)), total)
	if err != nil {
		return nil, fmt.Errorf("cannot extract stations: %w", err)
	}

	result := make([]models.Station, len(selected))
	for i, v := range selected {
		result[i], err = models.NewStationFromString(v)
		if err != nil {
			return nil, fmt.Errorf("cannot parse station: %w", err)
		}
	}

	return result, nil
}
