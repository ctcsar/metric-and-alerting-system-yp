package flags

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type serverFlags struct {
	URL           string
	StoreInterval int
	StoragePath   string
	Restore       bool
}

func NewServerFlags() *serverFlags {
	return &serverFlags{}
}
func (f *serverFlags) SetServerFlags() {
	flag.StringVar(&f.URL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&f.StoreInterval, "i", 300, "duration to save metrics in file")
	flag.StringVar(&f.StoragePath, "f", "storage.txt", "name of file to save metrics")
	flag.BoolVar(&f.Restore, "r", true, "restore metrics from file")
}

func (f *serverFlags) GetServerURL() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		f.URL = envRunAddr
	}

	return f.URL
}

func (f *serverFlags) GetStoragePath() string {
	if envStoragePath := os.Getenv("FILE_STORAGE_PATH"); envStoragePath != "" {
		f.StoragePath = envStoragePath
	}

	return f.StoragePath
}

func (f *serverFlags) GetStoreInterval() time.Duration {
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		f.StoreInterval, _ = strconv.Atoi(envStoreInterval)
	}

	return time.Duration(f.StoreInterval)
}

func (f *serverFlags) GetRestore() bool {
	return f.Restore
}
