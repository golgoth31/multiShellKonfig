/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "A brief description of your command",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !noID {
			konfGoReqID = os.Getenv("MSK_REQID")
			if konfGoReqID == "" {
				log.Debug().Msg("Request ID not set")

				cobra.CheckErr(errNoReqID)
			}
		}

		if len(args) == 1 {
			// get single context

			os.Exit(0)
		}

		contextList := []shell.ContextDef{}

		// get all available contexts
		for _, unitKonfig := range cfgData.Konfigs {
			log.Debug().Msgf("found config: %s", unitKonfig.Path)

			kubeConfig, err := konfig.Load(unitKonfig.Path, homedir)
			if err != nil {
				return err
			}

			for _, context := range kubeConfig.Contexts {
				log.Debug().Msgf("found context '%s@%s'", unitKonfig.Path, context.Name)

				contextList = append(contextList, shell.ContextDef{
					Name:     context.Name,
					FileID:   unitKonfig.ID,
					FilePath: unitKonfig.Path,
				})
			}
		}

		log.Debug().Msgf("context list: %v", contextList)

		context, err := shell.LoadPterm(contextList)
		if err != nil {
			return err
		}

		kubeConfig, err := konfig.Load(context.FilePath, homedir)
		if err != nil {
			return err
		}

		filePath, err := konfig.Generate(&context, kubeConfig, cfgContexts)
		if err != nil {
			return err
		}

		log.Debug().Msgf("KUBECONFIGTOUSE:" + filePath)

		if !noID {
			if err = os.WriteFile(
				fmt.Sprintf("/tmp/%s", konfGoReqID),
				[]byte("KUBECONFIGTOUSE:"+filePath),
				0666,
			); err != nil {
				return err
			}
		} else {
			log.Info().Msgf("KUBECONFIGTOUSE:" + filePath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
