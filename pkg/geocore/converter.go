package geo

import (
	"math"

	"github.com/ziprecruiter/h3-go/pkg/h3"
)

// h3EdgeLengthsKm contains the average edge length in km for each H3 resolution (0-15).
// Values sourced from H3 documentation for accurate radius estimation.
var h3EdgeLengthsKm = [16]float64{
	1107.7125910000001, // Res 0
	418.6760055000001,  // Res 1
	158.24465580000002, // Res 2
	59.810857199999994, // Res 3
	22.606379400000002, // Res 4
	8.544408276,        // Res 5
	3.229482772,        // Res 6
	1.220629759,        // Res 7
	0.461354684,        // Res 8
	0.174375668,        // Res 9
	0.065907807,        // Res 10
	0.024910561,        // Res 11
	0.009415526,        // Res 12
	0.003559893,        // Res 13
	0.001348575,        // Res 14
	0.000509713,        // Res 15
}

// haversineDistance calculates the great-circle distance between two points in km.
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLng := (lng2 - lng1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// GeoConverter handles conversions between the package's domain types and the H3 library.
//
// It encapsulates H3-specific encoding/decoding logic so the rest of the geo package can work
// with the domain `LatLng` and `Cell` types without importing H3 directly.
type GeoConverter struct{}

// LatLngToCell converts the coordinate in `cfg` into the H3 cell string for the requested resolution.
//
// The local `LatLng` type is represented in degrees, while the H3 package expects geographic
// coordinates in radians. This method uses `h3.NewLatLng` to perform that conversion internally.
func (gc *GeoConverter) LatLngToCell(cfg Config) (Cell, error) {
	h3Index, err := h3.NewCellFromLatLng(h3.NewLatLng(cfg.Coordinate.Lat, cfg.Coordinate.Lng), cfg.Res)
	if err != nil {
		return "", err
	}
	return Cell(h3Index.String()), nil
}

// CellToLatLng converts an H3 cell ID string back into the package's LatLng coordinates.
//
// The H3 library returns locations in radians, and this method converts them to the package's
// degrees-based `LatLng` representation using `NewLatLngRadians`.
func (gc *GeoConverter) CellToLatLng(cellID Cell) (LatLng, error) {
	h3Index, err := h3.NewCellFromString(string(cellID))
	if err != nil {
		return LatLng{}, err
	}

	latLng, err := h3Index.ToLatLng()
	if err != nil {
		return LatLng{}, err
	}
	return NewLatLngRadians(latLng.Latitude(), latLng.Longitude()), nil
}

// GetCellsInRadius returns the H3 cell IDs within a radius around the coordinate in `cfg`.
//
// The method converts the coordinate in `cfg` into an H3 cell at the requested resolution,
// estimates a k-ring count to cover the physical distance, uses H3's `GridDisk` to get candidates,
// then filters by actual geodesic distance to ensure precise physical radius.
func (gc *GeoConverter) GetCellsInRadius(cfg Config) ([]Cell, error) {
	h3Index, err := h3.NewCellFromLatLng(h3.NewLatLng(cfg.Coordinate.Lat, cfg.Coordinate.Lng), cfg.Res)
	if err != nil {
		return nil, err
	}

	k := distanceToKRing(float64(cfg.Distance), cfg.Res)
	cellIDs, err := h3Index.GridDisk(k)
	if err != nil {
		return nil, err
	}

	var cells []Cell
	for _, cellID := range cellIDs {
		latLng, err := cellID.ToLatLng()
		if err != nil {
			continue
		}
		cellCenter := NewLatLngRadians(latLng.Latitude(), latLng.Longitude())
		dist := haversineDistance(cfg.Coordinate.Lat, cfg.Coordinate.Lng, cellCenter.Lat, cellCenter.Lng)
		if dist <= float64(cfg.Distance) {
			cells = append(cells, Cell(cellID.String()))
		}
	}
	return cells, nil
}

// distanceToKRing estimates the H3 k-ring count needed to cover a physical distance in km.
// It uses the average edge length for the given resolution to calculate k = ceil(distanceKm / edgeKm).
// This ensures GridDisk(k) includes a superset of cells that may be within the radius for filtering.
func distanceToKRing(distanceKm float64, res int) int {
	if distanceKm <= 0 || res < 0 || res > 15 {
		return 0
	}
	edgeKm := h3EdgeLengthsKm[res]
	k := int(math.Ceil(distanceKm / edgeKm))
	return k
}
