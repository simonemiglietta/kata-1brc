package tool

import (
	"lvciot/go-pool-channel/internal/model"
	"sync"
)

type AggregateSafe struct {
	mu sync.Mutex
	A  model.Station
}

func NewAggregateSafe() AggregateSafe {
	return AggregateSafe{}
}

func (as *AggregateSafe) AddDetection(d model.Detection) {
	as.mu.Lock()
	as.A.AddDetection(d)
	as.mu.Unlock()
}

func (as *AggregateSafe) AddAggregate(a model.Station) {
	as.mu.Lock()
	as.A.AddStation(a)
	as.mu.Unlock()
}
