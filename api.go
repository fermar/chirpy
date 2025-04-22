package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/fermar/chirpy/internal/auth"
	"github.com/fermar/chirpy/internal/database"
)

type RespChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	chid, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "UUID error", err)
		return
	}
	slog.Debug("HIT chirpByID", "ID", chid)
	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "BD error", err)
		return
	}
	respondWithJSON(
		w,
		http.StatusOK,
		RespChirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HIT chirps")
	allChirpsBD, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "BD error", err)
		return
	}
	allChirps := []RespChirp{}
	for _, chirp := range allChirpsBD {
		allChirps = append(
			allChirps,
			RespChirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.CreatedAt,
				Body:      rechirp(chirp.Body),
				UserID:    chirp.UserID,
			},
		)
	}
	respondWithJSON(w, http.StatusOK, allChirps)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}
	type cleaned_chirp struct {
		CleanedBody string `json:"cleaned_body"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header", err)
		return
	}
	usrid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	msg := chirp{}
	err = decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}

	if len(msg.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body: msg.Body,
		// UserID: msg.User_ID,
		UserID: usrid,
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
		RespChirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      rechirp(newChirp.Body),
			UserID:    newChirp.UserID,
		},
	)
}

func (cfg *apiConfig) polkaUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type event struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	slog.Debug("HIT Upgrade User")
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error al leer api key", err)
		return
	}
	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "api key incorrecta", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	msg := event{}
	err = decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}

	if msg.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	updateUuid, err := uuid.Parse(msg.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}
	err = cfg.dbQueries.UpgradeUser(r.Context(), updateUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "BD error", err)
		return
	}
	// respondWithJSON(w, http.StatusOK, cleaned_chirp{CleanedBody: rechirp(msg.Body)})
	w.WriteHeader(http.StatusNoContent)
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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type UsrData struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	// ExpiresIn *int   `json:"expires_in_seconds,omitempty"`
}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HIT revoke")
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header", err)
		return
	}
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating token in DB", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HIT refresh")
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header", err)
		return
	}
	slog.Debug("HIT Refresh", "refresh Token (from auth header)", refreshToken)
	refTokenDb, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header from DB", err)
		return
	}
	accessToken, err := auth.MakeJWT(refTokenDb.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error srv", err)
		return
	}
	type tokenResponse struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, tokenResponse{Token: accessToken})
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	msg := UsrData{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}
	slog.Debug("HIT login", "user email", msg.Email)
	usr, err := cfg.dbQueries.GetUserByEmail(r.Context(), msg.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error DB", err)
		return
	}
	err = auth.CheckPasswordHash(usr.HashedPassword, msg.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error DB", err)
		return
	}
	defaultExp := "3600s"
	expTime, err := time.ParseDuration(defaultExp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error req", err)
		return
	}
	token, err := auth.MakeJWT(usr.ID, cfg.secret, expTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error srv", err)
		return
	}
	refreshToken, _ := auth.MakeRefreshToken()
	crtParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: usr.ID,
	}
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), crtParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error DB", err)
		return
	}
	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:           usr.ID,
			CreatedAt:    usr.CreatedAt,
			UpdatedAt:    usr.UpdatedAt,
			Email:        usr.Email,
			IsChirpyRed:  usr.IsChirpyRed,
			Token:        token,
			RefreshToken: refreshToken,
		},
	)
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	msg := UsrData{}
	err := decoder.Decode(&msg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}
	slog.Debug("HIT createUser", "user email", msg.Email)
	hp, err := auth.HashPassword(msg.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error srv", err)
	}
	cuparams := database.CreateUserParams{
		HashedPassword: hp,
		Email:          msg.Email,
	}
	usr, err := cfg.dbQueries.CreateUser(r.Context(), cuparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error DB", err)
		return
	}
	respondWithJSON(
		w,
		http.StatusCreated,
		User{
			ID:          usr.ID,
			CreatedAt:   usr.CreatedAt,
			UpdatedAt:   usr.UpdatedAt,
			Email:       usr.Email,
			IsChirpyRed: usr.IsChirpyRed,
		},
	)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HIT UPDATE User")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header", err)
		return
	}
	usrid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	newUserData := UsrData{}
	err = decoder.Decode(&newUserData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error en json decode", err)
		return
	}
	hashedpass, err := auth.HashPassword(newUserData.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "srv error", err)
		return
	}
	uuparams := database.UpdateUserParams{
		Email:          newUserData.Email,
		HashedPassword: hashedpass,
		ID:             usrid,
	}
	usr, err := cfg.dbQueries.UpdateUser(r.Context(), uuparams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "BD error", err)
		return
	}
	// respondWithJSON(w, http.StatusOK, cleaned_chirp{CleanedBody: rechirp(msg.Body)})
	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:          usr.ID,
			CreatedAt:   usr.CreatedAt,
			UpdatedAt:   usr.UpdatedAt,
			Email:       usr.Email,
			IsChirpyRed: usr.IsChirpyRed,
			// Token:     token,
		},
	)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chid, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "UUID error", err)
		return
	}
	slog.Debug("HIT Delete Chirp", "ID", chid)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error getting token header", err)
		return
	}
	usrid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "BD error", err)
		return
	}

	if chirp.UserID != usrid {
		respondWithError(w, http.StatusForbidden, "auth error", err)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), chid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "BD error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
