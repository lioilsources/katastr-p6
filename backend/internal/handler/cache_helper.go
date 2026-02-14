package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"katastr-p6/backend/internal/cache"
)

// CachedHandler provides cache-through helper for API handlers.
type CachedHandler struct {
	cache *cache.RedisCache
}

// NewCachedHandler creates a new CachedHandler. cache can be nil (no caching).
func NewCachedHandler(c *cache.RedisCache) *CachedHandler {
	return &CachedHandler{cache: c}
}

// CacheKey builds a deterministic cache key from a prefix and parameters.
func CacheKey(prefix string, params ...any) string {
	h := sha256.New()
	for _, p := range params {
		fmt.Fprintf(h, "%v:", p)
	}
	return fmt.Sprintf("cuzk:%s:%x", prefix, h.Sum(nil)[:8])
}

// GetOrFetch tries the cache first; on miss calls fallback, caches the result, and returns JSON bytes.
func (ch *CachedHandler) GetOrFetch(ctx context.Context, key string, ttl time.Duration, fallback func() (any, error)) ([]byte, error) {
	// Try cache
	if ch.cache != nil {
		cached, err := ch.cache.Get(ctx, key)
		if err == nil {
			return []byte(cached), nil
		}
		// Log non-miss errors but continue without cache.
		if cached == "" && err.Error() != "redis: nil" {
			slog.Warn("cache get error", "key", key, "error", err)
		}
	}

	// Call the upstream API
	data, err := fallback()
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	// Store in cache (best effort)
	if ch.cache != nil {
		if setErr := ch.cache.Set(ctx, key, string(jsonData), ttl); setErr != nil {
			slog.Warn("cache set error", "key", key, "error", setErr)
		}
	}

	return jsonData, nil
}
