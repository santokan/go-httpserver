package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/santokan/go-httpserver/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type webhookRequest struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil || apikey != cfg.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := webhookRequest{}
	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.SetUserPremiumByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
