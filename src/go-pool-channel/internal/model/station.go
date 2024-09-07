package model

import (
	"fmt"
	"math"
)

type Station struct {
	Name         string
	SamplesCount uint64
	Temperature  struct {
		Minimum, Maximum, Sum float32
	}
}

func NewStationAggregateFromDetection(d Detection) (sa Station) {
	sa.Name = d.StationName
	sa.SamplesCount = 1
	sa.Temperature.Sum = d.Temperature
	sa.Temperature.Minimum = d.Temperature
	sa.Temperature.Maximum = d.Temperature

	return
}

func (s *Station) AddDetection(d Detection) {
	t := d.Temperature

	s.SamplesCount++
	s.Temperature.Sum += t
	s.Temperature.Minimum = min(s.Temperature.Minimum, t)
	s.Temperature.Maximum = max(s.Temperature.Maximum, t)
}

func (s *Station) AddStation(o Station) {
	s.SamplesCount += o.SamplesCount
	s.Temperature.Sum += o.Temperature.Sum
	s.Temperature.Minimum = min(s.Temperature.Minimum, o.Temperature.Minimum)
	s.Temperature.Maximum = max(s.Temperature.Maximum, o.Temperature.Maximum)
}

func (s *Station) TemperatureMean() float32 {
	return s.Temperature.Sum / float32(s.SamplesCount)
}

func (s *Station) String() string {
	roundedMean := math.Round(float64(s.TemperatureMean()*10)) / 10

	return fmt.Sprintf(
		"%s=%.1f/%.1f/%.1f",
		s.Name, s.Temperature.Minimum, roundedMean, s.Temperature.Maximum,
	)
}
