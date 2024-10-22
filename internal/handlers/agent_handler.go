package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const updateURLFormat = "http://%s/update"

type sendMetrics struct {
	Id    string   `json:"id"`
	MType string   `json:"type"`
	Value *float64 `json:"value"`
	Delta *int64   `json:"delta"`
}

func SendMetric(sendURL string, metricType string, metricName string, metricValue string) error {
	var req sendMetrics

	req.Id = metricName
	req.MType = metricType
	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return err
		}
		req.Delta = &val
	} else if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return err
		}
		req.Value = &val
	}
	url := fmt.Sprintf(updateURLFormat, sendURL)

	jsonReq, _ := json.Marshal(req)

	_, err := resty.New().R().SetHeader("Content-Type", "application/json").SetBody(jsonReq).Post(url)
	if err != nil {
		return err
	}
	return nil
}
