/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// genidCmd represents the genid command
var genidCmd = &cobra.Command{
	Use:   "genid",
	Short: "Generate a uniq id used by the wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.CreateTemp("/tmp/", "msk-")
		if err != nil {
			log.Error().Err(err).Msg("")
		}

		fmt.Println(file.Name())
	},
}

func init() {
	rootCmd.AddCommand(genidCmd)
}
