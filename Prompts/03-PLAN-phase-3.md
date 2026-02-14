# Katastr Praha 6 — Phase 3: Map Display (P0 scope)

**Status: COMPLETED**

## Context

Phase 2 done — Go backend has CUZK REST API client, coordinate transform, 9 endpoints with Redis cache. Flutter app has basic map with OSM tiles, apiClient (version+health only). Now we add: CUZK cadastral tile overlay, GPS localization, tap-to-identify parcels. No backend changes needed — all work is Flutter + platform config. CUZK API key still missing — WMTS tiles work without it, tap-to-identify will fail gracefully.

**Convention: All code, file names, variable names, and comments in English.**

---

## Key decisions

1. **WMTS direct from Flutter** — no backend proxy. CUZK WMTS is free, no auth, EPSG:3857 (native for flutter_map).
2. **No provider yet** — all state is local to MapScreen (user position, overlay toggle, parcel result). Provider can be added in Phase 4+ when cross-screen state is needed.
3. **Location service in separate file** — encapsulates geolocator permissions + position retrieval.
4. **Bottom sheet inline** — simple `showModalBottomSheet` in map_screen.dart, no separate widget file.
5. **Parcel model in `models/parcel.dart`** — typed Dart class for JSON deserialization.

---

## Files

```
app/lib/
  config/constants.dart          # MODIFY: add cuzkWmtsUrl
  models/
    parcel.dart                  # NEW: Parcel + ParcelSearchResult classes
  services/
    api_client.dart              # MODIFY: add searchParcelsByPoint()
    location_service.dart        # NEW: GPS permission + getCurrentPosition()
  screens/
    map_screen.dart              # MODIFY: WMTS layer, GPS, tap handler, bottom sheet

app/android/app/src/main/AndroidManifest.xml  # MODIFY: location permissions
app/ios/Runner/Info.plist                      # MODIFY: NSLocation usage descriptions
```

---

## Implementation blocks (bottom-up, each compiles)

### Block 1: Platform permissions
- **AndroidManifest.xml** — add `ACCESS_FINE_LOCATION` + `ACCESS_COARSE_LOCATION` before `<application>`
- **Info.plist** — add `NSLocationWhenInUseUsageDescription`

### Block 2: WMTS cadastral overlay
- **constants.dart** — add `cuzkWmtsUrl` constant:
  `https://services.cuzk.gov.cz/wmts/local-km-wmts-google.asp?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=katuze_barvy&STYLE=default&FORMAT=image/png&TILEMATRIXSET=googlemapscompatible&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}`
- **map_screen.dart** — add second `TileLayer` with `tileBuilder` for 0.7 opacity (flutter_map v8 has no `opacity` param), `minZoom: 10`, `maxZoom: 20`, toggle via `_showCadastralOverlay` bool

### Block 3: Parcel data model
- **NEW `models/parcel.dart`** — `CadastralArea`, `Parcel` (with `fromJson` using CUZK JSON field names like `kmenoveCislo`, `katastralniUzemi`, `vymera`), `ParcelSearchResult`
- `displayNumber` getter: "1234/5" or "1234"
- Defensive null handling (fields may be missing from API)

### Block 4: ApiClient extension
- **api_client.dart** — add `searchParcelsByPoint(lat, lon, {radius=5})` calling `GET /api/parcels/polygon?lat=&lon=&radius=`
- Returns typed `ParcelSearchResult`

### Block 5: Location service
- **NEW `services/location_service.dart`** — `getCurrentPosition()` returns `LatLng?`
- Checks service enabled → checks/requests permission → gets position (high accuracy, 10s timeout)
- Returns null on denied/disabled (no exceptions)

### Block 6: MapScreen integration
- Add state: `_userPosition`, `_showCadastralOverlay`, `_isLocating`
- Add `_locateUser()` — calls location service, moves map, shows blue dot marker
- Add `_onMapTap(TapPosition, LatLng)` — calls `searchParcelsByPoint`, shows bottom sheet or error snackbar
- Add `_showParcelBottomSheet(Parcel, LatLng)` — modal bottom sheet with parcel number, cadastral area, area m², land type, usage, ownership sheet, coordinates
- Replace single FAB with Column: GPS button (my_location icon) + overlay toggle (layers icon)
- Remove old "Test API" FAB

---

## Verification

1. `cd app && flutter analyze` — zero errors ✓
2. `flutter test` — all tests pass ✓
3. `flutter run` — map shows OSM + colored cadastral overlay at zoom 14+
4. Tap GPS FAB → permission dialog → map centers on user position with blue dot
5. Tap layers FAB → cadastral overlay toggles on/off
6. Tap map with backend running (no CUZK key) → red snackbar "Parcel lookup failed" (no crash)
7. Tap map with backend + valid CUZK key → bottom sheet with parcel details
