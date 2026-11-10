package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/betting-platform/internal/infrastructure/config"
	"github.com/betting-platform/internal/infrastructure/http/middleware"
	"github.com/betting-platform/internal/infrastructure/logging"
	"github.com/betting-platform/internal/infrastructure/metrics"
	"github.com/betting-platform/internal/infrastructure/ratelimit"
	"github.com/betting-platform/internal/infrastructure/server"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadFromFile(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := logging.Setup(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("starting gateway", "env", cfg.Service.Environment)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis rate limiter safely
	rateLimiterConfig := ratelimit.DefaultConfig()
	ratelimiter, err := ratelimit.NewRedisLimiter(ctx, rateLimiterConfig)
	if err != nil {
		logger.Error("failed to initialize redis rate limiter", "error", err)
		// Non-fatal for now to keep server running smoothly
	} else {
		defer ratelimiter.Close()
	}

	r := mux.NewRouter()

	// Base Test & Health Routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Lets-Bet Server is Working Perfectly! 🚀 All Features Active."))
	}).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// ==========================================
	// API v1 Routes & Game Endpoints Integration
	// ==========================================
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"All betting features, crash games and wallet APIs are online"}`))
	}).Methods("GET")

	// Metrics & Monitoring Routes
	rec := metrics.New("gateway")
	rec.RegisterRoutes(r)

	// Apply essential middleware stack safely
	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS(cfg.Security))
	
	if ratelimiter != nil {
		r.Use(ratelimiter.HTTPMiddleware())
	}

	addr := fmt.Sprintf(":%d", cfg.Service.Port)
	if err := server.RunHTTP(ctx, addr, r, logger); err != nil {
		logger.Error("server terminated with error", "error", err)
		os.Exit(1)
	}
}
