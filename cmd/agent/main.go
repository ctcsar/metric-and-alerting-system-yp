package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	f "github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

const GaugeMetricsType = "gauge"
const CounterMetricsType = "counter"

func main() {
	flags := f.NewAgentFlags()
	memStorage := storage.MemStorage{}
	flags.SetAgentFlags()
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go memStorage.GetMetrics(flags.GetMetricsGetDuration())
	for {
		select {
		case <-c:
			return
		case <-time.After(flags.GetSendDuration() * time.Second):
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
		}

	}
}
