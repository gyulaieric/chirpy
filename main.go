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
	jwtSecret      string
	platform       string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Couldn't load .env: %v", err)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf(`Couldn't connect to database at "%s": %v`, dbURL, err)
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	apiCfg := apiConfig{
		db:        database.New(db),
		jwtSecret: jwtSecret,
		platform:  platform,
	}

	port := "8080"
	filepathRoot := http.Dir(".")

	mux := http.NewServeMux()

	// API
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("POST /api/users", apiCfg.handlerRegister())
	mux.Handle("PUT /api/users", apiCfg.handlerUpdateUsers())
	mux.Handle("POST /api/login", apiCfg.handlerLogin())
	mux.Handle("POST /api/refresh", apiCfg.handlerRefresh())
	mux.Handle("POST /api/revoke", apiCfg.handlerRevoke())
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
