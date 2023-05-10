/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// namespaceCmd represents the context command
var namespaceCmd = &cobra.Command{
	Use: "namespace",
	Aliases: []string{
		"ns",
	},
	Short: "Set the KUBECONFIG env variable to a specific namespace in the current context",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !noID {
			konfGoReqID = os.Getenv("MSK_REQID")
			if konfGoReqID == "" {
				log.Fatal().Msg("Request ID not set")
			}
		}

		curKubeConfig := os.Getenv("KUBECONFIG")
		if curKubeConfig == "" {
			log.Fatal().Msg("context not set")
		}

		log.Debug().Msgf("found config: %s", curKubeConfig)

		namespace := ""
		if len(args) == 1 {
			namespace = args[0]
		} else {
			namespaceList, err := konfig.GetNSList(curKubeConfig)
			cobra.CheckErr(err)

			ns, err := shell.LoadList(namespaceList)
			cobra.CheckErr(err)

			namespace = namespaceList[ns]
		}

		kubeConfig, err := konfig.Load(curKubeConfig, homedir)
		cobra.CheckErr(err)

		kubeConfig.Contexts[0].Context.Namespace = namespace

		filePath := path.Dir(curKubeConfig) + "/" + namespace + ".yaml"

		outContext, err := json.Marshal(kubeConfig)
		cobra.CheckErr(err)

		err = os.WriteFile(filePath, outContext, 0640)
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
	rootCmd.AddCommand(namespaceCmd)
}
