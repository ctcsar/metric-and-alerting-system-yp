// server_handler_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricValueHandler(t *testing.T) {
	// Create a test storage
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()

	// Create a test router
	r := chi.NewRouter()
	Webhook(r, g, c)

	// Create a test request
	req, err := http.NewRequest("GET", "/value/gauge/test", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllMetricsHandler(t *testing.T) {
	// Create a test storage
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()

	// Create a test router
	r := chi.NewRouter()
	Webhook(r, g, c)

	// Create a test request
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the response body
}

func TestUpdateHandler(t *testing.T) {
	// Create a test storage
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()

	// Create a test router
	r := chi.NewRouter()
	Webhook(r, g, c)

	// Create a test request
	req, err := http.NewRequest("POST", "/update/gauge/test/10.0", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestUpdateHandler_InvalidMetricType(t *testing.T) {
	// Create a test storage
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()

	// Create a test router
	r := chi.NewRouter()
	Webhook(r, g, c)

	// Create a test request
	req, err := http.NewRequest("POST", "/update/unknown/test/10.0", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestRun(t *testing.T) {
	// Create a test storage
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()

	// Create a test router
	r := chi.NewRouter()
	Webhook(r, g, c)

	// Start the server
	go func() {
		err := Run(r, g, c)
		assert.NoError(t, err)
	}()
	req, err := http.NewRequest("POST", "/update/gauge/test/10.0", nil)
	assert.NoError(t, err)
	// Create a test request
	// req, err := http.NewRequest("GET", "/value/gauge/test", nil)
	// assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)
}
