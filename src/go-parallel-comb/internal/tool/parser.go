package tool

import (
	"bufio"
	"fmt"
	"lvciot/go-seq/internal/model"
	"os"
)

const (
	RowSeparator = "\n"
)

func Parser(processId int, sf string, start int64, end int64, advancement model.AdvancementMutex) map[string]*model.StationAggregate {

	aggregates := make(map[string]*model.StationAggregate)

	srcFile, _ := os.Open(sf)
	defer srcFile.Close()

	start, err := findLineStartPosition(srcFile, start)
	if err != nil {
		fmt.Println("Error finding line start position:", err)
		return nil
	}

	_, err = srcFile.Seek(start, 0)
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
			a.AddMeasurement(d)
		} else {
			a := model.NewStationAggregateFromDetection(d)
			aggregates[d.Station] = &a
		}

		updateAdvancement(processId, advancement)

		position += int64(len([]byte(text)) + len(RowSeparator))
		if position >= end {
			break
		}
	}
	return aggregates
}

func updateAdvancement(processId int, advancement model.AdvancementMutex) {
	advancement.ShardLocks[processId].Lock()
	advancement.Shards[processId]++
	advancement.ShardLocks[processId].Unlock()
}

func findLineStartPosition(file *os.File, start int64) (int64, error) {
	buffer := make([]byte, 1)
	for {
		if start == 0 {
			break
		}
		start--
		_, err := file.Seek(start, 0)
		if err != nil {
			return 0, err
		}
		_, err = file.Read(buffer)
		if err != nil {
			return 0, err
		}
		if buffer[0] == '\n' {
			start++
			break
		}
	}
	return start, nil
}
