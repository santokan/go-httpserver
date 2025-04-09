package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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
