package model

import (
	"fmt"
	"math"
)

type StationAggregate struct {
	Station     string
	ItemCount   uint64
	Temperature struct {
		Minimum, Maximum, Sum float32
	}
}

func NewStationAggregateFromDetection(d Detection) (sa StationAggregate) {
	sa.Station = d.Station
	sa.ItemCount = 1
	sa.Temperature.Sum = d.Temperature
	sa.Temperature.Minimum = d.Temperature
	sa.Temperature.Maximum = d.Temperature

	return
}

func (sa *StationAggregate) AddDetection(d Detection) {
	t := d.Temperature

	sa.ItemCount++
	sa.Temperature.Sum += t
	sa.Temperature.Minimum = min(sa.Temperature.Minimum, t)
	sa.Temperature.Maximum = max(sa.Temperature.Maximum, t)
}

func (sa *StationAggregate) TemperatureMean() float32 {
	return sa.Temperature.Sum / float32(sa.ItemCount)
}

func (sa *StationAggregate) String() string {
	roundedMean := math.Round(float64(sa.TemperatureMean()*10)) / 10

	return fmt.Sprintf(
		"%s=%.1f/%.1f/%.1f",
		sa.Station, sa.Temperature.Minimum, roundedMean, sa.Temperature.Maximum,
	)
}
