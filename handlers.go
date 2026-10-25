package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	tmpl, err := getTemplate(AdminMetricsTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderedContent, err := tmpl.render(apiCfg.FileserverHits.Load())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(renderedContent)
}

func (apiCfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	apiCfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Metrics reset: %d\n", apiCfg.FileserverHits.Load())))
}

func chirpHandler(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
	}
	b := body{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		resp := errorResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if len(b.Body) > 140 {
		resp := errorResponse{Error: "Chirp body exceeds 140 characters"}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := successResponse{Valid: true}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}