package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/files"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/flags"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	metrics := storage.NewStorage()
	handler := chi.NewRouter()
	flags := flags.NewServerFlags()
	flags.SetServerFlags()
	flag.Parse()
	file := files.NewFile()

	if flags.GetRestore() {
		file.ReadFromFile(flags.StoragePath, metrics)
	}
	go func() {
		for {
			select {
			case <-c:
				file.WriteFile(metrics, flags.StoragePath)
				os.Exit(0)
			case <-time.After(time.Duration(flags.GetStoreInterval()) * time.Second):
				file.WriteFile(metrics, flags.StoragePath)
			}
		}
	}()

	if err := handlers.Run(flags.GetServerURL(), handler, metrics); err != nil {
		fmt.Println(err)
		return
	}
}
