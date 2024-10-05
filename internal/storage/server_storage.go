package storage

import (
	"fmt"
	"strconv"
)

type Storage struct {
	Gauge   map[string]float64
	Counter int64
}

func (g *Storage) String() string {
	return fmt.Sprintf("Storage{gauge: %+v, counter: %d}", g.Gauge, g.Counter)
}

func (g *Storage) SetStorage(v, t, n string) {
	switch t {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		g.Gauge = map[string]float64{n: val}
	case "counter":
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		g.Counter = val
	}
}
