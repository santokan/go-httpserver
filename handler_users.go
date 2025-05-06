package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/go-httpserver/internal/auth"
	"github.com/santokan/go-httpserver/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := createUserRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if params.Password == "" || params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type updatePasswordRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token: ", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := updatePasswordRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if params.Password == "" || params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	dbParams := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.UpdateUser(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
