package engine

import (
	"bufio"
	models "lvciot/generate/internal/entites"
	"math/rand"
	"os"
)

func GenerateMeasurements(f *os.File, stations []models.Station, total int, c *int) error {
	w := bufio.NewWriter(f)

	for i := 0; i < total; i++ {
		station := stations[rand.Intn(len(stations))]

		_, err := w.WriteString(station.NewRandomMeasure().String() + "\n")
		if err != nil {
			return err
		}

		*c = *c + 1
	}

	err := w.Flush()
	if err != nil {
		return err
	}

	return nil
}
