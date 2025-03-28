package main

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type validResp struct {
		Valid bool `json:"valid"`
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

	respondWithJSON(w, http.StatusOK, validResp{Valid: true})
}
