// agent_storage_test.go
package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetCounter(t *testing.T) {
	// Create a new MemStorage instance
	storage := MemStorage{
		Metrics: Metrics{
			Counter: 0,
		},
	}

	// Test that SetCounter sets the expected value
	storage.SetCounter(10)
	assert.Equal(t, int64(10), storage.Metrics.Counter)
}

func TestMemStorage_SetStorage(t *testing.T) {
	// Create a new MemStorage instance
	storage := MemStorage{
		Metrics: Metrics{
			Gauge: map[string]float64{},
		},
	}

	// Test that SetStorage sets the expected value
	storage.SetStorage(10.0)
	assert.NotNil(t, storage.Metrics.Gauge)
}

func TestMemStorage_GetMetrics(t *testing.T) {
	// Create a new MemStorage instance
	storage := MemStorage{
		Metrics: Metrics{
			Counter: 0,
		},
	}

	// Test that GetMetrics returns a MemStorage instance
	metrics := storage.GetMetrics()
	assert.NotNil(t, metrics)
}

func TestMemStorage_String(t *testing.T) {
	// Create a new MemStorage instance
	storage := MemStorage{
		Metrics: Metrics{
			Counter: 0,
		},
	}

	// Test that String returns a non-empty string
	str := storage.String()
	assert.NotEmpty(t, str)
}
