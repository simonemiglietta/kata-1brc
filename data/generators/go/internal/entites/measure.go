package models

import "fmt"

type Measure struct {
	StationName string
	Value       float64
}

func (m Measure) String() string {
	return fmt.Sprintf("%s;%.1f", m.StationName, m.Value)
}
