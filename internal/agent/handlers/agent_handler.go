package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
)

const maxRetries = 3

var retryDelays = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	5 * time.Second,
}

type sendMetrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta"`
	Value float64 `json:"value"`
}

// func SendMetric(sendURL string, metrics *storage.Metrics) error {

// 	var req []sendMetrics

// 	for k, v := range metrics.Gauge {
// 		req = append(req, sendMetrics{
// 			ID:    k,
// 			MType: "gauge",
// 			Value: v,
// 		})
// 	}
// 	for k, v := range metrics.Counter {
// 		req = append(req, sendMetrics{
// 			ID:    k,
// 			MType: "counter",
// 			Delta: v,
// 		})
// 	}
// 	url := url.URL{
// 		Scheme: "http",
// 		Host:   sendURL,
// 		Path:   "/updates/",
// 	}
// 	jsonReq, err := json.Marshal(req)
// 	if err != nil {
// 		return err
// 	}

// 	var buf bytes.Buffer
// 	gz := gzip.NewWriter(&buf)
// 	_, err = gz.Write(jsonReq)
// 	if err != nil {
// 		return err
// 	}
// 	err = gz.Close()
// 	if err != nil {
// 		return err
// 	}

// 	client := resty.New()

// 	resp, err := client.R().
// 		SetBody(buf.Bytes()).
// 		SetHeader("Content-Encoding", "gzip").
// 		Post(url.String())

// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode() != http.StatusOK {
// 		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
// 	}

// 	return nil
// }

func SendMetric(sendURL string, metrics *storage.Metrics) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = sendData(sendURL, metrics)
		if err != nil {
			fmt.Println("Failed to send data:", err)
			time.Sleep(retryDelays[i])
			continue
		}

		break
	}

	if err != nil {
		return fmt.Errorf("failed to send metrics after 3 retries: %w", err)
	}

	return nil
}

func sendData(sendURL string, metrics *storage.Metrics) error {
	var req []sendMetrics

	for k, v := range metrics.Gauge {
		req = append(req, sendMetrics{
			ID:    k,
			MType: "gauge",
			Value: v,
		})
	}
	for k, v := range metrics.Counter {
		req = append(req, sendMetrics{
			ID:    k,
			MType: "counter",
			Delta: v,
		})
	}
	url := url.URL{
		Scheme: "http",
		Host:   sendURL,
		Path:   "/updates/",
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err = gz.Write(jsonReq)
	if err != nil {
		return err
	}
	err = gz.Close()
	if err != nil {
		return err
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(buf.Bytes()).
		SetHeader("Content-Encoding", "gzip").
		Post(url.String())

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned non-OK status code: %d", resp.StatusCode())
	}

	return nil
}
