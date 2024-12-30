package agent

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
)

type agentFlags struct {
	url           string
	sendTime      int
	getMetricTime int
	key           string
}

var err error

func NewAgentFlags() *agentFlags {
	return &agentFlags{}
}
func (f *agentFlags) SetAgentFlags() {
	flag.StringVar(&f.url, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.sendTime, "r", 10, "time in seconds to send metrics")
	flag.IntVar(&f.getMetricTime, "p", 2, "time in seconds to get metrics")
	flag.StringVar(&f.key, "k", "secret", "secret key")
}
func (f agentFlags) GetURLForSend() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.url = envRunAddr
	}
	return f.url

}

func (f agentFlags) GetSendDuration() time.Duration {
	if envGetTime := os.Getenv("REPORT_INTERVAL"); envGetTime != "" {
		f.sendTime, err = strconv.Atoi(envGetTime)
		if err != nil {
			logger.Log.Fatal("cannot convert REPORT_INTERVAL")
		}
	}
	return time.Duration(f.sendTime)
}

func (f agentFlags) GetMetricsGetDuration() time.Duration {
	if envGetTime := os.Getenv("POLL_INTERVAL"); envGetTime != "" {
		f.getMetricTime, err = strconv.Atoi(envGetTime)
		if err != nil {
			logger.Log.Fatal("cannot convert POLL_INTERVAL")
		}
	}
	return time.Duration(f.getMetricTime)
}

func (f agentFlags) GetKey() string {
	if envKey := os.Getenv("KEY"); envKey != "" {
		f.key = envKey
	}
	return f.key
}
