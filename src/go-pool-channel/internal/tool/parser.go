package tool

import (
	"bufio"
	"lvciot/go-pool-channel/internal/model"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	ChanDetectionSize = 10
	ChanAggregateSize = 5
	FinalProcesses    = 2
)

func Parser(sf string, df string, c *atomic.Uint32) {
	nCpus := runtime.NumCPU()
	rows := make(chan string, ChanDetectionSize)
	aggregateMaps := make(chan StationMap, ChanAggregateSize)
	wg := sync.WaitGroup{}

	for i := 0; i < nCpus; i++ {
		wg.Add(1)
		go aggregator(rows, aggregateMaps, &wg)
	}

	stationsMapSafe := NewStationsMapSafe()
	for i := 0; i < FinalProcesses; i++ {
		go func() {
			superAggregator(&stationsMapSafe, aggregateMaps)
		}()
	}

	_ = fileRowsScanner(sf, rows)
	close(rows)

	dstFile, _ := os.Create(df)
	defer dstFile.Close()
	dstWriter := bufio.NewWriter(dstFile)

	wg.Wait()

	// todo: use fileResultsWriter
	stationsMap := stationsMapSafe.M

	totalStations := len(stationsMap)
	stations := make([]string, totalStations)
	aggregateRows := make([]string, totalStations)

	j := 0
	for station, _ := range stationsMap {
		stations[j] = station
		j++
	}
	sort.Strings(stations)
	for j, station := range stations {
		aggregate := stationsMap[station].A
		aggregateRows[j] = aggregate.String()
	}

	_, _ = dstWriter.WriteString("{")
	_, _ = dstWriter.WriteString(strings.Join(aggregateRows, ", "))
	_, _ = dstWriter.WriteString("}\n")
	_ = dstWriter.Flush()
}

func fileRowsScanner(file string, rows chan string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	wg := sync.WaitGroup{}
	defer wg.Wait()

	s := bufio.NewScanner(f)
	for s.Scan() {
		wg.Add(1)

		go func() {
			rows <- s.Text()
			wg.Done()
		}()
	}

	return nil
}

func fileResultsWriter(file string, stations StationMap) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	sortedRows := stations.ToSortedRows()

	w := bufio.NewWriter(f)

	_, _ = w.WriteString("{")
	_, _ = w.WriteString(strings.Join(sortedRows, ", "))
	_, _ = w.WriteString("}\n")
	_ = w.Flush()

	return nil
}

func aggregator(rows <-chan string, aggregateMaps chan StationMap, wg *sync.WaitGroup) {
	aggregateMap := make(StationMap)

	for row := range rows {
		detection := model.NewDetectionFromRow(row)
		a, exist := aggregateMap[detection.StationName]

		if exist {
			a.AddDetection(detection)
		} else {
			a := model.NewStationAggregateFromDetection(detection)
			aggregateMap[detection.StationName] = &a
		}
	}

	aggregateMaps <- aggregateMap
	wg.Done()
}

func superAggregator(aggregateMap *StationsMapSafe, aggregateMaps <-chan StationMap) {
	for am := range aggregateMaps {
		for _, a := range am {
			aggregateMap.AddAggregate(*a)
		}
	}
}
