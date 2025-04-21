package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/fermar/chirpy/internal/database"
)

func main() {
	log.SetOutput(os.Stderr) // Deshabilitar log cambiando a io.Discard
	slog.SetLogLoggerLevel(slog.LevelDebug)
	port := "8080"
	dir := "."
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	slog.Debug("env", "dbURL", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("bd", "err", err)
		os.Exit(1)
	}
	apiCfg := apiConfig{}
	apiCfg.fileserverHits.Store(0)
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.secret = os.Getenv("JWTSECRET")
	slog.Debug("env", "PLATFORM", apiCfg.platform)
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		apiCfg.midwMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(dir)))),
	)
	// mux.Handle("/healthz/", http.HandlerFunc(readiness))
	mux.HandleFunc("GET /admin/healthz", readiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)
	// mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByID)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/login", apiCfg.login)
	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.revoke)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)

	httpsrv := &http.Server{}
	httpsrv.Handler = mux
	httpsrv.Addr = ":" + port
	fmt.Printf("escuchando en puerto %v\nsirviendo archivos en %v\n", port, dir)
	log.Fatal(httpsrv.ListenAndServe())
}

func readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
