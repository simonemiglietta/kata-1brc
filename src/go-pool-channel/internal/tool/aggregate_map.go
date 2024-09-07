package tool

import (
	"lvciot/go-pool-channel/internal/model"
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
