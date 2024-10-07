package storage

// import (
// 	"math/rand/v2"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestMemStorage_String(t *testing.T) {
// 	// Create a new MemStorage instance
// 	m := &MemStorage{
// 		Metrics: Metrics{
// 			Gauge: map[string]float64{
// 				"test": 10.0,
// 			},
// 			Counter: map[string]int64{
// 				"test": 10,
// 			},
// 		},
// 	}

// 	// Test that String returns the expected value
// 	expected := "Metrics: Gauge: map[test:10] Counter: map[test:10]"
// 	assert.Equal(t, expected, m.String())
// }

// func TestMemStorage_SetStorage(t *testing.T) {
// 	// Create a new MemStorage instance
// 	m := &MemStorage{}

// 	// Test that SetStorage sets the expected value
// 	randVal := rand.Float64()
// 	m.SetStorage(randVal)
// 	assert.Equal(t, randVal, m.Metrics.Gauge["RandomValue"])
// }

// func TestMemStorage_SetCounter(t *testing.T) {
// 	// Create a new MemStorage instance
// 	m := &MemStorage{}

// 	// Test that SetCounter sets the expected value
// 	count := int64(10)
// 	m.SetCounter(count)
// 	assert.Equal(t, count, m.Metrics.Counter["counter"])
// }

// func TestMemStorage_GetMetrics(t *testing.T) {
// 	// Create a new MemStorage instance
// 	m := &MemStorage{
// 		Metrics: Metrics{
// 			Gauge: map[string]float64{
// 				"test": 10.0,
// 			},
// 			Counter: map[string]int64{
// 				"test": 10,
// 			},
// 		},
// 	}

// 	// Test that GetMetrics returns the expected value
// 	metrics := m.GetMetrics()
// 	assert.Equal(t, m.Metrics, metrics)
// }
