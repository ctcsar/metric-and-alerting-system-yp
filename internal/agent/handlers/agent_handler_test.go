package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
)

func TestSendMetric(t *testing.T) {
	ctx := context.Background()
	// Create a test server to mock the HTTP request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request URL and headers
		assert.Equal(t, "/update/gauge/test/10.0", r.URL.String())
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		assert.Equal(t, "0", r.Header.Get("Content-Length"))

		// Write a response
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		assert.NoError(t, err)

		t.Run("Valid URL and metrics", func(t *testing.T) {
			metrics := &storage.Metrics{
				Gauge: map[string]float64{
					"test": 10.0,
				},
			}

			err := SendMetric(ctx, r.URL.String(), metrics)
			assert.NoError(t, err)
		})

		t.Run("Invalid URL", func(t *testing.T) {
			metrics := &storage.Metrics{
				Gauge: map[string]float64{
					"test": 10.0,
				},
			}

			err := SendMetric(ctx, "invalidURL", metrics)
			assert.Error(t, err)
		})

	}))

	defer ts.Close()

}

func TestSendData(t *testing.T) {
	// Create a test server to mock the HTTP request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request URL and headers
		assert.Equal(t, "/update/gauge/test/10.0", r.URL.String())
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		assert.Equal(t, "0", r.Header.Get("Content-Length"))

		// Write a response
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		assert.NoError(t, err)

		t.Run("Valid URL and metrics", func(t *testing.T) {
			metrics := &storage.Metrics{
				Gauge: map[string]float64{
					"test": 10,
				},
			}

			err := sendData(r.URL.String(), metrics)
			assert.NoError(t, err)
		})

		t.Run("Invalid URL", func(t *testing.T) {
			metrics := &storage.Metrics{
				Gauge: map[string]float64{
					"test": 10,
				},
			}

			err := sendData("invalidURL", metrics)
			assert.Error(t, err)
		})
		t.Run("Different types of metrics", func(t *testing.T) {
			metrics := &storage.Metrics{
				Gauge: map[string]float64{
					"test": 10.0,
				},
				Counter: map[string]int64{
					"count": 10,
				},
			}

			err := sendData(r.URL.String(), metrics)
			assert.NoError(t, err)
		})

	}))

	defer ts.Close()

}
