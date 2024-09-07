package tool

import (
	"lvciot/go-pool-channel/internal/model"
	"sync"
)

type StationsMapSafe struct {
	mu sync.RWMutex
	M  map[string]*AggregateSafe
}

func NewStationsMapSafe() StationsMapSafe {
	return StationsMapSafe{M: make(map[string]*AggregateSafe)}
}

func (ams *StationsMapSafe) AddAggregate(a model.Station) {
	ams.mu.RLock()

	as, exist := ams.M[a.Name]

	if exist {
		ams.mu.RUnlock()
		as.AddAggregate(a)
	} else {
		as := NewAggregateSafe()
		ams.mu.Lock()
		ams.M[a.Name] = &as
		ams.mu.Unlock()
	}
}
