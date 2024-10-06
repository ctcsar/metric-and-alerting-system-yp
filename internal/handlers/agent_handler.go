package handlers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func SendMetric(metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	req, err := resty.New().R().SetHeader("Content-Type", "text/plain").SetHeader("Content-Length", "0").Post(url)
	if err != nil {
		return err
	}

	fmt.Println("request: ", req)
	return nil
}
