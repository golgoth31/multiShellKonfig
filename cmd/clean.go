/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean old cluster files",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Cleaning context(s)")

		// Clean all contexts from files not in config
		dirList, err := os.ReadDir(cfgContextsPath)
		cobra.CheckErr(err)

		kubeConfigFileNumber := 0
		for _, v := range dirList {
			log.Debug().Msg(v.Name())

			toDelete := true

			if !cleanAll {
				for _, konfig := range cfgData.Konfigs {
					if v.Name() == konfig.ID {
						toDelete = false

						break
					}
				}
			}

			if toDelete {
				log.Debug().Msgf("deleting %s", v.Name())

				err := os.RemoveAll(cfgContextsPath + "/" + v.Name())
				cobra.CheckErr(err)

				kubeConfigFileNumber++
			}
		}

		if kubeConfigFileNumber == 0 {
			log.Info().Msg("Nothing to clean")
		} else {
			log.Info().Msgf("All contexts cleaned for %d kubeConfig file", kubeConfigFileNumber)
		}
		// list all context, remove all other
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Clean all known contexts")
}
