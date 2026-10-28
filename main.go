package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	serCfg, apiCfg := getServerConfig(8080)

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
	endpoints := getEndpoints(apiCfg)
	for _, endpoint := range endpoints {
		for pattern, handler := range endpoint {
			mux.HandleFunc(fmt.Sprintf(pattern, prefix), logRequestMiddlewareFunc(handler))
		}
	}
}

func getEndpoints(apiCfg *apiConfig) []map[string]http.HandlerFunc {
	return []map[string]http.HandlerFunc{
		{
			"GET %s/healthz": readinessHandler,
			"POST %s/users":  apiCfg.handleCreateUser,
			"POST %s/login":  apiCfg.handleLogin,
			"POST %s/validate_chirp": chirpHandler,
			"POST %s/chirps": apiCfg.handleCreateChirp,
			"GET %s/chirps": apiCfg.handleGetAllChrips,
			"GET %s/chirps/{chirpID}": apiCfg.handleGetChirpByID,
		},
	}
}

func registerAdminHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	prefix := "/admin"
	mux.HandleFunc("GET "+prefix+"/metrics", logRequestMiddlewareFunc(apiCfg.metricHandler))
	mux.HandleFunc("POST "+prefix+"/reset", logRequestMiddlewareFunc(apiCfg.handleUserReset))
}
