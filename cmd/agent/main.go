package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

const GaugeMetricsType = "gauge"
const CounterMetricsType = "counter"

func main() {
	handlers.GracefulShutdown()
	memStorage := storage.MemStorage{}
	flags.SetAgentFlags()
	flag.Parse()

	go memStorage.GetMetrics(flags.GetMetricsGetDuration())
	for {
		metrics := memStorage.Metrics
		for k, v := range metrics.Gauge {

			err := handlers.SendMetric(flags.GetURLForSend(), GaugeMetricsType, k, fmt.Sprintf("%f", v))
			if err != nil {
				fmt.Println(err)
			}
		}
		for k, v := range metrics.Counter {
			err := handlers.SendMetric(flags.GetURLForSend(), CounterMetricsType, k, fmt.Sprintf("%d", v))
			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(flags.GetSendDuration() * time.Second)
	}
}
