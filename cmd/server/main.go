package main

import (
	"net/http"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
)

func main() {
	http.HandleFunc("/update/{type}/{name}/{value}", handlers.Webhook)
	//Запускем сервер
	if err := handlers.Run(); err != nil {
		panic(err)
	}
}
