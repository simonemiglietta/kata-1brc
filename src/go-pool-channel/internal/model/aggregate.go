package model

import (
	"fmt"
	"math"
)

type Aggregate struct {
	Station     string
	ItemCount   uint64
	Temperature struct {
		Minimum, Maximum, Sum float32
	}
}

func NewStationAggregateFromDetection(d Detection) (sa Aggregate) {
	sa.Station = d.Station
	sa.ItemCount = 1
	sa.Temperature.Sum = d.Temperature
	sa.Temperature.Minimum = d.Temperature
	sa.Temperature.Maximum = d.Temperature

	return
}

func (a *Aggregate) AddDetection(d Detection) {
	t := d.Temperature

	a.ItemCount++
	a.Temperature.Sum += t
	a.Temperature.Minimum = min(a.Temperature.Minimum, t)
	a.Temperature.Maximum = max(a.Temperature.Maximum, t)
}

func (a *Aggregate) AddAggregate(o Aggregate) {
	a.ItemCount += o.ItemCount
	a.Temperature.Sum += o.Temperature.Sum
	a.Temperature.Minimum = min(a.Temperature.Minimum, o.Temperature.Minimum)
	a.Temperature.Maximum = max(a.Temperature.Maximum, o.Temperature.Maximum)
}

func (a *Aggregate) TemperatureMean() float32 {
	return a.Temperature.Sum / float32(a.ItemCount)
}

func (a *Aggregate) String() string {
	roundedMean := math.Round(float64(a.TemperatureMean()*10)) / 10

	return fmt.Sprintf(
		"%s=%.1f/%.1f/%.1f",
		a.Station, a.Temperature.Minimum, roundedMean, a.Temperature.Maximum,
	)
}
