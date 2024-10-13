package flags

import (
	"flag"
	"os"
	"strconv"
	"time"
)

var urlSend string
var sendTime int
var getMetricsGetDuration int

func SetAgentFlags() {

	flag.StringVar(&urlSend, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&sendTime, "r", 10, "time in seconds to send metrics")
	flag.IntVar(&getMetricsGetDuration, "p", 2, "time in seconds to get metrics")
}
func GetURLForSend() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		urlSend = envRunAddr
	}
	return urlSend

}

func GetSendDuration() time.Duration {
	if envGetTime := os.Getenv("REPORT_INTERVAL"); envGetTime != "" {
		sendTime, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(sendTime)
}

func GetMetricsGetDuration() time.Duration {
	if envGetTime := os.Getenv("POLL_INTERVAL"); envGetTime != "" {
		getMetricsGetDuration, _ = strconv.Atoi(envGetTime)
	}
	return time.Duration(getMetricsGetDuration)
}
