package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type sendMetrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func SendMetric(sendURL string, metricType string, metricName string, metricValue string) error {
	var req sendMetrics

	req.ID = metricName
	req.MType = metricType
	switch metricType {
	case "counter":
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		req.Delta = &val
	case "gauge":
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		req.Value = &val
	}
	url := url.URL{
		Scheme: "http",
		Host:   sendURL,
		Path:   "/update/",
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Compress the request body using gzip
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

	// Set the request body and header
	resp, err := client.R().
		SetBody(buf.Bytes()).
		SetHeader("Content-Encoding", "gzip").
		Post(url.String())

	// Check the response status code
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
