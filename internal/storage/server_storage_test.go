// server_storage_test.go
package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStorageInstance(t *testing.T) {
	// Test that GetStorageInstance returns a non-nil Storage instance
	instance := GetStorageInstance()
	assert.NotNil(t, instance)

	// Test that GetStorageInstance returns the same instance on subsequent calls
	instance2 := GetStorageInstance()
	assert.Equal(t, instance, instance2)
}

func TestStorage_GetAllMetrics(t *testing.T) {
	// Create a new Storage instance
	m := &Storage{
		Gauge: map[string]float64{
			"test": 10.0,
		},
		Counter: 0,
	}

	// Test that GetAllMetrics returns the expected metrics
	metrics := m.GetAllMetrics()
	assert.Equal(t, map[string]map[string]float64{
		"gauge": {
			"test": 10.0,
		},
		"counter": {
			"counter": 0,
		},
	}, metrics)
}

func TestStorage_GetMetricValue(t *testing.T) {
	// Create a new Storage instance
	m := &Storage{
		Gauge: map[string]float64{
			"test": 10.0,
		},
		Counter: 0,
	}

	// Test that GetMetricValue returns the expected value for a gauge metric
	value, ok := m.GetMetricValue("gauge", "test")
	assert.Equal(t, 10.0, value)
	assert.True(t, ok)

	// Test that GetMetricValue returns the expected value for a counter metric
	value, ok = m.GetMetricValue("counter", "")
	assert.Equal(t, 0.0, value)
	assert.True(t, ok)

	// Test that GetMetricValue returns an error for an unknown metric type
	value, ok = m.GetMetricValue("unknown", "test")
	assert.Equal(t, 0.0, value)
	assert.False(t, ok)
}

func TestStorage_SetStorage(t *testing.T) {
	// Create a new Storage instance
	m := &Storage{
		Gauge: map[string]float64{},
	}

	// Test that SetStorage sets the expected value
	m.SetStorage("10.0", "", "")
	assert.Equal(t, map[string]float64{
		"": 10.0,
	}, m.Gauge)
}
