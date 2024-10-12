package tool

import (
	"lvciot/go-pool-channel/internal/models"
	"math"
	"sync"
	"sync/atomic"
)

func newWorkersPool(workers int, chunks <-chan [RowsChunkSize]string, counter *atomic.Uint32) (*sync.WaitGroup, *[]*models.StationMap) {
	wg := sync.WaitGroup{}
	sm := make([]*models.StationMap, workers)

	wg.Add(workers)

	go func() {
		for i := 0; i < workers; i++ {
			sm[i] = rowsToStationMap(chunks, &wg, counter)
		}
	}()

	return &wg, &sm
}

func rowsToStationMap(chunks <-chan [RowsChunkSize]string, wg *sync.WaitGroup, counter *atomic.Uint32) *models.StationMap {
	defer wg.Done()

	sm := models.NewStationMap()

	for rows := range chunks {
		for _, row := range rows {
			if row == "" {
				// no more rows
				break
			}

			d := models.NewDetectionFromRow(row)
			sm.AddDetection(d)

			go func() {
				counter.Add(1)
			}()
		}
	}

	return sm
}

// aggregateStationMaps merges concurrently an array of StationMap splitting them in couple and reiterating recursively
// until just one StationMap will be left
func aggregateStationMaps(oldArray *[]*models.StationMap) *models.StationMap {
	size := len(*oldArray)

	if size == 1 {
		return (*oldArray)[0]
	}

	newSize := int(math.Ceil(float64(size) / 2))
	newArray := make([]*models.StationMap, newSize)
	wg := sync.WaitGroup{}

	for i := 0; i < newSize; i++ {
		// coupled items indexes of the oldArray
		a := i * 2
		b := (i + 1) * 2

		if b >= newSize {
			// the last one, has no coupled item
			newArray[i] = (*oldArray)[a]
		} else {
			wg.Add(1)
			go func() {
				defer wg.Done()

				sm := (*oldArray)[a]
				sm.AddStations((*oldArray)[a])
				newArray[a] = sm
			}()
		}
	}

	wg.Wait()

	return aggregateStationMaps(&newArray)
}
