package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

const GAUGE_METRICS_TYPE = "gauge"
const COUNTER_METRICS_TYPE = "counter"

func main() {
	handlers.GracefulShutdown()
	memStorage := storage.MemStorage{}
	flags.SetAgentFlags()
	flag.Parse()

	go memStorage.GetMetrics(flags.GetMetricsGetDuration())
	for {
		metrics := memStorage.Metrics
		for k, v := range metrics.Gauge {

			err := handlers.SendMetric(flags.GetURLForSend(), GAUGE_METRICS_TYPE, k, fmt.Sprintf("%f", v))
			if err != nil {
				fmt.Println(err)
			}
		}
		for k, v := range metrics.Counter {
			err := handlers.SendMetric(flags.GetURLForSend(), COUNTER_METRICS_TYPE, k, fmt.Sprintf("%d", v))
			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(flags.GetSendDuration() * time.Second)
	}
}
