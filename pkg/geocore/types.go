package geo

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
