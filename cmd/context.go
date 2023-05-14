/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use: "context",
	Aliases: []string{
		"ctx",
	},
	Short: "Set the KUBECONFIG env variable to a specific context",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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

		contextList := shell.ShellContextList{}
		contextListString := []string{}

		// get all available contexts
		for _, unitKonfig := range cfgData.Konfigs {
			log.Debug().Msgf("found config: %s", unitKonfig.Path)

			kubeConfig, err := konfig.Load(unitKonfig.Path, homedir)
			cobra.CheckErr(err)

			for _, context := range kubeConfig.Contexts {
				log.Debug().Msgf("found context '%s@%s'", unitKonfig.Path, context.Name)

				contextList = append(
					contextList,
					shell.ContextDef{
						Name:     context.Name,
						FileID:   unitKonfig.ID,
						FilePath: unitKonfig.Path,
					},
				)
			}
		}

		// Sort contextList by context name
		sort.Stable(contextList)

		for _, v := range contextList {
			contextListString = append(
				contextListString,
				fmt.Sprintf(
					"%s (file: %s)",
					v.Name,
					v.FilePath,
				),
			)
		}

		log.Debug().Msgf("context list: %v", contextList)

		contextID, err := shell.LoadList(contextListString)
		cobra.CheckErr(err)

		kubeConfig, err := konfig.Load(contextList[contextID].FilePath, homedir)
		cobra.CheckErr(err)

		filePath, fileData, err := konfig.Generate(&contextList[contextID], kubeConfig, cfgContextsPath)
		cobra.CheckErr(err)

		err = konfig.SaveContextFile(filePath, fileData, false)
		cobra.CheckErr(err)

		log.Debug().Msgf("KUBECONFIGTOUSE:" + filePath)

		if !noID {
			err := os.WriteFile(
				fmt.Sprintf("/tmp/%s", konfGoReqID),
				[]byte("KUBECONFIGTOUSE:"+filePath),
				0666,
			)
			cobra.CheckErr(err)
		} else {
			log.Info().Msgf("KUBECONFIGTOUSE:" + filePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
