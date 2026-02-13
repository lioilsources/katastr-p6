package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"katastr-p6/backend/internal/cache"
)

type HealthHandler struct {
	cache *cache.RedisCache
}

func NewHealthHandler(c *cache.RedisCache) *HealthHandler {
	return &HealthHandler{cache: c}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	redisStatus := "connected"
	if h.cache != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := h.cache.Ping(ctx); err != nil {
			redisStatus = "disconnected"
		}
	} else {
		redisStatus = "not configured"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"redis":  redisStatus,
	})
}
