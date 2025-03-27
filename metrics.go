package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	slog.Debug("reset", "hits", cfg.fileserverHits.Load())
	w.Write([]byte("Reset Hits"))
}
