package agent

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

type MemStorage struct {
	Metrics Metrics
	Mutex   sync.RWMutex
}

type Metrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *MemStorage) String() string {
	return fmt.Sprintf("MemStorage{Metrics: %+v}", m.Metrics)
}
func (m *MemStorage) SetStorage(rand float64) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Metrics.Gauge = map[string]float64{
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
		"RandomValue":   rand,
	}
}

func (m *MemStorage) SetCounter(count int64) {

	m.Metrics.Counter = map[string]int64{
		"PollCount": count,
	}
}

func (m *MemStorage) SetMoreMetrics(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c:
			fmt.Println("Received Interrupt, stopping metric collection.")
			return
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping metric collection.")
			return
		case <-ticker.C:
			m.updateMetrics()
		}
	}
}

func (m *MemStorage) updateMetrics() {
	virtualMem, _ := mem.VirtualMemory()
	CPUInfo, err := cpu.Percent(0, false)
	if err != nil {
		logger.Log.Fatal("cannot get CPU percent", zap.Error(err))
	}

	m.Metrics.Gauge = map[string]float64{
		"TotalMemory":     float64(virtualMem.Total),
		"FreeMemory":      float64(virtualMem.Free),
		"CPUutilization1": CPUInfo[0],
	}
}
func (m *MemStorage) GetMetrics(duratiomTime time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var counter int64 = 0
	for {
		select {
		case <-c:
			return
		case <-time.After(duratiomTime * time.Second):
			RandomValue := rand.Float64()
			m.SetStorage(RandomValue)
			m.SetCounter(counter)
			time.Sleep(duratiomTime * time.Second)
			counter++
		}
	}
}
