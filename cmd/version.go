/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/golgoth31/multiShellKonfig/internal/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show msk version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"version: %s\nbuildDate: %s\nBuiltBy: %s\n",
			config.Version,
			config.Date,
			config.BuiltBy,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
