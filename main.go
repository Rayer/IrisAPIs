package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	startTime  = time.Now()
	config     atomic.Pointer[Config]
	checkers   atomic.Pointer[[]ServiceChecker]
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.Store(cfg)
	checkers.Store(&cfg.Services)

	// Watch config for changes
	go watchConfig("config.yaml")

	mux := http.NewServeMux()

	// Health check — always returns 200
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// API health — version + uptime
	mux.HandleFunc("GET /api/health", handleHealth)

	// All services status
	mux.HandleFunc("GET /api/status", handleAllStatus)

	// Single service status
	mux.HandleFunc("GET /api/status/{name}", handleServiceStatus)

	// System metrics (CPU, memory, disk, load)
	mux.HandleFunc("GET /api/system", handleSystem)

	// Middleware: JSON content type + request logging
	handler := withJSON(withLogging(mux))

	log.Printf("IrisAPIs v2 starting on %s", cfg.Host)
	if err := http.ListenAndServe(cfg.Host, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	cfg := config.Load()
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"version": cfg.Version,
		"uptime":  time.Since(startTime).String(),
	})
}

func handleAllStatus(w http.ResponseWriter, r *http.Request) {
	svcs := *checkers.Load()

	// Run checks concurrently
	ch := make(chan ServiceStatus, len(svcs))
	for _, svc := range svcs {
		go func(s ServiceChecker) {
			ch <- s.Check()
		}(svc)
	}

	results := make([]ServiceStatus, 0, len(svcs))
	for range svcs {
		results = append(results, <-ch)
	}

	writeJSON(w, http.StatusOK, results)
}

func handleServiceStatus(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	svcs := *checkers.Load()
	for _, svc := range svcs {
		if strings.EqualFold(svc.GetName(), name) {
			writeJSON(w, http.StatusOK, svc.Check())
			return
		}
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "service not found: " + name})
}

func handleSystem(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, GetSystemInfo())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// --- Middleware ---

func withJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// watchConfig polls for config changes (simple alternative to fsnotify)
func watchConfig(path string) {
	var lastMod time.Time
	for {
		time.Sleep(5 * time.Second)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().After(lastMod) {
			lastMod = info.ModTime()
			cfg, err := LoadConfig(path)
			if err != nil {
				log.Printf("Config reload failed: %v", err)
				continue
			}
			config.Store(cfg)
			checkers.Store(&cfg.Services)
			log.Printf("Config reloaded: %d services", len(cfg.Services))
		}
	}
}
