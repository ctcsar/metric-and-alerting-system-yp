package main

import (
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
)

func main() {
	g := storage.NewGaugeStorage()
	c := storage.NewCounterStorage()
	r := chi.NewRouter()
	//Запускем сервер
	if err := handlers.Run(r, g, c); err != nil {
		panic(err)
	}
}
