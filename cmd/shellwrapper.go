/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// shellwrapperCmd represents the shellwrapper command
var shellwrapperCmd = &cobra.Command{
	Use:   "shellwrapper",
	Short: "Wrap a shell function around msk-bin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var wrapper string
		var zsh = `
msk() {
	export MSK_REQID="$(msk-bin genid)"
  msk-bin $@
	res=$(cat /tmp/${MSK_REQID})
  # only change $KUBECONFIG if instructed by konf-go
  if [[ $res == "KUBECONFIGTOUSE:"* ]]
  then
    # this basically takes the line and cuts out the KUBECONFIGTOUSE Part
    export KUBECONFIG="${res#*KUBECONFIGTOUSE:}"
  else
    # this makes --help work
    echo "${res}"
  fi
	rm -f /tmp/${MSK_REQID}
	unset MSK_REQID
}
`

		var bash = `
konf() {
  res=$(msk-bin $@)
  # only change $KUBECONFIG if instructed by konf-go
  if [[ $res == "KUBECONFIGTOUSE:"* ]]
  then
    # this basically takes the line and cuts out the KUBECONFIGTOUSE Part
    export KUBECONFIG="${res#*KUBECONFIGTOUSE:}"
  else
    # this makes --help work
    echo "${res}"
  fi
}
`

		switch args[0] {
		case "zsh":
			wrapper = zsh
		case "bash":
			wrapper = bash
		default:
			return fmt.Errorf("multiShellKonfig currently does not support %s", args[0])
		}

		fmt.Println(wrapper)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(shellwrapperCmd)
}
