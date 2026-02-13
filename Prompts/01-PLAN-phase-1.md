# Katastr Praha 6 — Fáze 1: Základní infrastruktura

**Stav: DOKONČENO**

## Context

Cílem je vytvořit mobilní aplikaci (Android + iOS) pro nahlížení do katastru nemovitostí se zaměřením na Prahu 6. Projekt startoval od nuly.

Fáze 1 pokrývá: monorepo strukturu, Go backend s health endpointem, Docker + Redis, Flutter app s mapovou obrazovkou, a ověření komunikace mezi nimi.

**API klíč ČÚZK zatím není k dispozici** — připraveny placeholdery a mock-ready kód.

---

## Struktura monorepo

```
katastr-p6/
├── backend/                         # Go API server
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── handler/
│   │   │   ├── health.go
│   │   │   └── version.go
│   │   ├── cache/redis.go
│   │   └── middleware/
│   │       ├── logging.go
│   │       └── cors.go
│   ├── Dockerfile
│   ├── Makefile
│   ├── .env.example
│   └── go.mod
├── app/                             # Flutter mobilní aplikace
│   ├── lib/
│   │   ├── main.dart
│   │   ├── app.dart
│   │   ├── config/constants.dart
│   │   ├── screens/map_screen.dart
│   │   ├── services/api_client.dart
│   │   └── widgets/
│   ├── pubspec.yaml
│   └── test/
├── docker-compose.yml
├── .gitignore
├── Makefile
└── Prompts/
```

---

## Implementační kroky

### Krok 1: Git init + .gitignore
- `git init`
- `.gitignore` pro Go, Flutter, IDE, .env, .DS_Store

### Krok 2: Go backend
- `go mod init katastr-p6/backend`
- **Závislosti:** chi/v5, chi/cors, go-redis/v9, godotenv
- **cmd/server/main.go** — chi router, graceful shutdown (SIGINT/SIGTERM), middleware chain (logger, cors)
- **internal/config/config.go** — struct Config, načítání z env: PORT (8080), REDIS_URL, CUZK_API_KEY, CUZK_BASE_URL
- **internal/handler/health.go** — `GET /health` → `{"status":"ok","redis":"connected|disconnected"}`
- **internal/handler/version.go** — `GET /api/version` → `{"version":"0.1.0","backend":"go"}`
- **internal/cache/redis.go** — NewClient, Get/Set s TTL, Ping health check, graceful fallback pokud Redis nedostupný
- **internal/middleware/logging.go** — request logging přes `log/slog` (method, path, status, duration)
- **internal/middleware/cors.go** — AllowAll pro dev
- **.env.example** — šablona s komentáři
- **Makefile** — run, test, build, lint targets

### Krok 3: Docker
- **backend/Dockerfile** — multi-stage: `golang:1.23-alpine` builder → `alpine:3.19` runner
- **docker-compose.yml** — services: backend (8080), redis (6379, alpine 7)

### Krok 4: Flutter app
- `flutter create --org cz.katastrp6 --project-name katastr_p6 app`
- **pubspec.yaml dependencies:**
  - `flutter_map: ^8.2.2` + `latlong2: ^0.9.1` — mapa
  - `http: ^1.2.0` — HTTP klient
  - `geolocator: ^13.0.2` + `permission_handler: ^11.3.1` — GPS
  - `provider: ^6.1.2` — state management
- **lib/main.dart** — Provider setup, runApp
- **lib/app.dart** — MaterialApp, theme, MapScreen jako home
- **lib/config/constants.dart** — apiBaseUrl, Praha 6 bounding box, kódy KÚ
- **lib/screens/map_screen.dart** — FlutterMap s OSM tile layer, AppBar, FAB "Test API"
- **lib/services/api_client.dart** — getVersion(), checkHealth()

### Krok 5: Ověření E2E
- Backend běží (docker compose up nebo go run)
- Flutter app volá /api/version → zobrazí SnackBar s odpovědí

### Krok 6: Root Makefile
- `make up` / `make down` — docker compose
- `make backend-run` — lokální go run
- `make app-run` — flutter run

---

## Verifikace

- [x] `cd backend && go build ./...` — kompiluje bez chyb
- [x] `docker compose up -d && curl http://localhost:8080/health` → `{"status":"ok"}`
- [x] `curl http://localhost:8080/api/version` → `{"version":"0.1.0"}`
- [x] `cd app && flutter analyze` — bez chyb
- [x] `cd app && flutter test` — testy prochází
- [x] `cd app && flutter run` — zobrazí mapu (OSM tiles)
- [x] Klik na FAB "Test API" → SnackBar s verzí z backendu

---

## Technologie a verze

| Komponenta | Knihovna | Verze |
|---|---|---|
| Go router | go-chi/chi | v5.2.5 |
| Go Redis | redis/go-redis | v9.17.3 |
| Go .env | joho/godotenv | v1.5.1 |
| Flutter map | flutter_map | ^8.2.2 |
| Flutter GPS | geolocator | ^13.0.2 |
| Flutter state | provider | ^6.1.2 |
| Flutter HTTP | http | ^1.2.0 |
| Flutter coords | latlong2 | ^0.9.1 |
