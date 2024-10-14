package main

import (
	"flag"
	"fmt"

	"github.com/go-chi/chi"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func main() {
	metrics := storage.NewStorage()
	handler := chi.NewRouter()
	flags.SetServerFlags()
	flag.Parse()

	if err := handlers.Run(flags.GetServerURL(), handler, metrics); err != nil {
		fmt.Println(err)
		return
	}

}
