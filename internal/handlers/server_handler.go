package handlers

import (
	"fmt"
	"net/http"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m := storage.Storage{}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.PathValue("name") == "none" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.PathValue("type") != "gauge" && r.PathValue("type") != "counter" || r.PathValue("value") == "none" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m.SetStorage(r.PathValue("value"), r.PathValue("type"), r.PathValue("name"))
	w.WriteHeader(http.StatusOK)

	for k, v := range m.Gauge {
		fmt.Fprintf(w, "%s: %f\n", k, v)
	}
	fmt.Fprintf(w, "counter: %d\n", m.Counter)
}

func Run() error {
	return http.ListenAndServe(`:8080`, nil)
}
