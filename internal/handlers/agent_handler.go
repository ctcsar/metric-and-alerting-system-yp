package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendMetric(metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
