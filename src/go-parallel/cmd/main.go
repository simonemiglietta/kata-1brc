package main

import (
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"lvciot/go-seq/internal/model"
	"lvciot/go-seq/internal/tool"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

const (
	MaxRows = 1_000_000
	SrcFile = "../../../../data/measurements.txt"
	DstFile = "../../measurements.out"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	srcFile := filepath.Join(b, SrcFile)

	bar := progressbar.Default(MaxRows)
	ticker := time.NewTicker(time.Second)

	t := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = bar.Set(t)
			}
		}
	}()

	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	numCores := 1

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	fileSize := fileInfo.Size()

	// Calculate the size of each partition
	partitionSize := fileSize / int64(numCores)

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
			results[i] = tool.Parser(srcFile, start, end, &t)
		}(i)
	}

	wg.Wait()
	aggregates := make(map[string]*model.StationAggregate)

	for _, result := range results {
		a := result
		fmt.Print(a)
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

	_ = bar.Set(MaxRows)
	_, _ = fmt.Print(json.Marshal(aggregates))
}
