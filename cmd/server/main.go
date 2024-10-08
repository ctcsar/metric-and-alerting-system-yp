package main

import (
	"fmt"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
)

func main() {
	parseFlags()
	fmt.Println("Running server on", flagRunAddr)

	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()
	r := chi.NewRouter()
	//Запускем сервер
	if err := handlers.Run(flagRunAddr, r, g, c); err != nil {
		panic(err)
	}
}
