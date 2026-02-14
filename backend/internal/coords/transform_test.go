package coords

import (
	"math"
	"testing"
)

func TestWGS84ToSJTSK(t *testing.T) {
	// Old Town Square, Prague: WGS-84 (50.088, 14.421)
	x, y := WGS84ToSJTSK(50.088, 14.421)

	// Verify values are positive and in expected range for Prague
	if x <= 1000000 || x >= 1100000 {
		t.Errorf("X out of Prague range: %.0f", x)
	}
	if y <= 700000 || y >= 800000 {
		t.Errorf("Y out of Prague range: %.0f", y)
	}
}

func TestSJTSKToWGS84(t *testing.T) {
	// Convert a known S-JTSK point back to WGS-84
	// Use values from forward transform for consistency
	x, y := WGS84ToSJTSK(50.088, 14.421)
	lat, lon := SJTSKToWGS84(x, y)

	if math.Abs(lat-50.088) > 0.001 || math.Abs(lon-14.421) > 0.001 {
		t.Errorf("SJTSKToWGS84(%.0f, %.0f) = (%.4f, %.4f), expected ~(50.088, 14.421)", x, y, lat, lon)
	}
}

func TestRoundTrip(t *testing.T) {
	points := [][2]float64{
		{50.100, 14.390},  // Dejvice, Prague 6
		{50.088, 14.421},  // Old Town Square
		{50.075, 14.437},  // Vysehrad
		{50.108, 14.340},  // Vokovice, Prague 6
	}

	for _, p := range points {
		x, y := WGS84ToSJTSK(p[0], p[1])
		lat, lon := SJTSKToWGS84(x, y)

		if math.Abs(lat-p[0]) > 0.0001 || math.Abs(lon-p[1]) > 0.0001 {
			t.Errorf("round-trip (%.4f, %.4f) -> (%.0f, %.0f) -> (%.4f, %.4f)",
				p[0], p[1], x, y, lat, lon)
		}
	}
}

func TestPositiveValues(t *testing.T) {
	x, y := WGS84ToSJTSK(50.100, 14.390)
	if x <= 0 || y <= 0 {
		t.Errorf("expected positive values, got (%.0f, %.0f)", x, y)
	}
}
