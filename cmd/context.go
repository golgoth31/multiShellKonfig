/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

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
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		outputs := []string{"hello zzzz", "moto"}
		return outputs, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !noID {
			konfGoReqID = os.Getenv("MSK_REQID")
			if konfGoReqID == "" {
				log.Debug().Msg("Request ID not set")

				cobra.CheckErr(errNoReqID)
			}
		}

		contextList := shell.ShellContextList{}
		contextListString := []string{}
		konfigList := []*konfig.Konfig{}

		// get all available contexts
		for _, unitKonfig := range cfgData.Konfigs {
			log.Debug().Msgf("found config: %s", unitKonfig.Path)

			kubeConfig, err := konfig.Load(unitKonfig.Path, homedir)
			cobra.CheckErr(err)

			curKonfig := konfig.Konfig{
				FileID:   unitKonfig.ID,
				FilePath: unitKonfig.Path,
				Content:  kubeConfig,
			}

			konfigList = append(konfigList, &curKonfig)

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

		// Generate list of context for select
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

		log.Debug().Msgf("%d context(s) found", len(contextListString))
		log.Debug().Msgf("context list: %v", contextList)

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
			default:
				var err error

				contextID, err = shell.LoadList("context", contextListString)
				cobra.CheckErr(err)
			}
		}

		curKonfig := konfig.Konfig{}

		// extract context and file path from selection returned
		contextSplit := strings.Split(contextID, " (file: ")

		if len(contextSplit) == 1 {
			cobra.CheckErr(errors.New("context name badly formated"))
		}

		contextName := contextSplit[0]
		log.Debug().Msgf("selected context: %s", contextName)

		contextFileID := strings.Trim(contextSplit[1], "()")
		log.Debug().Msgf("selected file: %s", contextFileID)

		// select the right file
		for _, konfigUnit := range konfigList {
			log.Debug().Msgf("%s", contextFileID)

			if konfigUnit.FilePath == contextFileID {
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
