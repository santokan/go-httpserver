package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		filePathRoot = "."
		port         = "8080"
	)
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/healthz", http.HandlerFunc(healthCheck))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from '%s' on http://localhost:%s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
