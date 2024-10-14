package flags

import (
	"flag"
	"os"
)

var flagRunAddr string

func SetServerFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
}

func GetServerURL() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	return flagRunAddr
}
