package model

import "sync"

type AdvancementMutex struct {
	ShardLocks []sync.Mutex
	Shards     []int
}

type MeasurementInterface interface {
	GetStation() string
	GetItemCount() uint64
	GetTemperatureSum() float32
	GetTemperatureMin() float32
	GetTemperatureMax() float32
}
