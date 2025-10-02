package main

import (
	"log"
	"net/http"

	"github.com/MarianGheorghiu/chirpy/api"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	apiCfg := &api.APIConfig{}

	files := http.FileServer(http.Dir("app"))
	dirFilesHandler := http.StripPrefix("/app/", files)
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(dirFilesHandler))

	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerAdminMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("listening on %s ...", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
