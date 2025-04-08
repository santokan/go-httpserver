package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
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

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
		return
	}

	if len(requestBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	cleanedBody := cleanChirp(requestBody.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseBody := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: cleanedBody,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(responseBody)
	if err != nil {
		log.Printf("Error encoding response JSON: %v", err)
		// Attempt to send a generic error if encoding fails
		http.Error(w, `{"error": "Something went wrong"}`, http.StatusInternalServerError)
	}
}

// cleanChirp replaces profane words in a string with "****".
func cleanChirp(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	cleanedWords := make([]string, len(words))

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		isProfane := slices.Contains(profaneWords, lowerWord)
		if isProfane {
			cleanedWords[i] = "****"
		} else {
			cleanedWords[i] = word
		}
	}

	return strings.Join(cleanedWords, " ")
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
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on http://localhost:%s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
