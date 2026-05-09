package main

import (
	"fmt"

	geo "github.com/expoR93/go-geo/pkg/geocore"
	"github.com/spf13/cobra"
)

var (
	lat float64
	lng float64
	res int
)

var cellCmd = &cobra.Command{
	Use:   "cell",
	Short: "Convert a coordinate to an H3 cell ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := geo.NewConfigBuilder().WithCoordinate(lat, lng).WithResolution(res).Build()
		if err != nil {
			return err
		}

		converter := geo.GeoConverter{}
		cell, err := converter.LatLngToCell(cfg)
		if err != nil {
			return err
		}

		fmt.Println(cell)
		return nil
	},
}

func init() {
	cellCmd.Flags().Float64Var(&lat, "lat", 0, "Latitude in degrees")
	cellCmd.Flags().Float64Var(&lng, "lng", 0, "Longitude in degrees")
	cellCmd.Flags().IntVar(&res, "res", 9, "H3 resolution")
	_ = cellCmd.MarkFlagRequired("lat")
	_ = cellCmd.MarkFlagRequired("lng")

	addCommand(cellCmd)
}
