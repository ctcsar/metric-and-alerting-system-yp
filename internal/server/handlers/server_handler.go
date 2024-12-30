package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/compress"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Handler struct {
	MemStorage *storage.Storage
}

func NewHandler(metrics *storage.Storage) *Handler {
	return &Handler{MemStorage: metrics}
}

func (h *Handler) GetMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetJSONMetricValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var buff Metrics
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&buff)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch buff.MType {
	case "gauge":
		val, ok := h.MemStorage.GetGaugeValue(buff.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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
		_, err = w.Write(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		val, ok := h.MemStorage.GetCounterValue(buff.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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
		_, err = w.Write(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
func (h *Handler) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.MemStorage.SetGauge(name, val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.MemStorage.SetCounter(name, val)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) JSONUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var buff Metrics
	var resp Metrics

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&buff)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if buff.MType != "gauge" && buff.MType != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch buff.MType {

	case "gauge":
		err = h.MemStorage.SetGauge(buff.ID, *buff.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		val, ok := h.MemStorage.GetGaugeValue(buff.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp = Metrics{
			ID:    buff.ID,
			MType: buff.MType,
			Value: &val,
		}
		r, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "counter":
		err = h.MemStorage.SetCounter(buff.ID, *buff.Delta)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		val, ok := h.MemStorage.GetCounterValue(buff.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resp = Metrics{
			ID:    buff.ID,
			MType: buff.MType,
			Delta: &val,
		}
		r, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) JSONUpdateAllMetricsHandler(secretKey string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var buff []Metrics

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&buff)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respHash := r.Header.Get("HashSHA256")
	hashData, err := hex.DecodeString(respHash)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := hmac.New(sha256.New, []byte(secretKey))
	jsonBuff, err := json.Marshal(buff)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash.Write(jsonBuff)
	dst := hash.Sum(nil)
	if !hmac.Equal(hashData, dst) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, metric := range buff {
		switch metric.MType {
		case "gauge":
			err = h.MemStorage.SetGauge(metric.ID, *metric.Value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			err = h.MemStorage.SetCounter(metric.ID, *metric.Delta)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) PingHandler(ctx context.Context, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if err := db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Routers(ctx context.Context, handler chi.Router, metrics *storage.Storage, db *sql.DB, secretKey string) {
	h := Handler{
		MemStorage: metrics,
	}
	handler.Get("/value/{type}/{name}", h.GetMetricValueHandler)
	handler.Get("/", h.GetAllMetricsHandler)
	handler.Post("/update/{type}/{name}/{value}", h.UpdateHandler)
	handler.Post("/update/", h.JSONUpdateHandler)
	handler.Post("/updates/", func(w http.ResponseWriter, r *http.Request) { h.JSONUpdateAllMetricsHandler(secretKey, w, r) })
	handler.Post("/value/", h.GetJSONMetricValueHandler)
	handler.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		h.PingHandler(ctx, db, w, r)
	})
}

func Run(ctx context.Context, url string, handler chi.Router, metrics *storage.Storage, db *sql.DB, secretKey string) error {
	logger.Log.Info("starting server", zap.String("url", url))
	handler = logger.RequestLogger(handler)
	Routers(ctx, handler, metrics, db, secretKey)
	return http.ListenAndServe(url, compress.GzipMiddleware(handler))
}
