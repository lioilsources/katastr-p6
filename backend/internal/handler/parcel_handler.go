package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"katastr-p6/backend/internal/cache"
	"katastr-p6/backend/internal/coords"
	"katastr-p6/backend/internal/cuzk"
)

// ParcelHandler handles parcel-related API endpoints.
type ParcelHandler struct {
	client *cuzk.Client
	ch     *CachedHandler
}

// NewParcelHandler creates a new ParcelHandler.
func NewParcelHandler(client *cuzk.Client, c *cache.RedisCache) *ParcelHandler {
	return &ParcelHandler{
		client: client,
		ch:     NewCachedHandler(c),
	}
}

// Search handles GET /api/parcels/search?area={code}&number={num}
func (h *ParcelHandler) Search(w http.ResponseWriter, r *http.Request) {
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

	key := CacheKey("parcels:search", areaCode, number)
	data, err := h.ch.GetOrFetch(r.Context(), key, 1*time.Minute, func() (any, error) {
		return h.client.SearchParcels(r.Context(), areaCode, number)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Get handles GET /api/parcels/{id}
func (h *ParcelHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("parcel", id)
	data, err := h.ch.GetOrFetch(r.Context(), key, 5*time.Minute, func() (any, error) {
		return h.client.GetParcel(r.Context(), id)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Polygon handles GET /api/parcels/polygon?lat={lat}&lon={lon}&radius={m}
func (h *ParcelHandler) Polygon(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radiusStr := r.URL.Query().Get("radius")

	if latStr == "" || lonStr == "" {
		http.Error(w, `{"error":"missing required parameters: lat, lon"}`, http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid lat"}`, http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid lon"}`, http.StatusBadRequest)
		return
	}

	radius := 5 // default 5m
	if radiusStr != "" {
		if r, err := strconv.Atoi(radiusStr); err == nil {
			radius = r
		}
	}

	// Convert WGS-84 to S-JTSK for the CUZK API.
	x, y := coords.WGS84ToSJTSK(lat, lon)

	key := CacheKey("parcels:polygon", x, y, radius)
	data, err := h.ch.GetOrFetch(r.Context(), key, 1*time.Minute, func() (any, error) {
		return h.client.PolygonParcels(r.Context(), x, y, radius)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// Neighbors handles GET /api/parcels/neighbors/{id}
func (h *ParcelHandler) Neighbors(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	key := CacheKey("parcels:neighbors", id)
	data, err := h.ch.GetOrFetch(r.Context(), key, 5*time.Minute, func() (any, error) {
		return h.client.NeighborParcels(r.Context(), id)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
