package tool

import (
	"lvciot/go-pool-channel/internal/model"
	"maps"
	"slices"
)

type StationMap map[string]*model.Station

func NewAggregateMap() StationMap {
	return make(StationMap)
}

func (sm StationMap) AddAggregate(s *model.Station) {
	as, exist := sm[s.Name]

	if exist {
		as.AddStation(*s)
	} else {
		sm[s.Name] = s
	}
}

func (sm StationMap) ToSortedRows() []string {
	totalStations := len(sm)
	sortedRows := make([]string, totalStations)

	stationNames := maps.Keys(sm)
	for i, stationName := range slices.Sorted(stationNames) {
		sortedRows[i] = sm[stationName].String()
	}

	return sortedRows
}
