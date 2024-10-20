package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

type Handler struct {
	http.Handler
	MemStorage *storage.Storage
}

func GetMetricValueHandler(metrics *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		h := Handler{MemStorage: metrics}
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		switch metricType {
		case "gauge":
			val, ok := h.MemStorage.GetGaugeValue(metricName)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err := w.Write([]byte(fmt.Sprintf("%g", val)))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "counter":
			val, ok := h.MemStorage.GetCounterValue(metricName)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err := w.Write([]byte(fmt.Sprintf("%d", val)))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func GetAllMetricsHandler(metrics *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := Handler{MemStorage: metrics}
		gauge := h.MemStorage.GetAllGaugeMetrics()
		counter := h.MemStorage.GetAllCounterMetrics()
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
		_, err := w.Write([]byte(html))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	if metricType == "gauge" {
		err := h.MemStorage.SetGauge(name, value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if metricType == "counter" {
		err := h.MemStorage.SetCounter(name, value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func Routers(handler chi.Router, metrics *storage.Storage) {
	h := Handler{MemStorage: metrics}
	handler.Get("/value/{type}/{name}", GetMetricValueHandler(metrics))
	handler.Get("/", GetAllMetricsHandler(metrics))
	handler.Post("/update/{type}/{name}/{value}", h.UpdateHandler)
}

func Run(url string, handler chi.Router, metrics *storage.Storage) error {
	logger.Log.Info("starting server", zap.String("url", url))
	handler = logger.RequestLogger(handler)
	Routers(handler, metrics)
	return http.ListenAndServe(url, handler)
}
