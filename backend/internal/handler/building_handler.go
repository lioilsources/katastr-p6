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

// BuildingHandler handles building-related API endpoints.
type BuildingHandler struct {
	client *cuzk.Client
	ch     *CachedHandler
}

// NewBuildingHandler creates a new BuildingHandler.
func NewBuildingHandler(client *cuzk.Client, c *cache.RedisCache) *BuildingHandler {
	return &BuildingHandler{
		client: client,
		ch:     NewCachedHandler(c),
	}
}

// Search handles GET /api/buildings/search?area={code}&number={num}
func (h *BuildingHandler) Search(w http.ResponseWriter, r *http.Request) {
	areaStr := r.URL.Query().Get("area")
	number := r.URL.Query().Get("number")

	if areaStr == "" || number == "" {
		http.Error(w, `{"error":"missing required parameters: area, number"}`, http.StatusBadRequest)
		return
	}

	areaCode, err := strconv.Atoi(areaStr)
	if err != nil {
		http.Error(w, `{"error":"invalid area parameter"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("buildings:search", areaCode, number)
	data, err := h.ch.GetOrFetch(r.Context(), key, 1*time.Minute, func() (any, error) {
		return h.client.SearchBuildings(r.Context(), areaCode, number)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Get handles GET /api/buildings/{id}
func (h *BuildingHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("building", id)
	data, err := h.ch.GetOrFetch(r.Context(), key, 5*time.Minute, func() (any, error) {
		return h.client.GetBuilding(r.Context(), id)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
