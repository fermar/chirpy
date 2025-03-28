package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	slog.Warn("respondiendo error", "code", code, "msg", msg, "error", err)
	type errorResp struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResp{Error: msg})
	// resp, err := json.Marshal(errorResp{Error: msg})
	// if err != nil {
	// 	slog.Error("error en marshal", "error", err)
	// 	resp = []byte("")
	// }
	//
	// w.WriteHeader(code)
	// w.Write(resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		slog.Error("error en marshal", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
