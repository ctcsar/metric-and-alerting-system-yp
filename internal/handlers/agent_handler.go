package handlers

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-resty/resty/v2"
)

func SendMetric(sendURL string, metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf("http://"+sendURL+"/update/%s/%s/%s", metricType, metricName, metricValue)
	_, err := resty.New().R().SetHeader("Content-Type", "text/plain").SetHeader("Content-Length", "0").Post(url)
	if err != nil {
		return err
	}
	return nil
}

func GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Shutting down gracefully")
		os.Exit(0)
	}()
}
