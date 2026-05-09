package main

import (
	"fmt"

	geo "github.com/expoR93/go-geo/pkg/geocore"
	"github.com/spf13/cobra"
)

var radiusCmd = &cobra.Command{
	Use:   "radius",
	Short: "Get H3 cell IDs within a specified radius of a coordinate",
	RunE: func(cmd *cobra.Command, args []string) error {
		u := geo.UnitKilometers
		if unit == "m" {
			u = geo.UnitMiles
		}
		cfg, err := geo.NewConfigBuilder().WithCoordinate(lat, lng).WithRadius(radius, u).Build()
		if err != nil {
			return err
		}

		converter := geo.GeoConverter{}
		cells, err := converter.GetCellsInRadius(cfg)
		if err != nil {
			return err
		}

		for _, cell := range cells {
			fmt.Println(cell)
		}
		return nil
	},
}

var (
	radius float64
	unit   string
)

func init() {
	radiusCmd.Flags().Float64Var(&lat, "lat", 0, "Latitude in degrees")
	radiusCmd.Flags().Float64Var(&lng, "lng", 0, "Longitude in degrees")
	radiusCmd.Flags().Float64Var(&radius, "radius", 0, "Radius distance")
	radiusCmd.Flags().StringVar(&unit, "unit", "m", "Unit of distance (m for meters, km for kilometers)")
	_ = radiusCmd.MarkFlagRequired("lat")
	_ = radiusCmd.MarkFlagRequired("lng")
	_ = radiusCmd.MarkFlagRequired("radius")

	addCommand(radiusCmd)
}
