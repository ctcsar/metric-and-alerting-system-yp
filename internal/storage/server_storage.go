package storage

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	storageInstance *Storage
	once            sync.Once
)

type Storage struct {
	Gauge   map[string]float64
	Counter int64
}

func GetStorageInstance() *Storage {
	once.Do(func() {
		storageInstance = &Storage{
			Gauge: make(map[string]float64),
		}
	})
	return storageInstance
}

func (g *Storage) String() string {
	return fmt.Sprintf("Storage{gauge: %+v, counter: %d}", g.Gauge, g.Counter)
}
func (m *Storage) GetMetricValue(metricType string, metricName string) (float64, bool) {
	if metricType == "gauge" {
		value, ok := m.Gauge[metricName]
		return value, ok
	} else if metricType == "counter" {
		return float64(m.Counter), true
	}
	return 0, false
}

func (m *Storage) GetAllMetrics() map[string]map[string]float64 {
	metrics := make(map[string]map[string]float64)
	metrics["gauge"] = m.Gauge
	metrics["counter"] = map[string]float64{"counter": float64(m.Counter)}
	return metrics
}
func (g *Storage) SetStorage(v, t, n string) {
	switch t {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		g.Gauge[n] = val
	case "counter":
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		g.Counter = val
	}
}
