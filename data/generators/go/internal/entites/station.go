package models

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type value struct {
	Mean float64
	Min  float64
	Max  float64
}

type Station struct {
	Name  string
	value value
}

const (
	valueMin       float64 = -99.9
	valueMax       float64 = 99.9
	valueTolerance float64 = 25
)

func NewStationFromString(s string) (Station, error) {
	tokens := strings.Split(s, ";")

	if len(tokens) != 2 {
		return Station{}, fmt.Errorf("invalid station format: %s", s)
	}

	name := tokens[0]
	valueMean, e := strconv.ParseFloat(tokens[1], 32)

	if e != nil {
		return Station{}, fmt.Errorf("invalid measurement value: %s", tokens[1])
	}

	return Station{
		Name: name,
		value: value{
			Mean: valueMean,
			Min:  math.Max(valueMin, valueMean-valueTolerance),
			Max:  math.Min(valueMax, valueMean+valueTolerance),
		},
	}, nil
}

func (s Station) NewRandomMeasure() Measure {
	rn := rand.Float64()

	value := rn*(s.value.Max-s.value.Min) + s.value.Min

	return Measure{StationName: s.Name, Value: value}
}
