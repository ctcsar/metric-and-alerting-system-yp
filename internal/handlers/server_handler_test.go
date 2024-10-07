// agent_handler_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMetric(t *testing.T) {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		metricValue    string
		expectedStatus int
	}{
		{
			name:           "Send gauge metric",
			metricType:     "gauge",
			metricName:     "test",
			metricValue:    "10.0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Send counter metric",
			metricType:     "counter",
			metricName:     "test",
			metricValue:    "10",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Send invalid metric type",
			metricType:     "unknown",
			metricName:     "test",
			metricValue:    "10.0",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test server to mock the HTTP request
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.expectedStatus)
			}))
			defer ts.Close()

			// Update the URL in the SendMetric function to point to the test server
			// url := fmt.Sprintf("%s/update/%s/%s/%s", ts.URL, test.metricType, test.metricName, test.metricValue)

			// Call the SendMetric function
			err := SendMetric(test.metricType, test.metricName, test.metricValue)

			// Check if the error is nil
			if err != nil {
				t.Errorf("expected nil error, got %v", err)
			}
		})
	}
}
