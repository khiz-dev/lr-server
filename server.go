package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"github.com/khiz-dev/goserve/internal/database"
	_ "github.com/lib/pq"
)

type cfg struct {
	Port   int
	Server *http.Server
	db     *sql.DB
}

type apiConfig struct {
	FileserverHits atomic.Int32
	db             *database.Queries
	secret         string
	polkaKey       string
}

func getServerConfig(port int) (*cfg, *apiConfig) {
	mux := http.NewServeMux()

	apiCfg := apiConfig{
		FileserverHits: atomic.Int32{},
		db: nil,
		secret: os.Getenv("SECRET"),
		polkaKey: os.Getenv("POLKA_KEY"),
	}

	cfg := &cfg{
		Port: port,
	}

	cfg.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	db, err := getDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg.db = db
	dbQueries := database.New(db)
	apiCfg.db = dbQueries

	return cfg, &apiCfg
}

func getDatabaseConfig() (*sql.DB, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL is not set; add it to a .env file or export it in your shell")
	}

	// lib/pq defaults to sslmode=require. Local Postgres containers usually
	// don't have SSL enabled, so default to disable for local connections.
	if parsed, err := url.Parse(dbURL); err == nil && !parsed.Query().Has("sslmode") {
		host := parsed.Hostname()
		if host == "localhost" || host == "127.0.0.1" || host == "::1" {
			query := parsed.Query()
			query.Set("sslmode", "disable")
			parsed.RawQuery = query.Encode()
			dbURL = parsed.String()
		}
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
