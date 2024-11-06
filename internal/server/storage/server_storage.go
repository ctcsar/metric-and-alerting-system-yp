package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
	mutex   sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *Storage) String() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return fmt.Sprintf("Storage{gauge: %+v, counter: %+v}", m.Gauge, m.Counter)
}

func (m *Storage) GetAllMertrics() *Storage {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	metrics := &Storage{}
	metrics.Gauge = m.Gauge
	metrics.Counter = m.Counter
	return metrics
}

func (m *Storage) SetGauge(key string, val float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Gauge[key] = val
	return nil
}

func (m *Storage) SetCounter(key string, val int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Counter[key] = m.Counter[key] + val
	return nil
}
func (m *Storage) GetGaugeValue(metricName string) (float64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Gauge[metricName]
	return value, ok
}

func (m *Storage) GetCounterValue(metricName string) (int64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.Counter[metricName]
	return value, ok
}

func (m *Storage) GetAllGaugeMetrics() map[string]map[string]float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	metrics := make(map[string]map[string]float64)
	metrics["gauge"] = m.Gauge
	return metrics
}

func (m *Storage) GetAllCounterMetrics() map[string]map[string]int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	metrics := make(map[string]map[string]int64)
	metrics["counter"] = m.Counter
	return metrics
}

func (m *Storage) SetStorage(metrics *Storage) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Gauge = metrics.Gauge
	m.Counter = metrics.Counter
}
