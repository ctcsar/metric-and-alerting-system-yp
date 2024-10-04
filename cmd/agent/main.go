package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand/v2"
	"net/http"
	"runtime"
	"time"
)

type MemStorage struct {
	Metrics Metrics
}

type Metrics struct {
	gauge   map[string]float64
	counter int64
}

func (m *MemStorage) String() string {
	// implement the string representation of MemStorage
	return fmt.Sprintf("MemStorage{Metrics: %+v}", m.Metrics)
}
func (m *MemStorage) setStorage(r float64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Metrics.gauge = map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   r,
	}

}

func (m *MemStorage) SetCounter(count int64) {

	m.Metrics.counter = count
}

func sendMetric(metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func (m *MemStorage) GetMetrics() MemStorage {
	var counter int64 = 0
	for {
		RandomValue := rand.Float64()
		m.setStorage(RandomValue)
		m.SetCounter(counter)
		time.Sleep(2 * time.Second)
		counter++
	}
}
func main() {

	memStorage := MemStorage{}

	go memStorage.GetMetrics()
	for {
		metrics := memStorage.Metrics
		for k, v := range metrics.gauge {

			// fmt.Printf("%s: %f\n", k, v)
			sendMetric("gauge", k, fmt.Sprintf("%f", v))
		}
		sendMetric("counter", "counter", fmt.Sprintf("%d", metrics.counter))

		time.Sleep(12 * time.Second)
	}

}
