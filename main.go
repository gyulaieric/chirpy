package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := "8080"
	filepathRoot := http.Dir(".")

	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	// API
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	// ADMIN
	mux.Handle("POST /admin/reset", apiCfg.handlerResetMetrics())

	// File Server
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix(
				"/app",
				http.FileServer(filepathRoot),
			),
		),
	)

	// ADMIN
	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
