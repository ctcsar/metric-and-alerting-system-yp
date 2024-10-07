package main

import (
	"fmt"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func main() {

	memStorage := storage.MemStorage{}

	go memStorage.GetMetrics()
	for {
		metrics := memStorage.Metrics
		for k, v := range metrics.Gauge {

			handlers.SendMetric("gauge", k, fmt.Sprintf("%f", v))
		}
		for k, v := range metrics.Counter {
			handlers.SendMetric("counter", k, fmt.Sprintf("%d", v))
		}

		time.Sleep(12 * time.Second)
	}

}
