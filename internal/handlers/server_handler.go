package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
	"github.com/go-chi/chi"
)

type Handler struct {
	http.Handler
	GaugeStorage   *storage.Storage
	CounterStorage *storage.Storage
}

func GetMetricValueHandler(g *storage.Storage, c *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		h := Handler{GaugeStorage: g, CounterStorage: c}
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		switch metricType {
		case "gauge":
			val, ok := h.GaugeStorage.GetGaugeValue(metricName)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%g", val)
			return

		case "counter":
			val, ok := h.CounterStorage.GetCounterValue(metricName)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
			}
			fmt.Fprintf(w, "%d", val)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func GetAllMetricsHandler(g *storage.Storage, c *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := Handler{GaugeStorage: g, CounterStorage: c}
		gauge := h.GaugeStorage.GetAllGaugeMetrics()
		counter := h.CounterStorage.GetAllCounterMetrics()
		html := "<html><body>"
		for metricType, metricValues := range gauge {
			html += fmt.Sprintf("<h1>%s</h1>", metricType)
			for metricName, value := range metricValues {
				html += fmt.Sprintf("<p>%s: %g</p>", metricName, value)
			}
		}
		for metricType, metricValues := range counter {
			html += fmt.Sprintf("<h1>%s</h1>", metricType)
			for metricName, value := range metricValues {
				html += fmt.Sprintf("<p>%s: %d</p>", metricName, value)
			}
		}
		html += "</body></html>"
		w.Write([]byte(html))
	}
}

func UpdateHandler(r chi.Router, g *storage.Storage, c *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		h := Handler{GaugeStorage: g, CounterStorage: c}

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

		if metricType == "gauge" {
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			h.GaugeStorage.SetGauge(name, val)
		} else if metricType == "counter" {
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			h.CounterStorage.SetCounter(name, val)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func Webhook(r chi.Router, g *storage.Storage, c *storage.Storage) {
	r.Get("/value/{type}/{name}", GetMetricValueHandler(g, c))
	r.Get("/", GetAllMetricsHandler(g, c))
	r.Post("/update/{type}/{name}/{value}", UpdateHandler(r, g, c))
}

func Run(url string, r chi.Router, g *storage.Storage, c *storage.Storage) error {
	Webhook(r, g, c)
	return http.ListenAndServe(url, r)
}
