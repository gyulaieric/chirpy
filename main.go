package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	filepathRoot := http.Dir(".")

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(filepathRoot)))
	mux.HandleFunc("/healthz", handlerHealthz)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
