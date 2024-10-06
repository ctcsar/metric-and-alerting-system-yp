package main

import (
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	m := storage.Storage{}
	handlers.Webhook(r, m)
	//Запускем сервер
	if err := handlers.Run(r, m); err != nil {
		panic(err)
	}
}
