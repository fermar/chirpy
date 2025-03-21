package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) // Deshabilitar log cambiando a io.Discard
	mux := http.NewServeMux()
	httpsrv := &http.Server{}
	httpsrv.Handler = mux
	port := "8080"
	httpsrv.Addr = ":" + port
	log.Printf("escuchando en puerto %v...\n", port)
	log.Fatal(httpsrv.ListenAndServe())
}
