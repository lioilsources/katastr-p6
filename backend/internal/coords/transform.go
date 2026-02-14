package coords

import (
	"math"

	"github.com/wroge/wgs84/v2"
)

// WGS84ToSJTSK converts WGS-84 (lat, lon) to S-JTSK (x, y).
// Returns POSITIVE values as expected by the CUZK API.
func WGS84ToSJTSK(lat, lon float64) (x, y float64) {
	transform := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(5514))
	east, north, _ := transform(lon, lat, 0)
	// S-JTSK natively uses negative coordinates; CUZK API expects positive.
	return math.Abs(north), math.Abs(east)
}

// SJTSKToWGS84 converts S-JTSK (x, y) to WGS-84 (lat, lon).
// Input values are POSITIVE (as returned by CUZK API).
func SJTSKToWGS84(x, y float64) (lat, lon float64) {
	transform := wgs84.Transform(wgs84.EPSG(5514), wgs84.EPSG(4326))
	// Convert to negative for S-JTSK projection convention.
	outLon, outLat, _ := transform(-y, -x, 0)
	return outLat, outLon
}
