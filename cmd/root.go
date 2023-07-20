/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path"

	"github.com/golgoth31/multiShellKonfig/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "msk",
	Short: "A brief description of your application",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolVar(&noID, "no-id", false, "disable request id usage")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Find home directory.
	var err error

	homedir, err = os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".msk" (without extension).
	cfgDir = homedir + "/.kube/msk"
	cfgFile = cfgDir + "/config.yaml"
	cfgContextsPath = cfgDir + "/contexts"

	if _, errOsStat := os.Stat(cfgContextsPath); err != nil {
		if errors.Is(errOsStat, fs.ErrNotExist) {
			err = os.MkdirAll(cfgContextsPath, 0755)
			cobra.CheckErr(err)
		}
	}

	cfgDataByte, err := os.ReadFile(cfgFile)
	if err != nil {
		if _, errOsStat := os.Stat(path.Dir(cfgFile)); err != nil {
			if errors.Is(errOsStat, fs.ErrNotExist) {
				err = os.MkdirAll(path.Dir(cfgFile), 0755)
				cobra.CheckErr(err)
			}
		}

		cfgDataByte, err = yaml.Marshal(config.DefaultConfig)
		cobra.CheckErr(err)

		err = os.WriteFile(cfgFile, cfgDataByte, 0640)
		cobra.CheckErr(err)
	}

	err = yaml.Unmarshal(cfgDataByte, &cfgData)
	cobra.CheckErr(err)
}
