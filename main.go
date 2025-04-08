package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlTemplate, err := os.ReadFile("admin/metrics.html")
	if err != nil {
		log.Printf("Error reading admin template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hits := cfg.fileserverHits.Load()
	htmlContent := fmt.Sprintf(string(htmlTemplate), hits)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(htmlContent))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Metrics reset"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func handlerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func main() {
	const (
		filePathRoot = "."
		port         = "8080"
	)

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handerMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on http://localhost:%s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
