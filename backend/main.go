package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	count int
	mu    sync.Mutex
)

var countGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "app_count",
	Help: "Current value of the application count.",
})

var countChanges = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "app_count_changes_total",
	Help: "Total number of count changes.",
})

func sendCount(w http.ResponseWriter, value int) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"count": value,
	})
}

func changeCount(delta int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count += delta
		current := count

		countGauge.Set(float64(current))
		countChanges.Inc()
		mu.Unlock()

		sendCount(w, current)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Match this to the address of your Vite frontend.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	prometheus.MustRegister(countGauge, countChanges)

	mux := http.NewServeMux()
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /api/count", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		current := count
		mu.Unlock()

		sendCount(w, current)
	})

	mux.HandleFunc("POST /api/increment", changeCount(1))
	mux.HandleFunc("POST /api/decrement", changeCount(-1))

	log.Println("Server running at http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", cors(mux)))
}
