package storage

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetStorage(t *testing.T) {
	// Create a new MemStorage instance
	m := &MemStorage{}

	// Test that SetStorage sets the expected value
	randVal := rand.Float64()
	m.SetStorage(randVal)
	assert.NotNil(t, m.Metrics.Gauge["Alloc"])
	assert.NotNil(t, m.Metrics.Gauge["BuckHashSys"])
	assert.NotNil(t, m.Metrics.Gauge["Frees"])
	assert.NotNil(t, m.Metrics.Gauge["GCCPUFraction"])
	assert.NotNil(t, m.Metrics.Gauge["GCSys"])
	assert.NotNil(t, m.Metrics.Gauge["HeapAlloc"])
	assert.NotNil(t, m.Metrics.Gauge["HeapIdle"])
	assert.NotNil(t, m.Metrics.Gauge["HeapInuse"])
	assert.NotNil(t, m.Metrics.Gauge["HeapObjects"])
	assert.NotNil(t, m.Metrics.Gauge["HeapReleased"])
	assert.NotNil(t, m.Metrics.Gauge["HeapSys"])
	assert.NotNil(t, m.Metrics.Gauge["LastGC"])
	assert.NotNil(t, m.Metrics.Gauge["Lookups"])
	assert.NotNil(t, m.Metrics.Gauge["MCacheInuse"])
	assert.NotNil(t, m.Metrics.Gauge["MCacheSys"])
	assert.NotNil(t, m.Metrics.Gauge["MSpanInuse"])
	assert.NotNil(t, m.Metrics.Gauge["MSpanSys"])
	assert.NotNil(t, m.Metrics.Gauge["Mallocs"])
	assert.NotNil(t, m.Metrics.Gauge["NextGC"])
	assert.NotNil(t, m.Metrics.Gauge["NumForcedGC"])
	assert.NotNil(t, m.Metrics.Gauge["NumGC"])
	assert.NotNil(t, m.Metrics.Gauge["OtherSys"])
	assert.NotNil(t, m.Metrics.Gauge["PauseTotalNs"])
	assert.NotNil(t, m.Metrics.Gauge["StackInuse"])
	assert.NotNil(t, m.Metrics.Gauge["StackSys"])
	assert.NotNil(t, m.Metrics.Gauge["Sys"])
	assert.NotNil(t, m.Metrics.Gauge["TotalAlloc"])
	assert.NotNil(t, m.Metrics.Gauge["RandomValue"])
}

func TestMemStorage_SetCounter(t *testing.T) {
	// Create a new MemStorage instance
	m := &MemStorage{}

	// Test that SetCounter sets the expected value
	count := int64(10)
	m.SetCounter(count)
	assert.Equal(t, count, m.Metrics.Counter["PollCount"])
}

func TestGetMetrics(t *testing.T) {
	// Create a test MemStorage instance
	m := &MemStorage{}

	// Create a test duration
	duration := 1 * time.Second

	// Start the GetMetrics goroutine
	go m.GetMetrics(duration)

	// Wait for a short period of time to allow the goroutine to run
	time.Sleep(500 * time.Millisecond)

	// Check that the storage has been updated
	assert.NotNil(t, m.Metrics.Gauge["RandomValue"])

	// Wait for the goroutine to exit
	time.Sleep(500 * time.Millisecond)

	// Check that the goroutine has exited
	assert.Equal(t, int64(0), m.Metrics.Counter["PollCount"])
}
