package geo

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidLatitude   = errors.New("latitude must be between -90 and 90 degrees")
	ErrInvalidLongitude  = errors.New("longitude must be between -180 and 180 degrees")
	ErrInvalidResolution = errors.New("H3 resolution must be between 0 and 15")
)

// Unit type for explicit measurement systems
type Unit int

const (
	UnitKilometers Unit = iota
	UnitMiles
)

// LatLng represents a geographical point.
// Standardized for use in search filters.
type LatLng struct {
	Lat, Lng float64
}

// NewLatLngDegrees is a constructor that ensures consistent handling of latitude and longitude values in degrees.
func NewLatLngDegrees(lat, lng float64) LatLng {
	return LatLng{Lat: lat, Lng: lng}
}

// NewLatLngRadians is a constructor that converts radians to degrees before creating a LatLng instance.
func NewLatLngRadians(latRad, lngRad float64) LatLng {
	return LatLng{
		Lat: latRad * (180.0 / math.Pi),
		Lng: lngRad * (180.0 / math.Pi),
	}
}

// Cell is a domain-specific type for the H3 Index string.
// Chosen over uint64 to ensure seamless integration with gomobile
// and to avoid 64-bit integer precision issues in mobile runtimes.
type Cell string

// CellResult represents a single cell and its resolution.
type CellResult struct {
	ID  Cell
	Res int
}

// CellList is a collection of cell results, useful for bulk processing.
type CellList struct {
	IDs []CellResult
}

// Distance handles internal conversion logic
type Distance float64

// NewDistance is the constructor that handles unit conversion logic.
// It ensures the internal engine always receives metric values.
func NewDistance(value float64, unit Unit) Distance {
	if unit == UnitMiles {
		return Distance(value * 1.60934)
	}
	return Distance(value)
}

// Config provides a structured way to pass parameters to the CLI or App.
type Config struct {
	Coordinate LatLng
	Res        int
	Distance   Distance
	Unit       Unit // Tracks user preference (Metric vs Imperial)
}

// ConfigBuilder facilitates the step-by-step construction of a Config.
type ConfigBuilder struct {
	config Config
}

// NewConfigBuilder initializes a builder with default values.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: Config{
			Res:  9,              // Defaulting to Res 9 for a good balance of precision and performance
			Unit: UnitKilometers, // Defaulting to Metric
		},
	}
}

// WithCoordinate sets the central LatLng for the search.
func (b *ConfigBuilder) WithCoordinate(lat, lng float64) *ConfigBuilder {
	b.config.Coordinate = NewLatLngDegrees(lat, lng)
	return b
}

// WithResolution sets the H3 grid resolution (Default is 9).
func (b *ConfigBuilder) WithResolution(res int) *ConfigBuilder {
	b.config.Res = res
	return b
}

// WithRadius sets the search radius and handles unit conversion internally.
func (b *ConfigBuilder) WithRadius(val float64, unit Unit) *ConfigBuilder {
	b.config.Distance = NewDistance(val, unit)
	b.config.Unit = unit
	return b
}

// Build validates the configuration and returns the final Config struct.
func (b *ConfigBuilder) Build() (Config, error) {
	// Validate Latitude
	if b.config.Coordinate.Lat < -90 || b.config.Coordinate.Lat > 90 {
		return Config{}, fmt.Errorf("%w: received %f", ErrInvalidLatitude, b.config.Coordinate.Lat)
	}

	// Validate Longitude
	if b.config.Coordinate.Lng < -180 || b.config.Coordinate.Lng > 180 {
		return Config{}, fmt.Errorf("%w: received %f", ErrInvalidLongitude, b.config.Coordinate.Lng)
	}

	// Validate H3 Resolution
	if b.config.Res < 0 || b.config.Res > 15 {
		return Config{}, fmt.Errorf("%w: received %d", ErrInvalidResolution, b.config.Res)
	}

	return b.config, nil
}
