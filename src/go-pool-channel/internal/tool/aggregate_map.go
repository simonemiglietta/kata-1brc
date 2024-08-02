package tool

import (
	"lvciot/go-pool-channel/internal/model"
)

type AggregateMap map[string]*model.Aggregate

func NewAggregateMap() AggregateMap {
	return make(AggregateMap)
}

func (am AggregateMap) AddAggregate(a *model.Aggregate) {
	as, exist := am[a.Station]

	if exist {
		as.AddAggregate(*a)
	} else {
		am[a.Station] = a
	}
}
