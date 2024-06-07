package tool

import (
	"lvciot/go-pool-channel/internal/model"
	"sync"
)

type AggregateMapSafe struct {
	mu sync.RWMutex
	M  map[string]*AggregateSafe
}

func NewAggregateMapSafe() AggregateMapSafe {
	return AggregateMapSafe{M: make(map[string]*AggregateSafe)}
}

func (ams *AggregateMapSafe) AddAggregate(a model.Aggregate) {
	ams.mu.RLock()

	as, exist := ams.M[a.Station]

	if exist {
		ams.mu.RUnlock()
		as.AddAggregate(a)
	} else {
		as := NewAggregateSafe()
		ams.mu.Lock()
		ams.M[a.Station] = &as
		ams.mu.Unlock()
	}
}
