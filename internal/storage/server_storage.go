package storage

import (
	"fmt"
)

type Storage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *Storage) String() string {
	return fmt.Sprintf("Storage{gauge: %+v, counter: %+v}", m.Gauge, m.Counter)
}

func (m *Storage) SetGauge(key string, val float64) error {
	m.Gauge[key] = val
	return nil
}

func (m *Storage) SetCounter(key string, val int64) error {
	m.Counter[key] = m.Counter[key] + val
	fmt.Println(m.Counter)
	return nil
}
func (m *Storage) GetGaugeValue(metricName string) (float64, bool) {
	value, ok := m.Gauge[metricName]
	return value, ok
}

func (m *Storage) GetCounterValue(metricName string) (int64, bool) {
	value, ok := m.Counter[metricName]
	return value, ok
}

func (m *Storage) GetAllGaugeMetrics() map[string]map[string]float64 {
	metrics := make(map[string]map[string]float64)
	metrics["gauge"] = m.Gauge
	return metrics
}

func (m *Storage) GetAllCounterMetrics() map[string]map[string]int64 {
	metrics := make(map[string]map[string]int64)
	metrics["counter"] = m.Counter
	return metrics
}
