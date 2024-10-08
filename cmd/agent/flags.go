package main

import (
	"flag"
)

var URLSend string
var GetTime int
var SendTime int

func parseFlags() {
	flag.StringVar(&URLSend, "a", "http://localhost:8080", "address and port to run server")
	flag.IntVar(&GetTime, "p", 2, "time in seconds to get metrics")
	flag.IntVar(&SendTime, "r", 10, "time in seconds to send metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
