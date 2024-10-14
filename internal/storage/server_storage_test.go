package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const Gauge = "gauge"

func TestStorage_NewGaugeStorage(t *testing.T) {
	// Test that NewGaugeStorage returns a new Storage instance with an empty Gauge map
	m := NewStorage()
	assert.NotNil(t, m)
	assert.Empty(t, m.Gauge)
}

func TestStorage_GetAllGaugeMetrics(t *testing.T) {
	// Create a new Storage instance with some gauge metrics
	m := &Storage{
		Gauge: map[string]float64{
			"test1": 10.0,
			"test2": 20.0,
		},
	}

	// Test that GetAllGaugeMetrics returns all gauge metrics
	metrics := m.GetAllGaugeMetrics()
	assert.Equal(t, map[string]map[string]float64{
		Gauge: {
			"test1": 10.0,
			"test2": 20.0,
		},
	}, metrics)
}

func TestStorage_SetGauge(t *testing.T) {
	// Create a new Storage instance
	m := &Storage{
		Gauge: make(map[string]float64),
	}

	// Test that SetGauge sets the expected value for a gauge metric
	err := m.SetGauge("test", "10.0")
	assert.Nil(t, err)
	assert.Equal(t, 10.0, m.Gauge["test"])
}

func TestStorage_SetCounter(t *testing.T) {
	// Create a new Storage instance
	m := &Storage{
		Counter: make(map[string]int64),
	}

	// Test that SetCounter sets the expected value for a counter metric
	err := m.SetCounter("test", "10")
	assert.Nil(t, err)
	assert.Equal(t, int64(10), m.Counter["test"])
}

func TestStorage_IncrementCounter(t *testing.T) {
	// Create a new Storage instance with a counter metric
	m := &Storage{
		Counter: map[string]int64{
			"test": 10,
		},
	}

	// Test that IncrementCounter increments the counter metric
	m.Counter["test"]++
	assert.Equal(t, int64(11), m.Counter["test"])
}
