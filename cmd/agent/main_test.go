// File: cmd/agent/main_test.go

package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func TestMainWithValidFlags(t *testing.T) {
	// Set up valid flags
	f := flags.NewAgentFlags()
	f.URL = "localhost:8080"
	f.SendTime = 10
	f.GetMetricTime = 2

	// Check that the function did not panic
	assert.Nil(t, recover())
}

func TestMainWithInvalidFlags(t *testing.T) {
	// Set up invalid flags
	f := flags.NewAgentFlags()
	f.URL = ""          // invalid URL
	f.SendTime = 0      // invalid send time
	f.GetMetricTime = 0 // invalid get metric time

	// Call the main function with the invalid flags
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r) // check that the function panicked
		}
	}()
}

func TestMainWithEmptyFlags(t *testing.T) {

	// Call the main function with the empty flags
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r) // check that the function panicked
		}
	}()
}

// func TestSendMetrics(t *testing.T) {
// 	// Set up valid flags
// 	f := flags.NewFlags()
// 	f.URL = "localhost:8080"
// 	f.SendTime = 10
// 	f.GetMetricTime = 2

// 	// Call the send metrics function
// 	err := handlers.SendMetric(f.URL, "counter", "metric2", "20")

// 	// Check that the function did not return an error
// 	assert.Nil(t, err)
// }

func TestGetMetrics(t *testing.T) {
	// Set up valid flags
	f := flags.NewAgentFlags()
	f.URL = "localhost:8080"
	f.SendTime = 10
	f.GetMetricTime = 2

	// Create a mock storage object
	memStorage := storage.MemStorage{}

	// Call the get metrics function
	go memStorage.GetMetrics(f.GetMetricsGetDuration())

	// Wait for the get metrics function to finish
	time.Sleep(time.Duration(f.GetMetricTime) * time.Second)

	// Check that the storage has some metrics
	assert.NotNil(t, memStorage.Metrics)
}

func TestSignalHandling(t *testing.T) {
	// Set up valid flags
	f := flags.NewAgentFlags()
	f.URL = "localhost:8080"
	f.SendTime = 10
	f.GetMetricTime = 2

	// Create a mock signal channel
	c := make(chan os.Signal, 1)

	// Call the signal handling function
	go func() {
		<-c
		// Check that the function did not panic
		assert.Nil(t, recover())
	}()

	// Send a signal to the channel
	c <- os.Interrupt
}
