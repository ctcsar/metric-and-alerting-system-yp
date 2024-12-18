package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	f "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/flags"
	handlers "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/handlers"
	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	"go.uber.org/zap"
)

const GaugeMetricsType = "gauge"
const CounterMetricsType = "counter"

func main() {
	flags := f.NewAgentFlags()
	memStorage := storage.MemStorage{}
	flags.SetAgentFlags()
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go memStorage.GetMetrics(flags.GetMetricsGetDuration())
	for {
		select {
		case <-c:
			fmt.Println("Agent stopped")
			os.Exit(0)
		case <-time.After(flags.GetSendDuration() * time.Second):
			metrics := memStorage.Metrics
			if metrics.Gauge != nil || metrics.Counter != nil {
				err := handlers.SendMetric(ctx, flags.GetURLForSend(), &metrics)
				if err != nil {
					logger.Log.Info("cannot send metric:", zap.Error(err))
					return
				}
			}
		}

	}
}
