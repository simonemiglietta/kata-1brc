package value_objects

import (
	"fmt"
	"math"
)

type StationAggregate struct {
	Station     string
	ItemCount   uint64
	Temperature struct {
		Minimum, Maximum, Sum float64
	}
}

func NewStationAggregateFromDetection(d Detection) (sa StationAggregate) {
	sa.Station = d.Station
	sa.ItemCount = 1
	sa.Temperature.Sum = float64(d.Temperature)
	sa.Temperature.Minimum = float64(d.Temperature)
	sa.Temperature.Maximum = float64(d.Temperature)

	return
}

func (sa *StationAggregate) AddDetection(d Detection) {
	t := float64(d.Temperature)

	sa.ItemCount++
	sa.Temperature.Sum += t
	sa.Temperature.Minimum = min(sa.Temperature.Minimum, t)
	sa.Temperature.Maximum = max(sa.Temperature.Minimum, t)
}

func (sa *StationAggregate) TemperatureMean() float64 {
	return sa.Temperature.Sum / float64(sa.ItemCount)
}

func (sa *StationAggregate) String() string {
	roundedMean := math.Round(sa.TemperatureMean()*10) / 10
	return fmt.Sprintf(
		"%s=%.1f/%.1f/%.1f",
		sa.Station, sa.Temperature.Minimum, roundedMean, sa.Temperature.Maximum,
	)
}
