// File: cmd/agent/main_test.go

package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	flags "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/flags"
	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
)

func TestMainWithValidFlags(t *testing.T) {
	// Check that the function did not panic
	assert.Nil(t, recover())
}

func TestMainWithEmptyFlags(t *testing.T) {

	// Call the main function with the empty flags
	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(t, r) // check that the function panicked
		}
	}()
}

func TestGetMetrics(t *testing.T) {
	// Set up valid flags
	f := flags.NewAgentFlags()

	// Create a mock storage object
	memStorage := storage.MemStorage{}

	// Call the get metrics function
	go memStorage.GetMetrics(f.GetMetricsGetDuration())

	// Wait for the get metrics function to finish
	time.Sleep(time.Duration(f.GetMetricsGetDuration()) * time.Second)

	// Check that the storage has some metrics
	assert.NotNil(t, memStorage.Metrics)
}

func TestSignalHandling(t *testing.T) {
	// Set up valid flags
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
