# Katastr Praha 6 — Phase 2: CUZK REST API Client in Go

**Status: COMPLETED**

## Context

Phase 1 done — Go backend with chi router (health + version), Redis cache, Docker, Flutter app with map. Now we need to connect backend to CUZK REST API: HTTP client, data models, endpoint handlers, S-JTSK coordinate transform, and Redis cache integration. No CUZK API key yet — server must start without it.

**Convention: All code, file names, variable names, and comments in English.**

---

## New files

```
backend/internal/
├── coords/
│   ├── transform.go          # WGS-84 <-> S-JTSK conversion
│   └── transform_test.go     # Unit test with Prague coordinates
├── cuzk/
│   ├── client.go             # HTTP client, retry, rate limit
│   ├── models.go             # Parcel, Building, Unit, Proceeding structs
│   ├── parcels.go            # SearchParcels, GetParcel, PolygonParcels, NeighborParcels
│   ├── buildings.go          # SearchBuildings, GetBuilding, AddressPoint
│   ├── units.go              # SearchUnits, GetUnit
│   └── proceedings.go        # GetProceeding
└── handler/
    ├── cache_helper.go       # Shared cache get-or-fetch logic
    ├── parcel_handler.go     # GET /api/parcels/...
    ├── building_handler.go   # GET /api/buildings/...
    ├── unit_handler.go       # GET /api/units/...
    └── proceeding_handler.go # GET /api/proceedings/...
```

Modified: `cmd/server/main.go` (new routes + CUZK client init)

---

## Implementation blocks (bottom-up, each compiles)

### Block 1: Coordinate transform
- `go get github.com/wroge/wgs84/v2` (zero-dependency Go library)
- **coords/transform.go** — `WGS84ToSJTSK(lat, lon)` -> positive (x, y), `SJTSKToWGS84(x, y)` -> (lat, lon)
- **coords/transform_test.go** — round-trip test with Prague point (Old Town Square ~ S-JTSK 1043200, 744700)
- Note: S-JTSK has negative coords, CUZK API expects positive -> `math.Abs()`

### Block 2: CUZK models
- **cuzk/models.go** — Go structs with JSON tags:
  - `Parcel` (ID, BaseNumber, Subdivision, NumberingType, CadastralArea, Area, LandType, UsageType, OwnershipSheetNo, ReferencePoint)
  - `Building` (ID, DescriptiveNo, EvidenceNo, BuildingType, MunicipalPart, CadastralArea, UsageType)
  - `Unit` (ID, UnitNumber, UnitType, CommonPartsShare, BuildingID)
  - `Proceeding` (ID, SequenceNumber, Year, Office, Status, Type, FilingDate)
  - Response wrappers: `{Entity}SearchResponse` with `[]Entity` + `Total`
- Note: Structs are approximations — will adjust after connecting to real API

### Block 3: CUZK HTTP client
- `go get golang.org/x/time` (rate limiter)
- **cuzk/client.go** — `NewClient(baseURL, apiKey)`:
  - HTTP header `Api-Key: {key}`
  - Timeout 10s per request
  - Rate limit: 1 req/s (60/min)
  - Retry: max 3 attempts, exponential backoff (1s, 2s, 4s), retry on 5xx/429
  - Private methods `do(ctx, method, path)` -> `[]byte` and `get(ctx, path, target)` -> JSON decode

### Block 4: Endpoint wrappers
- **cuzk/parcels.go** — `SearchParcels(ctx, areaCode, number)`, `GetParcel(ctx, id)`, `PolygonParcels(ctx, x, y, radius)`, `NeighborParcels(ctx, id)`
- **cuzk/buildings.go** — `SearchBuildings(ctx, areaCode, number)`, `GetBuilding(ctx, id)`, `AddressPoint(ctx, id)`
- **cuzk/units.go** — `SearchUnits(ctx, areaCode, buildingNo, unitNo)`, `GetUnit(ctx, id)`
- **cuzk/proceedings.go** — `GetProceeding(ctx, id)`

### Block 5: Cache helper + handlers
- **handler/cache_helper.go** — shared `CachedHandler` struct with `GetOrFetch(ctx, key, ttl, fallback)`:
  - Cache hit -> return []byte
  - Cache miss -> call fallback, serialize JSON, store in cache, return []byte
  - Redis unavailable -> skip cache, call API directly
  - Cache key pattern: `cuzk:{entity}:{sha256(params)[:16]}`
- **handler/parcel_handler.go** — 4 endpoints:
  - `GET /api/parcels/search?area={code}&number={num}` (TTL 1min)
  - `GET /api/parcels/{id}` (TTL 5min)
  - `GET /api/parcels/polygon?lat={lat}&lon={lon}&radius={m}` — WGS->SJTSK internally (TTL 1min)
  - `GET /api/parcels/neighbors/{id}` (TTL 5min)
- **handler/building_handler.go** — `GET /api/buildings/search`, `GET /api/buildings/{id}`
- **handler/unit_handler.go** — `GET /api/units/search`, `GET /api/units/{id}`
- **handler/proceeding_handler.go** — `GET /api/proceedings/{id}`

### Block 6: Integration into main.go
- Init `cuzk.NewClient(cfg.CUZKBaseURL, cfg.CUZKAPIKey)` + warning if key missing
- Init all handlers with cuzkClient + redisCache
- Register routes in `/api` group
- `go build ./...` + `go test ./...`

---

## API endpoints (result)

| Method | Path | Description | Cache TTL |
|---|---|---|---|
| GET | /api/parcels/search?area=&number= | Search parcels | 1 min |
| GET | /api/parcels/{id} | Parcel detail | 5 min |
| GET | /api/parcels/polygon?lat=&lon=&radius= | Parcels near point (WGS-84 input) | 1 min |
| GET | /api/parcels/neighbors/{id} | Neighbor parcels | 5 min |
| GET | /api/buildings/search?area=&number= | Search buildings | 1 min |
| GET | /api/buildings/{id} | Building detail | 5 min |
| GET | /api/units/search?area=&buildingNo=&unitNo= | Search units | 1 min |
| GET | /api/units/{id} | Unit detail | 5 min |
| GET | /api/proceedings/{id} | Proceeding detail | 5 min |

---

## New Go dependencies

- `github.com/wroge/wgs84/v2` — S-JTSK <-> WGS-84 (zero-dependency)
- `golang.org/x/time` — rate limiter for CUZK client

---

## Verification

- [x] `cd backend && go build ./...` — compiles without errors
- [x] `go vet ./...` — no issues
- [x] `go test ./internal/coords/... -v` — coordinate round-trip test passes (4/4)
- [x] Server starts without CUZK_API_KEY (warning in log)
- [x] All 9 new endpoints registered and reachable
- [x] `curl localhost:8080/health` — still works
