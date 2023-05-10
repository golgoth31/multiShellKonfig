/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/golgoth31/multiShellKonfig/internal/config"
	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new kubeconfig file to the list of known files",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug().Msgf("importing file: %s", args[0])

			// read file to ensure it is available

			// make path relative to home dir if possible
			pathOrig := args[0]
			path := ""

			if strings.HasPrefix(pathOrig, homedir) {
				path = strings.TrimPrefix(pathOrig, homedir)
				path = "~" + path
			} else {
				path = pathOrig
			}

			// check if already exists
			add := true
			for _, konfig := range cfgData.Konfigs {
				if path == konfig.Path {
					add = false

					break
				}
			}

			if add {
				_, err := konfig.Load(path, homedir)
				cobra.CheckErr(err)

				localConf := config.Konfig{
					Path: path,
					ID:   xid.New().String(),
				}
				cfgData.Konfigs = append(cfgData.Konfigs, localConf)

				cfgDataByte, err := yaml.Marshal(cfgData)
				cobra.CheckErr(err)

				log.Debug().Msgf("%s", cfgDataByte)

				err = os.WriteFile(cfgFile, cfgDataByte, 0640)
				cobra.CheckErr(err)
			} else {
				log.Info().Msg("path already exists in config")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
