package geo

type CellIndexer interface {
	// LatLngToCell returns the Cell ID string for a given coordinate​
	LatLngToCell(latLng LatLng, res int) (Cell, error)

	// CellToLatLng returns the center coordinates of a Cell ID
	CellToLatLng(cellID Cell) (LatLng, error)

	// GetCellsInRadius returns all Cell IDs within a certain Km/Miles radius
	GetCellsInRadius(latlng LatLng, radius Distance, res int) ([]Cell, error)
}
