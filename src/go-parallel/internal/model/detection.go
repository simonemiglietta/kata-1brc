package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Detection struct {
	Station     string
	Temperature float32
}

func NewDetectionFromRow(row string) (d Detection) {
	split := strings.Split(row, ";")

	d.Station = split[0]
	t, _ := strconv.ParseFloat(split[1], 32)
	d.Temperature = float32(t)

	return
}

func (d Detection) String() string {
	return fmt.Sprintf("%s;%.1f", d.Station, d.Temperature)
}
