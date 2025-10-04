package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MarianGheorghiu/chirpy/api"
	"github.com/MarianGheorghiu/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db.Ping: %v", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apiCfg := &api.APIConfig{
		Queries:  dbQueries,
		Platform: platform,
	}

	files := http.FileServer(http.Dir("app"))
	dirFilesHandler := http.StripPrefix("/app/", files)
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(dirFilesHandler))

	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.HandlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)

	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerChirpCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerChirpsList)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerChirpGet)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerAdminMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Printf("listening on %s ...", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
