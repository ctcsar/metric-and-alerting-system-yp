package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Handler struct {
	http.Handler
	MemStorage *storage.Storage
}

func GetMetricValueHandler(metrics *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buff Metrics
		rj := false
		w.Header().Set("Content-Type", "application/json;")
		h := Handler{MemStorage: metrics}

		if r.Body == http.NoBody {
			buff.ID = chi.URLParam(r, "name")
			buff.MType = chi.URLParam(r, "type")
		} else {
			rj = true
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&buff)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		switch buff.MType {
		case "gauge":
			val, ok := h.MemStorage.GetGaugeValue(buff.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if rj {
				resp := Metrics{
					ID:    buff.ID,
					MType: buff.MType,
					Value: &val,
				}
				r, err := json.Marshal(resp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(r)
			} else {
				_, err := w.Write([]byte(fmt.Sprintf("%g", val)))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return

				}
			}
		case "counter":
			val, ok := h.MemStorage.GetCounterValue(buff.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if rj {
				resp := Metrics{
					ID:    buff.ID,
					MType: buff.MType,
					Delta: &val,
				}
				r, err := json.Marshal(resp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(r)
			} else {
				_, err := w.Write([]byte(fmt.Sprintf("%d", val)))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
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
	w.Header().Set("Content-Type", "application/json;")

	var buff Metrics

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Body == http.NoBody {
		buff.ID = chi.URLParam(r, "name")
		buff.MType = chi.URLParam(r, "type")
		if buff.MType == "counter" {
			val, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			buff.Delta = &val
		} else if buff.MType == "gauge" {
			val, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			buff.Value = &val
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&buff)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if buff.ID == "none" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if buff.MType == "gauge" {
		err := h.MemStorage.SetGauge(buff.ID, *buff.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if buff.MType == "counter" {
		err := h.MemStorage.SetCounter(buff.ID, *buff.Delta)
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
	handler.Post("/update", h.UpdateHandler)
	handler.Post("/value", GetMetricValueHandler(metrics))
}

func Run(url string, handler chi.Router, metrics *storage.Storage) error {
	logger.Log.Info("starting server", zap.String("url", url))
	handler = logger.RequestLogger(handler)
	Routers(handler, metrics)
	return http.ListenAndServe(url, handler)
}
