package telemetry

import (
	"scheduler/pkg/types"
	"sync"
)

// this is the telemetry data cache
// it has a read write mutex to allow multiple
// reads but only one write at a time
type TelemetryCache struct {
	sync.RWMutex
	data map[string]types.NodeTelemetryMetrics
}

// create a global instance of the cache
var globalCache = &TelemetryCache{
	data: make(map[string]types.NodeTelemetryMetrics),
}

func UpdateCache(newData map[string]types.NodeTelemetryMetrics) {
	globalCache.Lock() // full lock to write
	defer globalCache.Unlock()
	globalCache.data = newData
}

func GetNodeMetrics(nodeName string) (types.NodeTelemetryMetrics, bool) {
	globalCache.RLock() // read lock
	defer globalCache.RUnlock()
	val, ok := globalCache.data[nodeName]
	return val, ok
}
