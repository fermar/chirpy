package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/fermar/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
	apiKey         string
}

func (cfg *apiConfig) midwMetricsInc(next http.Handler) http.Handler {
	// cfg.fileserverHits.Add(1)
	// slog.Info("app", "hits", cfg.fileserverHits.Load())
	// return next
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		slog.Debug("app", "hits", cfg.fileserverHits.Load())
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	slog.Debug("stats", "hits", cfg.fileserverHits.Load())
	content := fmt.Sprintf(`<html> 
		<body> 
		<h1> Welcome, Chirpy Admin</h1>`+
		"<p>Chirpy has been visited %d times!</p>"+
		`</body>
		</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(content))
	// w.Write([]byte("Hits: " + strconv.FormatInt(int64(cfg.fileserverHits.Load()), 10)))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HIT reset")
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "", nil)
		return
	}
	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error DB", err)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	slog.Debug("reset", "hits", cfg.fileserverHits.Load())
	w.Write([]byte("Reset Hits"))
}
