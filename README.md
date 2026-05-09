# go-geo

A Go library for geospatial operations using the H3 hexagonal grid system. Provides coordinate-to-cell conversion, cell-to-coordinate lookup, and radius-based cell queries. Includes a command-line tool and supports mobile binding with gomobile.

## Features

- Convert latitude/longitude to H3 cell IDs
- Convert H3 cell IDs back to coordinates
- Find all cells within a radius of a point
- Pure Go implementation (no CGO dependencies)
- Command-line interface with Cobra
- Compatible with gomobile for mobile apps

## Installation

### Library

```bash
go get github.com/expoR93/go-geo/pkg/geocore
```

### CLI Tool

```bash
go install github.com/expoR93/go-geo/cmd/geocli
```

Or build from source:

```bash
git clone https://github.com/expoR93/go-geo.git
cd go-geo
go build ./cmd/geocli
```

## Usage

### Library

```go
import "github.com/expoR93/go-geo/pkg/geocore"

func main() {
    // Create a config
    cfg, _ := geocore.NewConfigBuilder().
        WithCoordinate(37.7749, -122.4194).
        WithResolution(9).
        WithRadius(5.0, geocore.UnitKilometers).
        Build()

    converter := geocore.GeoConverter{}

    // Get H3 cell ID
    cell, _ := converter.LatLngToCell(cfg)
    fmt.Println(cell) // e.g., "8928308280fffff"

    // Get coordinates from cell
    latLng, _ := converter.CellToLatLng(cell)
    fmt.Printf("Lat: %f, Lng: %f\n", latLng.Lat, latLng.Lng)

    // Get cells in radius
    cells, _ := converter.GetCellsInRadius(cfg)
    for _, c := range cells {
        fmt.Println(c)
    }
}
```

### CLI Tool

```bash
# Convert coordinate to H3 cell
geocli cell --lat 37.7749 --lng -122.4194 --res 9

# Convert H3 cell to coordinate
geocli coord 8928308280fffff

# Get cells within radius
geocli radius --lat 37.7749 --lng -122.4194 --radius 5 --unit km
```

## Mobile Binding

Build for Android/iOS using gomobile:

```bash
# Install gomobile
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init

# Build for Android
gomobile bind -target=android github.com/expoR93/go-geo/pkg/geocore

# Build for iOS
gomobile bind -target=ios github.com/expoR93/go-geo/pkg/geocore
```

This generates platform-specific bindings you can integrate into mobile apps.

## API Reference

### Types

- `LatLng`: Represents a geographic coordinate
- `Cell`: H3 cell ID as a string
- `Config`: Configuration for operations
- `GeoConverter`: Main converter struct

### Methods

- `NewConfigBuilder()`: Create a config builder
- `LatLngToCell(cfg Config)`: Convert coordinate to cell
- `CellToLatLng(cellID Cell)`: Convert cell to coordinate
- `GetCellsInRadius(cfg Config)`: Get cells in radius

## Testing

Run tests:

```bash
go test ./...
```

## License

See LICENSE file.