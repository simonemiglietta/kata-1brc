package tool

import (
	"bufio"
	"lvciot/go-pool-channel/internal/model"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const (
	ChanDetectionSize = 1_000
	ChanPoolSize      = 5
)

func Parser(sf string, df string, increment chan int) {
	nCpus := runtime.NumCPU()
	rows := make(chan string, ChanDetectionSize)
	aggregateMapPools := make(chan AggregateMap, ChanPoolSize)
	aggregateResultsMap := NewAggregateMap()
	var poolWg sync.WaitGroup

	// pool of workers
	// dividi et impera through channel
	for i := 0; i < nCpus; i++ {
		go aggregator(rows, aggregateMapPools, &poolWg)
	}

	// recompose results
	go func() {
		for aggregateMap := range aggregateMapPools {
			for _, a := range aggregateMap {
				aggregateResultsMap.AddAggregate(a)
			}
		}
	}()

	// file parsing
	srcFile, _ := os.Open(sf)
	defer srcFile.Close()
	srcScanner := bufio.NewScanner(srcFile)

	dstFile, _ := os.Create(df)
	defer dstFile.Close()
	dstWriter := bufio.NewWriter(dstFile)

	for srcScanner.Scan() {
		increment <- 1
		rows <- srcScanner.Text()
	}
	close(rows)
	poolWg.Wait()
	close(aggregateMapPools)

	// sort and print
	totalStations := len(aggregateResultsMap)
	stations := make([]string, totalStations)
	aggregateRows := make([]string, totalStations)

	j := 0
	for station, _ := range aggregateResultsMap {
		stations[j] = station
		j++
	}
	sort.Strings(stations)
	for j, station := range stations {
		aggregate := aggregateResultsMap[station]
		aggregateRows[j] = aggregate.String()
	}

	_, _ = dstWriter.WriteString("{")
	_, _ = dstWriter.WriteString(strings.Join(aggregateRows, ", "))
	_, _ = dstWriter.WriteString("}\n")
	_ = dstWriter.Flush()

	return
}

func aggregator(rows <-chan string, aggregateMaps chan AggregateMap, wg *sync.WaitGroup) {
	wg.Add(1)
	aggregateMap := NewAggregateMap()

	for row := range rows {
		detection := model.NewDetectionFromRow(row)
		a, exist := aggregateMap[detection.Station]

		if exist {
			a.AddDetection(detection)
		} else {
			a := model.NewStationAggregateFromDetection(detection)
			aggregateMap[detection.Station] = &a
		}
	}

	aggregateMaps <- aggregateMap
	wg.Done()
}
