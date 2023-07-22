/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"path"

	"github.com/golgoth31/multiShellKonfig/internal/namespace"
	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// namespaceCmd represents the context command
var (
	namespaceCmd = &cobra.Command{
		Use: "namespace",
		Aliases: []string{
			"ns",
		},
		Short: "Set the KUBECONFIG env variable to a specific namespace in the current context",
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			nsObj, err := namespace.New(os.Getenv("KUBECONFIG"))
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveError
			}

			output, err := nsObj.GetNsList()
			if err != nil {
				return []string{}, cobra.ShellCompDirectiveError
			}

			if len(args) != 0 {
				for _, ns := range output {
					if ns == args[0] {
						output = []string{
							ns,
						}
					}
				}
			}
			return output, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			nsObj, err := namespace.New(os.Getenv("KUBECONFIG"))
			cobra.CheckErr(err)

			nsObj.MskReqID = os.Getenv("MSK_REQID")
			if nsObj.MskReqID == "" {
				log.Fatal().Msg(errNoReqID.Error())
			}

			localNamespace := ""
			if len(args) == 1 {
				localNamespace = args[0]
			} else {
				namespaceList, errNamespaceList := nsObj.GetNsList()
				cobra.CheckErr(errNamespaceList)

				ns, errNs := shell.LoadList("namespace", namespaceList)
				cobra.CheckErr(errNs)

				localNamespace = ns
			}

			kubeConfig, err := konfig.Load(nsObj.CurKubeConfig, homedir)
			cobra.CheckErr(err)

			kubeConfig.Contexts[0].Context.Namespace = localNamespace

			filePath := path.Dir(nsObj.CurKubeConfig) + "/" + localNamespace + ".yaml"

			fileData, err := json.Marshal(kubeConfig)
			cobra.CheckErr(err)

			err = konfig.SaveContextFile(filePath, fileData)
			cobra.CheckErr(err)

			log.Debug().Msgf("KUBECONFIGTOUSE:" + filePath)

			err = os.WriteFile(
				nsObj.MskReqID,
				[]byte("KUBECONFIGTOUSE:"+filePath),
				filePerm,
			)
			cobra.CheckErr(err)
		},
	}
)

func init() {
	rootCmd.AddCommand(namespaceCmd)
}
