package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "geocli",
	Short: "go-geo command line tool",
	Long:  "A small CLI wrapper for the go-geo package.",
}

func Execute() error {
	return rootCmd.Execute()
}

func addCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
