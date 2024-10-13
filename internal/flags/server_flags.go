package flags

import (
	"flag"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func SetServerFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
}

func GetServerURL() string {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	return flagRunAddr
}
