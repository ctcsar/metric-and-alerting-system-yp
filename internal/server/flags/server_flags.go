package server

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
)

type serverFlags struct {
	url           string
	storeInterval int
	storagePath   string
	restore       bool
	databaseDSN   string
}

func NewServerFlags() *serverFlags {
	return &serverFlags{}
}

var err error

func (f *serverFlags) SetServerFlags() {
	flag.StringVar(&f.url, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.storeInterval, "i", 300, "duration to save metrics in file")
	flag.StringVar(&f.storagePath, "f", "storage.txt", "name of file to save metrics")
	flag.BoolVar(&f.restore, "r", true, "restore metrics from file")
	flag.StringVar(&f.databaseDSN, "d", "metrics:password@localhost:5432/metrics?sslmode=disable", "path to database")
}

func (f *serverFlags) GetServerURL() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.url = envRunAddr
	}

	return f.url
}

func (f *serverFlags) GetStoragePath() string {
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		f.storagePath = envStoragePath
	}

	return f.storagePath
}

func (f *serverFlags) GetStoreInterval() time.Duration {
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		f.storeInterval, err = strconv.Atoi(envStoreInterval)
		if err != nil {
			logger.Log.Fatal("cannot convert STORE_INTERVAL")
		}
	}

	return time.Duration(f.storeInterval)
}

func (f *serverFlags) GetRestore() bool {
	return f.restore
}

func (f *serverFlags) GetDatabasePath() string {
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		f.databaseDSN = envDatabaseDSN
	}
	return f.databaseDSN
}
