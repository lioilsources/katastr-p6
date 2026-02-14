package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"katastr-p6/backend/internal/cache"
	"katastr-p6/backend/internal/cuzk"
)

// ProceedingHandler handles proceeding-related API endpoints.
type ProceedingHandler struct {
	client *cuzk.Client
	ch     *CachedHandler
}

// NewProceedingHandler creates a new ProceedingHandler.
func NewProceedingHandler(client *cuzk.Client, c *cache.RedisCache) *ProceedingHandler {
	return &ProceedingHandler{
		client: client,
		ch:     NewCachedHandler(c),
	}
}

// Get handles GET /api/proceedings/{id}
func (h *ProceedingHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("proceeding", id)
	data, err := h.ch.GetOrFetch(r.Context(), key, 5*time.Minute, func() (any, error) {
		return h.client.GetProceeding(r.Context(), id)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
