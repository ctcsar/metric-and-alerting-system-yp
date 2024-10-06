package handlers

import (
	"fmt"
	"net/http"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
)

func GetMetricValueHandler(m storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		value, ok := m.GetMetricValue(metricType, metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%f", value)
	}
}

func GetAllMetricsHandler(m storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := m.GetAllMetrics()
		html := "<html><body>"
		for metricType, metricValues := range metrics {
			html += fmt.Sprintf("<h1>%s</h1>", metricType)
			for metricName, value := range metricValues {
				html += fmt.Sprintf("<p>%s: %f</p>", metricName, value)
			}
		}
		html += "</body></html>"
		w.Write([]byte(html))
	}
}

func UpdateHandler(r chi.Router, m storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		name := chi.URLParam(r, "name")
		metricType := chi.URLParam(r, "type")
		value := chi.URLParam(r, "value")

		if name == "none" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if metricType != "gauge" && metricType != "counter" || value == "none" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m.SetStorage(value, metricType, name)
		w.WriteHeader(http.StatusOK)

		res := m.GetAllMetrics()
		counterValue, ok := m.GetMetricValue("counter", "counter")
		if !ok {
			fmt.Fprintf(w, "counter: %+v\n", counterValue)
		}
		fmt.Fprintf(w, "gauge: %+v", res["gauge"])
	}
}

func Webhook(r chi.Router, m storage.Storage) {
	r.Get("/value/{type}/{name}", GetMetricValueHandler(m))
	r.Get("/", GetAllMetricsHandler(m))
	r.Post("/update/{type}/{name}/{value}", UpdateHandler(r, m))
}

func Run(r chi.Router, m storage.Storage) error {
	Webhook(r, m)
	return http.ListenAndServe(":8080", r)
}
