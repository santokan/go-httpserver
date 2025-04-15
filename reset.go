package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Platform is "dev", proceed with reset operations
	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting all users: %v", err)
		http.Error(w, "Error deleting all users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("Metrics reset and users deleted")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
