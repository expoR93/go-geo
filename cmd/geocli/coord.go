package main

import (
	"fmt"

	geo "github.com/expoR93/go-geo/pkg/geocore"
	"github.com/spf13/cobra"
)

var coordCmd = &cobra.Command{
	Use:   "coord <h3-id>",
	Short: "Convert an H3 cell ID to a coordinate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		converter := geo.GeoConverter{}

		latLng, err := converter.CellToLatLng(geo.Cell(args[0]))
		if err != nil {
			return err
		}

		fmt.Printf("Latitude: %f, Longitude: %f\n", latLng.Lat, latLng.Lng)
		return nil
	},
}

func init() {
	addCommand(coordCmd)
}
