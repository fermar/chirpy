package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) // Deshabilitar log cambiando a io.Discard
	port := "8080"
	dir := "./"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dir)))
	httpsrv := &http.Server{}
	httpsrv.Handler = mux
	httpsrv.Addr = ":" + port
	log.Printf("escuchando en puerto %v\nsirviendo archivos en %v\n", port, dir)
	log.Fatal(httpsrv.ListenAndServe())
}
