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

			// fmt.Printf("%s: %f\n", k, v)
			handlers.SendMetric("gauge", k, fmt.Sprintf("%f", v))
		}
		handlers.SendMetric("counter", "counter", fmt.Sprintf("%d", metrics.Counter))

		time.Sleep(12 * time.Second)
	}

}
