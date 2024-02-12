package tool

import (
	"bufio"
	"lvciot/go-seq/internal/model"
	"os"
	"sort"
	"strings"
)

func Parser(sf string, df string, i *int) {
	// todo: make map of reference
	aggregates := make(map[string]*model.StationAggregate)

	srcFile, _ := os.Open(sf)
	defer srcFile.Close()
	srcScanner := bufio.NewScanner(srcFile)

	dstFile, _ := os.Create(df)
	defer dstFile.Close()
	dstWriter := bufio.NewWriter(dstFile)

	for srcScanner.Scan() {
		*i++
		d := model.NewDetectionFromRow(srcScanner.Text())

		a, exist := aggregates[d.Station]

		if exist {
			a.AddDetection(d)
		} else {
			a := model.NewStationAggregateFromDetection(d)
			aggregates[d.Station] = &a
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

	_, _ = dstWriter.WriteString("{")
	_, _ = dstWriter.WriteString(strings.Join(aggregateRows, ", "))
	_, _ = dstWriter.WriteString("}\n")
	_ = dstWriter.Flush()

	return
}
