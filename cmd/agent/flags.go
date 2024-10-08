package main

import (
	"flag"
	"os"
	"strconv"
)

var URLSend string
var GetTime int
var SendTime int

func parseFlags() {
	flag.StringVar(&URLSend, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&SendTime, "r", 10, "time in seconds to send metrics")
	flag.IntVar(&GetTime, "p", 2, "time in seconds to get metrics")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		URLSend = envRunAddr
	}
	if envSetTime := os.Getenv("REPORT_INTERVAL"); envSetTime != "" {
		GetTime, _ = strconv.Atoi(envSetTime)
	}
	if envGetTime := os.Getenv("POLL_INTERVAL"); envGetTime != "" {
		SendTime, _ = strconv.Atoi(envGetTime)
	}

}
