package main

import (
	"fmt"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func main() {

	parseFlags()
	fmt.Println("Running agent on", UrlSend)
	fmt.Printf("get metrics every %d seconds", GetTime)
	fmt.Printf("send metrics every %d seconds", SendTime)
	memStorage := storage.MemStorage{}

	go memStorage.GetMetrics(GetTime)
	for {
		metrics := memStorage.Metrics
		for k, v := range metrics.Gauge {

			handlers.SendMetric(UrlSend, "gauge", k, fmt.Sprintf("%f", v))
		}
		for k, v := range metrics.Counter {
			handlers.SendMetric(UrlSend, "counter", k, fmt.Sprintf("%d", v))
		}

		time.Sleep(time.Duration(SendTime) * time.Second)
	}
}
