package models

import (
	"maps"
	"slices"
)

const StationsTotalGuessed int = 10_000

type StationMap map[string]*Station

func NewStationMap() *StationMap {
	sm := make(StationMap, StationsTotalGuessed)
	return &sm
}

func (sm *StationMap) AddStation(s *Station) {
	as, exist := (*sm)[s.Name]

	if exist {
		as.AddStation(s)
	} else {
		(*sm)[s.Name] = s
	}
}

func (sm *StationMap) AddDetection(d Detection) {
	s, exist := (*sm)[d.StationName]

	if exist {
		s.AddDetection(&d)
	} else {
		(*sm)[d.StationName] = NewStationFromDetection(&d)
	}
}

func (sm *StationMap) AddStations(osm *StationMap) {
	for _, station := range *osm {
		(*sm).AddStation(station)
	}
}

func (sm *StationMap) ToSortedRows() []string {
	totalStations := len(*sm)
	sortedRows := make([]string, totalStations)

	stationNames := maps.Keys(*sm)
	for i, stationName := range slices.Sorted(stationNames) {
		sortedRows[i] = (*sm)[stationName].String()
	}

	return sortedRows
}
