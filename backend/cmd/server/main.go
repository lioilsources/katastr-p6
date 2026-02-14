package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"katastr-p6/backend/internal/cache"
	"katastr-p6/backend/internal/config"
	"katastr-p6/backend/internal/cuzk"
	"katastr-p6/backend/internal/handler"
	"katastr-p6/backend/internal/middleware"
)

func main() {
	cfg := config.Load()

	// Redis (optional — graceful fallback)
	var redisCache *cache.RedisCache
	if cfg.RedisURL != "" {
		redisCache = cache.NewRedisCache(cfg.RedisURL)
		defer redisCache.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisCache.Ping(ctx); err != nil {
			slog.Warn("redis not available, running without cache", "error", err)
			redisCache = nil
		} else {
			slog.Info("redis connected", "addr", cfg.RedisURL)
		}
	}

	// CUZK API client
	cuzkClient := cuzk.NewClient(cfg.CUZKBaseURL, cfg.CUZKAPIKey)
	if cfg.CUZKAPIKey == "" {
		slog.Warn("CUZK_API_KEY not set, API calls to CUZK will fail")
	}

	// Handlers
	healthHandler := handler.NewHealthHandler(redisCache)
	parcelHandler := handler.NewParcelHandler(cuzkClient, redisCache)
	buildingHandler := handler.NewBuildingHandler(cuzkClient, redisCache)
	unitHandler := handler.NewUnitHandler(cuzkClient, redisCache)
	proceedingHandler := handler.NewProceedingHandler(cuzkClient, redisCache)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(middleware.CORS()))

	r.Get("/health", healthHandler.Health)
	r.Route("/api", func(r chi.Router) {
		r.Get("/version", handler.Version)

		// Parcels
		r.Get("/parcels/search", parcelHandler.Search)
		r.Get("/parcels/polygon", parcelHandler.Polygon)
		r.Get("/parcels/neighbors/{id}", parcelHandler.Neighbors)
		r.Get("/parcels/{id}", parcelHandler.Get)

		// Buildings
		r.Get("/buildings/search", buildingHandler.Search)
		r.Get("/buildings/{id}", buildingHandler.Get)

		// Units
		r.Get("/units/search", unitHandler.Search)
		r.Get("/units/{id}", unitHandler.Get)

		// Proceedings
		r.Get("/proceedings/{id}", proceedingHandler.Get)
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
