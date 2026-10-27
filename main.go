package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	serCfg, apiCfg := getServerConfig(8080)
	defer serCfg.db.Close()

	registerAppHandlers(serCfg.Server.Handler.(*http.ServeMux), apiCfg)
	registerApiHandlers(serCfg.Server.Handler.(*http.ServeMux), apiCfg)
	registerAdminHandlers(serCfg.Server.Handler.(*http.ServeMux), apiCfg)

	log.Printf("Serving on port: %d\n", serCfg.Port)

	if err := serCfg.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func registerAppHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	prefix := "/app"
	mux.Handle(
		prefix+"/",
		http.StripPrefix(prefix+"/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("./app/")))),
	)
}

func registerApiHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	prefix := "/api"
	mux.HandleFunc("GET "+prefix+"/healthz", logRequestMiddlewareFunc(readinessHandler))
	mux.HandleFunc("POST "+prefix+"/validate_chirp", logRequestMiddlewareFunc(chirpHandler))
	mux.HandleFunc("POST "+prefix+"/users", logRequestMiddlewareFunc(apiCfg.handleCreateUser))
}

func registerAdminHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	prefix := "/admin"
	mux.HandleFunc("GET "+prefix+"/metrics", logRequestMiddlewareFunc(apiCfg.metricHandler))
	mux.HandleFunc("POST "+prefix+"/reset", logRequestMiddlewareFunc(apiCfg.handleUserReset))
}
