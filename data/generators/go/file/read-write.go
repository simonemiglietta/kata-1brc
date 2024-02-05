package file

import (
	"bufio"
	"os"
)

func ReadData(fileName string) []DataRow {
	file, _ := os.Open(fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	dataRows := make(map[int]DataRow)
	l := 0

	for scanner.Scan() {
		row := NewDataRow(scanner.Text())

		dataRows[l] = row
		l++
	}

	r := make([]DataRow, l)

	for i := 0; i < l; i++ {
		r[i] = dataRows[i]
	}

	return r
}

func WriteData(fileName string, drs []DataRow) {
	file, _ := os.Create(fileName)
	writer := bufio.NewWriter(file)

	for _, dr := range drs {
		_, _ = writer.WriteString(dr.String() + "\n")
	}
}
