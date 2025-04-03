package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/fermar/chirpy/internal/database"
)

func (cfg *apiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}
	type cleaned_chirp struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type respChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	msg := chirp{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}

	if len(msg.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   msg.Body,
		UserID: msg.User_ID,
	}
	newChirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "BD error", err)
		return
	}
	// respondWithJSON(w, http.StatusOK, cleaned_chirp{CleanedBody: rechirp(msg.Body)})
	respondWithJSON(
		w,
		http.StatusCreated,
		respChirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      rechirp(newChirp.Body),
			UserID:    newChirp.UserID,
		},
	)
}

func rechirp(body string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	listaPalabras := strings.Split(body, " ")
	for i, p := range listaPalabras {
		if slices.Contains(bannedWords, strings.ToLower(p)) {
			listaPalabras[i] = "****"
		}
	}
	return strings.Join(listaPalabras, " ")
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type usrData struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	msg := usrData{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}
	slog.Debug("HIT createUser", "user email", msg.Email)
	usr, err := cfg.dbQueries.CreateUser(r.Context(), msg.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error DB", err)
		return
	}
	respondWithJSON(
		w,
		http.StatusCreated,
		User{ID: usr.ID, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt, Email: usr.Email},
	)
}
