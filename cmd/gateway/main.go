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
	"github.com/betting-platform/internal/infrastructure/tracing"
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

	// Initialize OpenTelemetry tracing
	tracerCfg := tracing.DefaultConfig("gateway")
	cleanup, err := tracing.InitTracer(ctx, tracerCfg)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// Initialize MaxMind GeoIP provider (optional)
	geoProvider, err := middleware.NewMaxMindProvider(middleware.DefaultDBPath())
	if err != nil {
		logger.Warn("failed to load maxmind database", "error", err)
	}

	// Initialize Redis rate limiter using DefaultConfig (reads REDIS_URL automatically from environment)
	rate
