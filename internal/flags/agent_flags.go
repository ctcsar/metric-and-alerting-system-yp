package flags

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type agentFlags struct {
	Url           string
	SendTime      int
	GetMetricTime int
}

func NewFlags() *agentFlags {
	return &agentFlags{}
}
func (f *agentFlags) SetAgentFlags() {
	flag.StringVar(&f.Url, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.SendTime, "r", 10, "time in seconds to send metrics")
	flag.IntVar(&f.GetMetricTime, "p", 2, "time in seconds to get metrics")
}
func (f agentFlags) GetURLForSend() string {
	urlSend := f.Url
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		urlSend = envRunAddr
	}
	return urlSend

}

func (f agentFlags) GetSendDuration() time.Duration {
	sendTime := f.SendTime
	if envGetTime := os.Getenv("REPORT_INTERVAL"); envGetTime != "" {
		sendTime, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(sendTime)
}

func (f agentFlags) GetMetricsGetDuration() time.Duration {
	getMetricsGetDuration := f.GetMetricTime
	if envGetTime := os.Getenv("POLL_INTERVAL"); envGetTime != "" {
		getMetricsGetDuration, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(getMetricsGetDuration)
}
