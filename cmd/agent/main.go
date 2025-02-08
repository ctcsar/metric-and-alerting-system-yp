package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	f "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/flags"
	handlers "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/handlers"
	storage "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/storage"
	workerpool "github.com/ctcsar/metric-and-alerting-system-yp/internal/agent/workers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
)

const GaugeMetricsType = "gauge"
const CounterMetricsType = "counter"

func main() {
	flags := f.NewAgentFlags()
	memStorage := storage.MemStorage{}
	flags.SetAgentFlags()
	flag.Parse()
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sendTicker := time.NewTicker(flags.GetSendDuration() * time.Second)

	workerPool := workerpool.NewWorkerPool(flags.GetRateLimit())

	workerPool.Start(ctx)

	defer workerPool.Stop()

	go memStorage.GetMetrics(flags.GetMetricsGetDuration())
	go memStorage.SetMoreMetrics(ctx)
	for {
		select {
		case <-c:
			fmt.Println("Agent stopped")
			return
		case <-sendTicker.C:
			workerPool.SubmitTask(func() {
				memStorage.Mutex.RLock()
				defer memStorage.Mutex.RUnlock()
				metrics := memStorage.Metrics
				if metrics.Gauge != nil || metrics.Counter != nil {
					err := handlers.SendMetric(ctx, flags.GetURLForSend(), &metrics, flags.GetKey())
					if err != nil {
						logger.Log.Error("cannot send metric:", zap.Error(err))
					}
				}
			})
		}
	}
}
