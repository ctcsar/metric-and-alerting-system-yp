// server_handler_test.go
package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

func dbConnect() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), "host=localhost user=metrics password=password dbname=metrics")
	if err != nil {
		panic(err)
	}
	return db
}
func TestGetMetricValueHandler(t *testing.T) {
	// Create a test storage
	m := storage.NewStorage()
	// Create a test router
	r := chi.NewRouter()
	f := dbConnect()
	Routers(r, m, f)

	// Add a test metric value to the storage
	m.Gauge["test"] = 10.5

	// Create a test request
	req, err := http.NewRequest("GET", "/value/gauge/test", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllMetricsHandler(t *testing.T) {
	// Create a test storage
	m := storage.NewStorage()
	// Create a test router
	r := chi.NewRouter()
	f := dbConnect()

	Routers(r, m, f)

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

func TestUpdateGaugeHandlerWithJSON(t *testing.T) {

	h := NewHandler(storage.NewStorage())
	// Create a test request
	req, err := http.NewRequest("POST", "/update/gauge/test/10.0", bytes.NewBufferString(`{"id":"test","type":"gauge","value":10.0}`))
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	h.JSONUpdateHandler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestUpdateCounterHandlerWithJSON(t *testing.T) {

	h := NewHandler(storage.NewStorage())
	// Create a test request
	req, err := http.NewRequest("POST", "/update/", bytes.NewBufferString(`{"id":"test","type":"counter","delta":10}`))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	w := httptest.NewRecorder()

	h.JSONUpdateHandler(w, req)

	defer req.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestUpdateNextCounterHandler(t *testing.T) {

	h := NewHandler(storage.NewStorage())

	err0 := h.MemStorage.SetCounter("test", 10)
	assert.NoError(t, err0)
	// Create a test request
	req, err := http.NewRequest("POST", "/update/", bytes.NewBufferString(`{"id":"test","type":"counter","delta":10}`))
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	h.JSONUpdateHandler(w, req)
	defer req.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	res, ok := h.MemStorage.GetCounterValue("test")
	assert.True(t, ok)
	assert.Equal(t, int64(20), res)
}

func TestUpdateHandlerWithoutJSON(t *testing.T) {

	h := NewHandler(storage.NewStorage())
	// Create a test request
	req, err := http.NewRequest("POST", "/update/gauge/test/10.0", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	h.UpdateHandler(w, req)
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateHandler_InvalidMetricType(t *testing.T) {
	h := NewHandler(storage.NewStorage())

	// Create a test request
	req, err := http.NewRequest("POST", "/update/unknown/test/10.0", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	h.UpdateHandler(w, req)
	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestRun(t *testing.T) {
	// Create a test storage
	m := storage.NewStorage()
	url := "localhost:8080"

	// Create a test router
	r := chi.NewRouter()
	f := dbConnect()

	Routers(r, m, f)

	// Start the server
	go func() {
		err := Run(url, r, m, f)
		assert.NoError(t, err)
	}()
	req, err := http.NewRequest("POST", "/update/gauge/test/10.0", nil)
	assert.NoError(t, err)

	// Create a test response recorder
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetGaugeMetricValueJsonHandler(t *testing.T) {
	// Create a test storage
	h := NewHandler(storage.NewStorage())

	err := h.MemStorage.SetGauge("test", 10.05)

	assert.NoError(t, err)
	// Create a test request
	resp, err := http.NewRequest("POST", "/value/", bytes.NewBuffer([]byte(`{"id":"test","type":"gauge"}`)))

	assert.NoError(t, err)

	w := httptest.NewRecorder()
	// Serve the request
	h.GetJSONMetricValueHandler(w, resp)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":"test","type":"gauge","value":10.05}`, w.Body.String())
}

func TestGetCounterMetricValueJsonHandler(t *testing.T) {
	// Create a test storage
	h := NewHandler(storage.NewStorage())

	err := h.MemStorage.SetCounter("test", 10)

	assert.NoError(t, err)
	// Create a test request
	resp, err := http.NewRequest("POST", "/value/", bytes.NewBuffer([]byte(`{"id":"test","type":"counter"}`)))

	assert.NoError(t, err)

	w := httptest.NewRecorder()
	// // Serve the request
	h.GetJSONMetricValueHandler(w, resp)
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"id":"test","type":"counter","delta":10}`, w.Body.String())
}
