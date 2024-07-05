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
	_, b, _, _ := runtime.Caller(0)
	srcFile := filepath.Join(b, SrcFile)
	dstFile := filepath.Join(b, DstFile)
	numCores := 16

	bar := progressbar.Default(MaxRows)
	ticker := time.NewTicker(time.Second)

	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	df, _ := os.Create(dstFile)
	defer df.Close()
	dstWriter := bufio.NewWriter(df)

	dfj, _ := os.Create(DstFileJson)
	defer df.Close()

	dstWriterJ := bufio.NewWriter(dfj)

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Calculate the size of each partition
	partitionSize := fileSize / int64(numCores)

	advancement := model.AdvancementMutex{
		ShardLocks: make([]sync.Mutex, numCores),
		Shards:     make([]int, numCores),
	}

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

	var wg sync.WaitGroup
	wg.Add(numCores)

	results := make([]map[string]*model.StationAggregate, numCores)

	for i := 0; i < numCores; i++ {
		go func(i int) {
			defer wg.Done()
			start := int64(i) * partitionSize
			end := start + partitionSize
			if i == numCores-1 {
				end = fileSize // Ensure the last partition goes to the end of the file
			}
			results[i] = tool.Parser(i, srcFile, start, end, advancement)
		}(i)
	}

	wg.Wait()
	aggregates := make(map[string]*model.StationAggregate)

	for _, result := range results {
		for _, station := range result {
			a, exist := aggregates[station.Station]
			if exist {
				a.AddMeasurement(station)
			} else {
				aggregates[station.Station] = station
			}
		}
	}

	totalStations := len(aggregates)
	stations := make([]string, totalStations)
	aggregateRows := make([]string, totalStations)

	j := 0
	for station, _ := range aggregates {
		stations[j] = station
		j++
	}
	sort.Strings(stations)
	for j, station := range stations {
		aggregate := aggregates[station]
		aggregateRows[j] = aggregate.String()
	}

	_, _ = dstWriter.WriteString(strings.Join(aggregateRows, "\n"))
	_ = dstWriter.Flush()

	jsonRows, err := json.Marshal(aggregateRows)
	_, _ = dstWriterJ.Write(jsonRows)
	_ = dstWriterJ.Flush()
}
