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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// namespaceCmd represents the context command
var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "A brief description of your command",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		konfGoReqID = os.Getenv("MSK_REQID")
		if konfGoReqID == "" {
			log.Fatal().Msg("Request ID not set")
		}

		curKubeConfig := os.Getenv("KUBECONFIG")
		if konfGoReqID == "" {
			log.Fatal().Msg("context not set")
		}

		namespace := "toto"

		if len(args) == 1 {
			// get single context

			os.Exit(0)
		}

		log.Debug().Msgf("found config: %s", curKubeConfig)
		kubeConfig, err := konfig.Load(curKubeConfig, homedir)
		cobra.CheckErr(err)

		kubeConfig.Contexts[0].Context.Namespace = namespace

		filePath := path.Dir(curKubeConfig) + "/" + namespace + ".yaml"

		outContext, err := json.Marshal(kubeConfig)
		cobra.CheckErr(err)

		err = os.WriteFile(filePath, outContext, 0640)
		cobra.CheckErr(err)

		fmt.Println("KUBECONFIGTOUSE:" + filePath)
		err = os.WriteFile(
			fmt.Sprintf("/tmp/%s", konfGoReqID),
			[]byte("KUBECONFIGTOUSE:"+filePath),
			0666,
		)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
}