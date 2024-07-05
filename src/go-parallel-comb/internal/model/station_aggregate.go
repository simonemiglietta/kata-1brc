package model

import (
	"fmt"
	"math"
)

type StationAggregate struct {
	Station     string `json:"station"`
	ItemCount   uint64 `json:"item_count"`
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

// use composite pattern to add measurement
func (sa *StationAggregate) AddMeasurement(m MeasurementInterface) {
	if m.GetStation() != sa.Station {
		return
	}
	sa.ItemCount += m.GetItemCount()
	sa.Temperature.Sum += m.GetTemperatureSum()
	sa.Temperature.Minimum = min(sa.Temperature.Minimum, m.GetTemperatureMin())
	sa.Temperature.Maximum = max(sa.Temperature.Maximum, m.GetTemperatureMax())
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

func (sa *StationAggregate) GetStation() string {
	return sa.Station
}

func (sa *StationAggregate) GetItemCount() uint64 {
	return sa.ItemCount
}

func (sa *StationAggregate) GetTemperatureSum() float32 {
	return sa.Temperature.Sum
}

func (sa *StationAggregate) GetTemperatureMin() float32 {
	return sa.Temperature.Minimum
}

func (sa *StationAggregate) GetTemperatureMax() float32 {
	return sa.Temperature.Maximum
}
