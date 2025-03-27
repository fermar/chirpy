package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) // Deshabilitar log cambiando a io.Discard
	slog.SetLogLoggerLevel(slog.LevelDebug)
	port := "8080"
	dir := "."
	apiCfg := apiConfig{}
	apiCfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		apiCfg.midwMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(dir)))),
	)
	// mux.Handle("/healthz/", http.HandlerFunc(readiness))
	mux.HandleFunc("/healthz/", readiness)
	mux.HandleFunc("/metrics/", apiCfg.metrics)
	mux.HandleFunc("/reset/", apiCfg.reset)

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
