package handlers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func SendMetric(sendUrl string, metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf(sendUrl+"/update/%s/%s/%s", metricType, metricName, metricValue)
	fmt.Println("Sending metric to " + url)
	_, err := resty.New().R().SetHeader("Content-Type", "text/plain").SetHeader("Content-Length", "0").Post(url)
	if err != nil {
		return err
	}
	return nil
}
