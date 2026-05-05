package geo

import (
	"testing"

	"github.com/ziprecruiter/h3-go/pkg/h3"
)

func TestGeoConverter_LatLngToCell(t *testing.T) {
	cfg := Config{
		Coordinate: NewLatLngDegrees(37.775938728915946, -122.41795063018799),
		Res:        9,
	}
	gc := &GeoConverter{}

	got, err := gc.LatLngToCell(cfg)
	if err != nil {
		t.Fatalf("LatLngToCell returned error: %v", err)
	}

	h3Index, err := h3.NewCellFromLatLng(h3.NewLatLng(cfg.Coordinate.Lat, cfg.Coordinate.Lng), cfg.Res)
	if err != nil {
		t.Fatalf("failed to create expected H3 index: %v", err)
	}
	want := Cell(h3Index.String())
	if got != want {
		t.Fatalf("LatLngToCell = %s, want %s", got, want)
	}
}

func TestGeoConverter_CellToLatLng(t *testing.T) {
	cfg := Config{
		Coordinate: NewLatLngDegrees(37.775938728915946, -122.41795063018799),
		Res:        9,
	}
	gc := &GeoConverter{}

	h3Index, err := h3.NewCellFromLatLng(h3.NewLatLng(cfg.Coordinate.Lat, cfg.Coordinate.Lng), cfg.Res)
	if err != nil {
		t.Fatalf("failed to create H3 index for test: %v", err)
	}
	cellID := Cell(h3Index.String())
	got, err := gc.CellToLatLng(cellID)
	if err != nil {
		t.Fatalf("CellToLatLng returned error: %v", err)
	}

	h3Center, err := h3.NewCellFromString(string(cellID))
	if err != nil {
		t.Fatalf("failed to parse expected cell string: %v", err)
	}

	latLng, err := h3Center.ToLatLng()
	if err != nil {
		t.Fatalf("failed to get expected cell center: %v", err)
	}

	want := NewLatLngRadians(latLng.Latitude(), latLng.Longitude())
	assertAlmostEqual(t, want.Lat, got.Lat, "Lat")
	assertAlmostEqual(t, want.Lng, got.Lng, "Lng")
}

func TestGeoConverter_GetCellsInRadius(t *testing.T) {
	cfg := Config{
		Coordinate: NewLatLngDegrees(37.775938728915946, -122.41795063018799),
		Res:        9,
		Distance:   NewDistance(0.1, UnitKilometers),
	}
	gc := &GeoConverter{}

	cells, err := gc.GetCellsInRadius(cfg)
	if err != nil {
		t.Fatalf("GetCellsInRadius returned error: %v", err)
	}
	if len(cells) == 0 {
		t.Fatal("GetCellsInRadius returned no cells")
	}

	for _, cell := range cells {
		cellLatLng, err := gc.CellToLatLng(cell)
		if err != nil {
			t.Fatalf("CellToLatLng returned error for cell %s: %v", cell, err)
		}

		distance := haversineDistance(cfg.Coordinate.Lat, cfg.Coordinate.Lng, cellLatLng.Lat, cellLatLng.Lng)
		if distance > float64(cfg.Distance)+tolerance {
			t.Fatalf("cell %s is outside radius: distance %f km", cell, distance)
		}
	}
}

func TestDistanceToKRing(t *testing.T) {
	testCases := []struct {
		name       string
		distance   float64
		resolution int
		want       int
	}{
		{"Zero distance", 0, 9, 0},
		{"Small radius at res 9", 0.1, 9, 1},
		{"Just above one edge at res 9", 0.191, 9, 2},
		{"Invalid resolution", 0.1, 99, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := distanceToKRing(tc.distance, tc.resolution)
			if got != tc.want {
				t.Fatalf("distanceToKRing(%f, %d) = %d, want %d", tc.distance, tc.resolution, got, tc.want)
			}
		})
	}
}
