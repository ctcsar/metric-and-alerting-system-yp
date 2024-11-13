package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/files"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	f "github.com/ctcsar/metric-and-alerting-system-yp/internal/server/flags"
	h "github.com/ctcsar/metric-and-alerting-system-yp/internal/server/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	metrics := storage.NewStorage()
	handler := chi.NewRouter()
	flags := f.NewServerFlags()
	flags.SetServerFlags()
	flag.Parse()
	file := files.NewFile()
	url := url.URL{
		Host: flags.GetServerURL(),
	}

	ps := fmt.Sprintf("%s sslmode=disable",
		flags.GetDatabasePath())

	db, err := pgxpool.New(context.Background(), ps)
	if err != nil {
		logger.Log.Info("cannot connect to database", zap.Error(err))
	}
	defer db.Close()

	if flags.GetRestore() {
		err := file.ReadFromFile(flags.GetStoragePath(), metrics)
		if err != nil {
			logger.Log.Info("cannot read file", zap.Error(err))
		}
	}
	go func() {
		for {
			select {
			case <-c:
				err := file.WriteFile(metrics, flags.GetStoragePath())
				if err != nil {
					fmt.Println(err)
				}
				os.Exit(0)
			case <-time.After(time.Duration(flags.GetStoreInterval()) * time.Second):
				err := file.WriteFile(metrics, flags.GetStoragePath())
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}()

	if err := h.Run(url.Host, handler, metrics, db); err != nil {
		fmt.Println(err)
		return
	}
}
