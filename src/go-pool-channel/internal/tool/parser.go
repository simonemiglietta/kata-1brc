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
	CHAN_DETECTION_SIZE = 10
	CHAN_AGGREGATE_SIZE = 5
	FINAL_PROCESSES     = 2
)

type AggregateMap map[string]*model.Aggregate

func Parser(sf string, df string, c *int) {
	nCpus := runtime.NumCPU()
	rows := make(chan string, CHAN_DETECTION_SIZE)
	aggregateMaps := make(chan AggregateMap, CHAN_AGGREGATE_SIZE)
	wg := sync.WaitGroup{}

	for i := 0; i < nCpus; i++ {
		wg.Add(1)
		go aggregator(rows, aggregateMaps, &wg)
	}

	aggregateMapSafe := NewAggregateMapSafe()
	for i := 0; i < FINAL_PROCESSES; i++ {
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
		*c++
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

func aggregator(rows <-chan string, aggregateMaps chan AggregateMap, wg *sync.WaitGroup) {
	aggregateMap := make(AggregateMap)

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

func superAggregator(aggregateMap *AggregateMapSafe, aggregateMaps <-chan AggregateMap) {
	for am := range aggregateMaps {
		for _, a := range am {
			aggregateMap.AddAggregate(*a)
		}
	}
}
