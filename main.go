package main

import (
	"errors"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	serCfg, apiCfg := getServerConfig(8080)

	mux := serCfg.Server.Handler.(*http.ServeMux)

	registerAppHandlers(mux, apiCfg)
	registerApiHandlers(mux, apiCfg)
	registerAdminHandlers(mux, apiCfg)

	log.Printf("Serving on port: %d\n", serCfg.Port)

	if err := serCfg.Server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func registerAppHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	mux.Handle(
		"/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix(
				"/app/",
				http.FileServer(http.Dir("./app")),
			),
		),
	)
}

func registerApiHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	endpoints := getEndpoints(apiCfg)

	// Register authenticated endpoints
	for _, endpoint := range endpoints.Auth {
		mux.HandleFunc(
			endpoint.Pattern,
			apiCfg.middlewareAuthenticatedUser(endpoint.Handler),
		)
	}

	// Register public endpoints
	for _, endpoint := range endpoints.Public {
		mux.HandleFunc(
			endpoint.Pattern,
			logRequestMiddlewareFunc(endpoint.Handler),
		)
	}
}

func registerAdminHandlers(mux *http.ServeMux, apiCfg *apiConfig) {
	mux.HandleFunc(
		"GET /admin/metrics",
		logRequestMiddlewareFunc(apiCfg.metricHandler),
	)

	mux.HandleFunc(
		"POST /admin/reset",
		logRequestMiddlewareFunc(apiCfg.handleUserReset),
	)
}

type Endpoint struct {
	Pattern string
	Handler http.HandlerFunc
}

type Endpoints struct {
	Auth   []Endpoint
	Public []Endpoint
}

func getEndpoints(apiCfg *apiConfig) Endpoints {
	return Endpoints{
		Auth: []Endpoint{
			{
				Pattern: "POST /api/chirps",
				Handler: apiCfg.handleCreateChirp,
			},
			{
				Pattern: "PUT /api/users",
				Handler: apiCfg.handleUserUpdate,
			},
			{
				Pattern: "DELETE /api/chirps/{chirpID}",
				Handler: apiCfg.handleDeleteChirp,
			},
		},

		Public: []Endpoint{
			{
				Pattern: "POST /api/login",
				Handler: apiCfg.handleLogin,
			},
			{
				Pattern: "POST /api/refresh",
				Handler: apiCfg.handleRefreshToken,
			},
			{
				Pattern: "POST /api/revoke",
				Handler: apiCfg.handleRevokeRefreshToken,
			},
			{
				Pattern: "POST /api/validate_chirp",
				Handler: chirpHandler,
			},
			{
				Pattern: "GET /api/healthz",
				Handler: readinessHandler,
			},
			{
				Pattern: "POST /api/users",
				Handler: apiCfg.handleCreateUser,
			},
			{
				Pattern: "GET /api/chirps",
				Handler: apiCfg.handleGetAllChrips,
			},
			{
				Pattern: "GET /api/chirps/{chirpID}",
				Handler: apiCfg.handleGetChirpByID,
			},
			{
				Pattern: "POST /api/polka/webhooks",
				Handler: apiCfg.handleWebhook,
			},
		},
	}
}
