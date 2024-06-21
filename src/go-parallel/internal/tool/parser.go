package tool

import (
	"bufio"
	"fmt"
	"lvciot/go-seq/internal/model"
	"os"
)

func Parser(sf string, start int64, end int64, t *int) map[string]*model.StationAggregate {

	aggregates := make(map[string]*model.StationAggregate)

	srcFile, _ := os.Open(sf)
	defer srcFile.Close()

	_, err := srcFile.Seek(start, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return nil
	}

	srcScanner := bufio.NewScanner(srcFile)
	position := start
	for srcScanner.Scan() {
		text := srcScanner.Text()
		d := model.NewDetectionFromRow(text)

		a, exist := aggregates[d.Station]

		if exist {
			a.AddDetection(d)
		} else {
			a := model.NewStationAggregateFromDetection(d)
			aggregates[d.Station] = &a
		}
		// retrieve \n length
		position += int64(len([]byte(text))) + 1
		*t++
		if position >= end {
			break
		}

	}
	return aggregates
}
