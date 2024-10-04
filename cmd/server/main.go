package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type Storage struct {
	gauge   map[string]float64
	counter int64
}

func (g *Storage) String() string {
	return fmt.Sprintf("Storage{gauge: %+v, counter: %d}", g.gauge, g.counter)
}

func (g *Storage) setStorage(v, t, n string) {
	switch t {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		g.gauge = map[string]float64{n: val}
	case "counter":
		val, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		g.counter = val
	}
}
func webhook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	m := Storage{}
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

	m.setStorage(r.PathValue("value"), r.PathValue("type"), r.PathValue("name"))
	w.WriteHeader(http.StatusOK)

	for k, v := range m.gauge {
		fmt.Fprintf(w, "%s: %f\n", k, v)
	}
	fmt.Fprintf(w, "counter: %d\n", m.counter)
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
