package tool

import (
	"lvciot/go-pool-channel/internal/model"
	"sync"
)

type AggregateSafe struct {
	mu sync.Mutex
	A  model.Aggregate
}

func NewAggregateSafe() AggregateSafe {
	return AggregateSafe{}
}

func (as *AggregateSafe) AddDetection(d model.Detection) {
	as.mu.Lock()
	as.A.AddDetection(d)
	as.mu.Unlock()
}

func (as *AggregateSafe) AddAggregate(a model.Aggregate) {
	as.mu.Lock()
	as.A.AddAggregate(a)
	as.mu.Unlock()
}
