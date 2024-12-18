package main

import (
	"context"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/files"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	database "github.com/ctcsar/metric-and-alerting-system-yp/internal/server/database"
	f "github.com/ctcsar/metric-and-alerting-system-yp/internal/server/flags"
	h "github.com/ctcsar/metric-and-alerting-system-yp/internal/server/handlers"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metrics := storage.NewStorage()
	handler := chi.NewRouter()
	flags := f.NewServerFlags()
	flags.SetServerFlags()
	flag.Parse()
	file := files.NewFile()
	url := url.URL{
		Host: flags.GetServerURL(),
	}
	db, err := database.DBConnect(ctx, flags.GetDatabasePath())
	if err != nil {
		logger.Log.Fatal("cannot connect to database", zap.Error(err))
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
				err = file.WriteFile(metrics, flags.GetStoragePath())
				if err != nil {
					logger.Log.Warn("cannot save to file", zap.Error(err))
					return
				}
				err := database.DBSaveMetrics(ctx, db, metrics)
				if err != nil {
					logger.Log.Error("cannot save metrics to database", zap.Error(err))
					return
				}
				os.Exit(0)
			case <-time.After(time.Duration(flags.GetStoreInterval()) * time.Second):
				err = file.WriteFile(metrics, flags.GetStoragePath())
				if err != nil {
					logger.Log.Warn("cannot save to file", zap.Error(err))
					return
				}
				err := database.DBSaveMetrics(ctx, db, metrics)
				if err != nil {
					logger.Log.Error("cannot save metrics to database", zap.Error(err))
					return
				}
			}
			cancel()
		}
	}()

	if err := h.Run(ctx, url.Host, handler, metrics, db); err != nil {
		logger.Log.Fatal("cannot run handlers", zap.Error(err))
		return
	}
}
