package flags

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type agentFlags struct {
	URL           string
	SendTime      int
	GetMetricTime int
}

func NewAgentFlags() *agentFlags {
	return &agentFlags{}
}
func (f *agentFlags) SetAgentFlags() {
	flag.StringVar(&f.URL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.SendTime, "r", 10, "time in seconds to send metrics")
	flag.IntVar(&f.GetMetricTime, "p", 2, "time in seconds to get metrics")
}
func (f agentFlags) GetURLForSend() string {
	urlSend := f.URL
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		urlSend = envRunAddr
	}
	return urlSend

}

func (f agentFlags) GetSendDuration() time.Duration {
	if envGetTime := os.Getenv("REPORT_INTERVAL"); envGetTime != "" {
		f.SendTime, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(f.SendTime)
}

func (f agentFlags) GetMetricsGetDuration() time.Duration {
	if envGetTime := os.Getenv("POLL_INTERVAL"); envGetTime != "" {
		f.GetMetricTime, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(f.GetMetricTime)
}
