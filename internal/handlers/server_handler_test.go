package handlers

import (
	"net/http"
	"testing"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi/v5"
)

func TestGetMetricValueHandler(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Create a test storage instance
	storageInstance := storage.Storage{
		Gauge: map[string]float64{
			"test": 10.0,
		},
		Counter: 0,
	}

	// Register the GetMetricValueHandler with the test storage instance
	r.Get("/value/{type}/{name}", GetMetricValueHandler(storageInstance))

	// Test cases
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get gauge metric value",
			url:            "/value/gauge/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "10.000000",
		},
		{
			name:           "Get counter metric value",
			url:            "/value/counter/",
			expectedStatus: http.StatusOK,
			expectedBody:   "0.000000",
		},
		{
			name:           "Get unknown metric type",
			url:            "/value/unknown/test",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "Get unknown metric name",
			url:            "/value/gauge/unknown",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// // Create a test request
			// req, err := http.NewRequest("GET", test.url, nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// // Create a test response recorder
			// w := httptest.NewRecorder()

			// // Serve the request
			// r.ServeHTTP(w, req)

			// // Check the response status code
			// assert.Equal(t, test.expectedStatus, w.Code)

			// // Check the response body
			// assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestGetAllMetricsHandler(t *testing.T) {
	// Create a new chi router
	r := chi.NewRouter()

	// Create a test storage instance
	storageInstance := storage.Storage{
		Gauge: map[string]float64{
			"test": 10.0,
		},
		Counter: 0,
	}

	// Register the GetAllMetricsHandler with the test storage instance
	r.Get("/", GetAllMetricsHandler(storageInstance))

	// Test cases
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get all metrics",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html><body><h1>gauge</h1><p>test: 10.000000</p><h1>counter</h1><p>counter: 0.000000</p></body></html>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request
			// req, err := http.NewRequest("GET", test.url, nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// // Create a test response recorder
			// w := httptest.NewRecorder()

			// // Serve the request
			// r.ServeHTTP(w, req)

			// // Check the response status code
			// assert.Equal(t, test.expectedStatus, w.Code)

			// // Check the response body
			// assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestWebhook(t *testing.T) {
	// Register the Webhook with the test storage instance

	// Test cases
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Get gauge metric value",
			url:            "/value/gauge/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "10.000000",
		},
		{
			name:           "Get counter metric value",
			url:            "/value/counter/",
			expectedStatus: http.StatusOK,
			expectedBody:   "0.000000",
		},
		{
			name:           "Get unknown metric type",
			url:            "/value/unknown/test",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "Get unknown metric name",
			url:            "/value/gauge/unknown",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "Get all metrics",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "<html><body><h1>gauge</h1><p>test: 10.000000</p><h1>counter</h1><p>counter: 0.000000</p></body></html>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// // Create a test request
			// req, err := http.NewRequest("GET", test.url, nil)
			// if err != nil {
			// 	t.Fatal(err)
			// }

			// // Create a test response recorder
			// w := httptest.NewRecorder()

			// // Serve the request
			// r.ServeHTTP(w, req)

			// // Check the response status code
			// assert.Equal(t, test.expectedStatus, w.Code)

			// // Check the response body
			// assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
