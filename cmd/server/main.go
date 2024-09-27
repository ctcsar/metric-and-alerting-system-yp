package main

import (
	"net/http"
	"strconv"
)

type MemStorage struct {
	gauge   float64
	counter []int64
}

func (g *MemStorage) setStorage(v, t string) {
	switch t {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		g.gauge = val
	case "counter":
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		g.counter = append(g.counter, val)
	}
}
func webhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	storrage := &MemStorage{}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.PathValue("name") == "none" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.PathValue("type") != "gauge" || r.PathValue("type") != "counter" || r.PathValue("value") == "none" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storrage.setStorage(r.PathValue("value"), r.PathValue("type"))
	w.WriteHeader(http.StatusOK)
}

func run() error {
	return http.ListenAndServe(`:8080`, nil)
}

func main() {
	http.HandleFunc("/update/{type}/{name}/{value}", webhook)
	//Запускем сервер
	if err := run(); err != nil {
		panic(err)
	}
}
