package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	type validResp struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	msg := chirp{}
	err := decoder.Decode(&msg)
	if err != nil {
		slog.Error("no se puede decodificar chirp", "err", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json") // normal header

	if len(msg.Body) > 140 {
		dat, err := json.Marshal(errorResp{Error: "Chirp is too long"})
		if err != nil {
			slog.Error("error marshalling", "err", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	dat, err := json.Marshal(validResp{Valid: true})
	if err != nil {
		slog.Error("error marshalling", "err", err)
		w.WriteHeader(500)
		return
	}
	// w.Header().Add("Content-Type", "application/json") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
