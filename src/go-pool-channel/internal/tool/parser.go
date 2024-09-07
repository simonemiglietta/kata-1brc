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

	aggregateMapSafe := NewStationsMapSafe()
	for i := 0; i < FinalProcesses; i++ {
		go func() {
			superAggregator(&aggregateMapSafe, aggregateMaps)
		}()
	}

	srcFile, _ := os.Open(sf)
	defer srcFile.Close()
	srcScanner := bufio.NewScanner(srcFile)

	dstFile, _ := os.Create(df)
	defer dstFile.Close()
	dstWriter := bufio.NewWriter(dstFile)

	for srcScanner.Scan() {
		c.Add(1)
		rows <- srcScanner.Text()
	}
	close(rows)

	wg.Wait()

	aggregatesSafe := aggregateMapSafe.M

	totalStations := len(aggregatesSafe)
	stations := make([]string, totalStations)
	aggregateRows := make([]string, totalStations)

	j := 0
	for station, _ := range aggregatesSafe {
		stations[j] = station
		j++
	}
	sort.Strings(stations)
	for j, station := range stations {
		aggregate := aggregatesSafe[station].A
		aggregateRows[j] = aggregate.String()
	}

	_, _ = dstWriter.WriteString("{")
	_, _ = dstWriter.WriteString(strings.Join(aggregateRows, ", "))
	_, _ = dstWriter.WriteString("}\n")
	_ = dstWriter.Flush()

	return
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
