package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Detection struct {
	Station     string
	Temperature float32
}

func NewDetectionFromRow(row string) (d Detection) {
	split := strings.Split(row, ";")
	if len(split) < 2 {
		fmt.Fprintf(os.Stderr, "Invalid row format: %s\n", row)
		return d
	}

	d.Station = split[0]
	t, _ := strconv.ParseFloat(split[1], 32)

	d.Temperature = float32(t)
	return
}

func (d Detection) String() string {
	return fmt.Sprintf("%s;%.1f", d.Station, d.Temperature)
}
