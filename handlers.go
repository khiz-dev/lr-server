package main

import (
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