package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type cleaned_chirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// type validResp struct {
	// 	Valid bool `json:"valid"`
	// }

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

	respondWithJSON(w, http.StatusOK, cleaned_chirp{CleanedBody: rechirp(msg.Body)})
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
