package file

import (
	"fmt"
	"strconv"
	"strings"
)

type DataRow struct {
	Station     string
	Temperature float32
}

func NewDataRow(row string) (dr DataRow) {
	parts := strings.Split(row, ";")

	dr.Station = parts[0]
	t, _ := strconv.ParseFloat(parts[1], 32)
	dr.Temperature = float32(t)

	return
}

func (dr DataRow) String() string {
	return fmt.Sprintf("%s;%.2f", dr.Station, dr.Temperature)
}
