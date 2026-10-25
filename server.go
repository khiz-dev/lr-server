package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type cfg struct {
	Port   int
	Server *http.Server
}

type apiConfig struct {
	FileserverHits atomic.Int32
}

func getServerConfig(port int) (*cfg, *apiConfig) {
	mux := http.NewServeMux()

	apiCfg := apiConfig{
		FileserverHits: atomic.Int32{},
	}

	cfg := &cfg{
		Port: port,
	}

	cfg.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	return cfg, &apiCfg
}
