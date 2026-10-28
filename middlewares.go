package main

import (
	"context"
	"log"
	"net/http"

	"github.com/khiz-dev/goserve/internal/auth"
)

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func logRequestMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next(w, r)
	}
}

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apiCfg *apiConfig) middlewareAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		jwtToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(jwtToken, apiCfg.secret)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// You can set the userID in the request context if needed
		ctx := context.WithValue(r.Context(), "userID", userID.String())
		next(w, r.WithContext(ctx))
	}
}
