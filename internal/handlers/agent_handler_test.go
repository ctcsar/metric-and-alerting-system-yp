package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	// Create a test server to mock the HTTP request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request URL and headers
		assert.Equal(t, "/update/gauge/test/10.0", r.URL.String())
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
		assert.Equal(t, "0", r.Header.Get("Content-Length"))

		// Write a response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer ts.Close()
}
