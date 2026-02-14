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

// UnitHandler handles unit-related API endpoints.
type UnitHandler struct {
	client *cuzk.Client
	ch     *CachedHandler
}

// NewUnitHandler creates a new UnitHandler.
func NewUnitHandler(client *cuzk.Client, c *cache.RedisCache) *UnitHandler {
	return &UnitHandler{
		client: client,
		ch:     NewCachedHandler(c),
	}
}

// Search handles GET /api/units/search?area={code}&buildingNo={bn}&unitNo={un}
func (h *UnitHandler) Search(w http.ResponseWriter, r *http.Request) {
	areaStr := r.URL.Query().Get("area")
	buildingNo := r.URL.Query().Get("buildingNo")
	unitNo := r.URL.Query().Get("unitNo")

	if areaStr == "" || buildingNo == "" || unitNo == "" {
		http.Error(w, `{"error":"missing required parameters: area, buildingNo, unitNo"}`, http.StatusBadRequest)
		return
	}

	areaCode, err := strconv.Atoi(areaStr)
	if err != nil {
		http.Error(w, `{"error":"invalid area parameter"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("units:search", areaCode, buildingNo, unitNo)
	data, err := h.ch.GetOrFetch(r.Context(), key, 1*time.Minute, func() (any, error) {
		return h.client.SearchUnits(r.Context(), areaCode, buildingNo, unitNo)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Get handles GET /api/units/{id}
func (h *UnitHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("unit", id)
	data, err := h.ch.GetOrFetch(r.Context(), key, 5*time.Minute, func() (any, error) {
		return h.client.GetUnit(r.Context(), id)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
