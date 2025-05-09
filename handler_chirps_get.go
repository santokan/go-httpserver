package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/santokan/go-httpserver/internal/database"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error

	authorIdStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if authorIdStr != "" {
		var userUuid uuid.UUID
		userUuid, err = uuid.Parse(authorIdStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUser(r.Context(), userUuid)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get chirps", err)
		return
	}

	// Only need to handle descending order explicitly since the DB query
	// already sorts in ascending order by default
	if sortOrder == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found!", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
