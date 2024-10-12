package tool

import (
	"bufio"
	"errors"
	"lvciot/go-pool-channel/internal/models"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	ChanDetectionSize = 10
	ChanAggregateSize = 5
	FinalProcesses    = 2
)

func Parser(sf string, df string, counter *atomic.Uint32) {
	nCpus := runtime.NumCPU()
	rows := make(chan string, ChanDetectionSize)

	wg, stationMaps := newWorkersPool(nCpus, rows, counter)

	_ = fileRowsScanner(sf, rows)

	close(rows)
	wg.Wait()

	stationMap := aggregateStationMaps(stationMaps)

	_ = fileResultsWriter(df, stationMap)
}

func fileRowsScanner(file string, rows chan string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		err = errors.Join(err, e)
	}()

	wg := sync.WaitGroup{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		wg.Add(1)

		go func(t string) {
			rows <- t
			wg.Done()
		}(s.Text())
	}

	wg.Wait()
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
