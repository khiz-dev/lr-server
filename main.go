package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func main() {
	serCfg, apiCfg := getServerConfig(8080)

	registerHandlers(serCfg.Server.Handler.(*http.ServeMux), apiCfg)

	log.Printf("Serving on port: %d\n", serCfg.Port)

	if err := serCfg.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func registerHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	mux.Handle(
		"/app/",
		http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))),
	)
	mux.HandleFunc("GET /healthz", logRequestMiddlewareFunc(readinessHandler))
	mux.HandleFunc("GET /metrics", logRequestMiddlewareFunc(apiCfg.metricHandler))
	mux.HandleFunc("POST /reset", logRequestMiddlewareFunc(apiCfg.resetHandler))
}

func (apiCfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d\n", apiCfg.FileserverHits.Load())))
}

func (apiCfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	apiCfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Metrics reset: %d\n", apiCfg.FileserverHits.Load())))
}
