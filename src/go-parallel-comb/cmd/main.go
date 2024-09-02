package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"lvciot/go-seq/internal/model"
	"lvciot/go-seq/internal/tool"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	MaxRows     = 1_000_000_000
	SrcFile     = "../../../../data/measurements.txt"
	DstFile     = "../../measurements.out"
	DstFileJson = "../../measurements.json"
)

func main() {
	numCores := runtime.NumCPU()
	_, b, _, _ := runtime.Caller(0)
	srcFile := filepath.Join(b, SrcFile)
	dstFile := filepath.Join(b, DstFile)
	dstFileJson := filepath.Join(b, DstFileJson)

	sortedStations := StationStats(srcFile, numCores)

	writeOnFile(sortedStations, dstFile, dstFileJson)
}

func StationStats(srcFile string, numCores int) []string {
	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return nil
	}

	// Calculate the size of each partition
	fileSize := fileInfo.Size()
	partitionSize := fileSize / int64(numCores)

	advancement := model.AdvancementMutex{
		ShardLocks: make([]sync.Mutex, numCores),
		Shards:     make([]int, numCores),
	}

	progressBar(numCores, advancement)

	var wg sync.WaitGroup
	wg.Add(numCores)
	stations := make([]map[string]*model.StationAggregate, numCores)
	for i := 0; i < numCores; i++ {
		go func(i int) {
			defer wg.Done()
			start := int64(i) * partitionSize
			end := start + partitionSize
			if i == numCores-1 {
				end = fileSize // Ensure the last partition goes to the end of the file
			}
			stations[i] = tool.Parser(i, srcFile, start, end, advancement)
		}(i)
	}
	wg.Wait()

	reconstructedStations := reconstructSolution(stations)
	sortedStations := sortStations(reconstructedStations)

	return sortedStations
}

func writeOnFile(sortedStations []string, rawFileDest string, jsonDestFile string) {
	df, _ := os.Create(rawFileDest)
	defer df.Close()
	dstWriter := bufio.NewWriter(df)
	_, _ = dstWriter.WriteString(strings.Join(sortedStations, "\n"))
	_ = dstWriter.Flush()

	dfj, _ := os.Create(jsonDestFile)
	defer df.Close()
	dstWriterJ := bufio.NewWriter(dfj)
	jsonRows, _ := json.Marshal(sortedStations)
	_, _ = dstWriterJ.Write(jsonRows)
	_ = dstWriterJ.Flush()
}

func sortStations(aggregates map[string]*model.StationAggregate) []string {
	j := 0
	totalStations := len(aggregates)
	stations := make([]string, totalStations)
	sortedRows := make([]string, totalStations)

	for station, _ := range aggregates {
		stations[j] = station
		j++
	}
	sort.Strings(stations)
	for j, station := range stations {
		aggregate := aggregates[station]
		sortedRows[j] = aggregate.String()
	}
	return sortedRows

}

func progressBar(numCores int, advancement model.AdvancementMutex) {
	bar := progressbar.Default(MaxRows)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:

				totalProgress := 0
				for i := 0; i < numCores; i++ {
					totalProgress += advancement.Shards[i]
				}
				_ = bar.Set(totalProgress)
			}
		}
	}()
}

func reconstructSolution(processesStations []map[string]*model.StationAggregate) map[string]*model.StationAggregate {
	aggregates := make(map[string]*model.StationAggregate)

	for _, stationsAggregates := range processesStations {
		for _, sa := range stationsAggregates {
			a, exist := aggregates[sa.Station]
			if exist {
				a.AddMeasurement(sa)
			} else {
				aggregates[sa.Station] = sa
			}
		}
	}
	return aggregates
}
