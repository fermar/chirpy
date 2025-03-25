package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) // Deshabilitar log cambiando a io.Discard
	port := "8080"
	dir := "."
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(dir))))
	// mux.Handle("/healthz/", http.HandlerFunc(readiness))
	mux.HandleFunc("/healthz/", readiness)

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
