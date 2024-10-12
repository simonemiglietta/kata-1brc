package tool

import (
	"bufio"
	"errors"
	"lvciot/go-pool-channel/internal/models"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

const (
	RowsChanSize  int = 1_000
	RowsChunkSize int = 1_000
)

func Parser(sf string, df string, counter *atomic.Uint32) {
	nCpus := runtime.NumCPU()
	chunks := make(chan [RowsChunkSize]string, RowsChanSize)

	wg, stationMaps := newWorkersPool(nCpus, chunks, counter)

	_ = fileRowsScanner(sf, chunks)

	close(chunks)
	wg.Wait()

	stationMap := aggregateStationMaps(stationMaps)

	_ = fileResultsWriter(df, stationMap)
}

func fileRowsScanner(file string, chunks chan [RowsChunkSize]string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		err = errors.Join(err, e)
	}()

	s := bufio.NewScanner(f)

main:
	for {
		var chunk [RowsChunkSize]string

		for i := 0; i < RowsChunkSize; i++ {
			if !s.Scan() {
				chunks <- chunk
				break main
			}

			chunk[i] = s.Text()
		}

		chunks <- chunk
	}

	return err
}

func fileResultsWriter(file string, stations *models.StationMap) (err error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		err = errors.Join(err, e)
	}()

	sortedRows := stations.ToSortedRows()

	w := bufio.NewWriter(f)

	_, _ = w.WriteString("{")
	_, _ = w.WriteString(strings.Join(sortedRows, ", "))
	_, _ = w.WriteString("}\n")
	_ = w.Flush()

	return nil
}
