package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"

	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/retry"
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

func SendMetric(ctx context.Context, sendURL string, metrics *storage.Metrics, secretKey string) error {

	err := retry.Retry(func() error {
		return sendData(sendURL, metrics, secretKey)
	}, maxRetries, ctx, retryDelays)

	if err != nil {
		return fmt.Errorf("failed to send metrics after %d retries: %w", maxRetries, err)
	}

	return nil
}

func sendData(sendURL string, metrics *storage.Metrics, secretKey string) error {
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

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(jsonReq)
	dst := h.Sum(nil)

	resp, err := client.R().
		SetBody(buf.Bytes()).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("HashSHA256", fmt.Sprintf("%x", dst)).
		Post(url.String())

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server returned non-OK status code: %d", resp.StatusCode())
	}

	return nil
}
