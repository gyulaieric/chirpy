package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gyulaieric/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Couldn't load .env: %v", err)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf(`Couldn't connect to database at "%s": %v`, dbURL, err)
	}
	apiCfg := apiConfig{
		db:       database.New(db),
		platform: os.Getenv("PLATFORM"),
	}

	port := "8080"
	filepathRoot := http.Dir(".")

	mux := http.NewServeMux()

	// API
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("POST /api/users", apiCfg.handlerRegister())
	mux.Handle("POST /api/login", apiCfg.handlerLogin())
	mux.Handle("GET /api/chirps", apiCfg.handlerGetChirps())
	mux.Handle("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp())
	mux.Handle("POST /api/chirps", apiCfg.handlerCreateChirp())
	// ADMIN
	mux.Handle("POST /admin/reset", apiCfg.handlerReset())
	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())

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

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
