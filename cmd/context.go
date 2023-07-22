/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golgoth31/multiShellKonfig/internal/context"
	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var (
	contextCmd = &cobra.Command{
		Use: "context",
		Aliases: []string{
			"ctx",
		},
		Short: "Set the KUBECONFIG env variable to a specific context",
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ctxObj, err := context.New(cfgData.Konfigs, homedir)
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveError
			}

			contextListString, err := ctxObj.GetContextList()
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveError
			}

			return contextListString, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctxObj, err := context.New(cfgData.Konfigs, homedir)
			cobra.CheckErr(err)

			ctxObj.MskReqID = os.Getenv("MSK_REQID")
			if ctxObj.MskReqID == "" {
				log.Fatal().Msg(errNoReqID.Error())
			}

			contextListString, err := ctxObj.GetContextList()
			cobra.CheckErr(err)

			log.Debug().Msgf("%d context(s) found", len(contextListString))

			contextID := ""

			if len(args) == 1 {
				contextID = args[0]
			} else {
				switch len(contextListString) {
				case 0:
					log.Info().Msg("No context found")
					os.Exit(0)
				case 1:
					log.Info().Msg("Only one context found, using it")
					contextID = contextListString[0]
				default:
					var errContextID error

					contextID, errContextID = shell.LoadList("context", contextListString)
					cobra.CheckErr(errContextID)
				}
			}

			curKonfig := konfig.Konfig{}

			// extract context and file path from selection returned
			contextSplit := strings.Split(contextID, "@")

			if len(contextSplit) == 1 {
				cobra.CheckErr(errors.New("context name badly formated"))
			}

			contextName := contextSplit[0]
			log.Debug().Msgf("selected context: %s", contextName)

			// contextFilePath := strings.Trim(contextSplit[1], "()")
			contextFilePath := contextSplit[1]
			log.Debug().Msgf("selected file: %s", contextFilePath)

			// select the right file
			for _, konfigUnit := range ctxObj.KonfigList {
				log.Debug().Msgf("%q", contextName)
				log.Debug().Msgf("%q", contextFilePath)
				log.Debug().Msgf("%s", konfigUnit.FilePath)

				// if konfigUnit.FilePath == contextFilePath {
				if konfigUnit.FilePath == contextFilePath {
					curKonfig = *konfigUnit

					break
				}
			}

			if curKonfig.FileID == "" {
				cobra.CheckErr(errors.New("Konfig not found"))
			}

			filePath, fileData, err := curKonfig.Generate(contextName, cfgContextsPath)
			cobra.CheckErr(err)

			err = konfig.SaveContextFile(filePath, fileData, false)
			cobra.CheckErr(err)

			log.Debug().Msgf("KUBECONFIGTOUSE:" + filePath)

			if !noID {
				err := os.WriteFile(
					fmt.Sprintf("/tmp/%s", ctxObj.MskReqID),
					[]byte("KUBECONFIGTOUSE:"+filePath),
					0666,
				)
				cobra.CheckErr(err)
			} else {
				log.Info().Msgf("KUBECONFIGTOUSE:" + filePath)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(contextCmd)
}
