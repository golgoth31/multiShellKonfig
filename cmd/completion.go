/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh]",
	Short: "Generate completion script",
	Long: fmt.Sprintf(`To load completions:

Bash:

  $ source <(%[1]s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
  # macOS:
  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

  # You will need to start a new shell for this setup to take effect.
`, rootCmd.Name()),
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {

		case "zsh":
			// This allows for also using 'source <(konf completion zsh)' with zsh, similar to bash.
			// Basically it just adds the compdef command so it can be run. Taken from kubectl, who
			// do a similar thing
			zshHeader := "#compdef _msk msk\ncompdef _msk msk\n"

			// So per default cobra makes use of the words[] array that zsh provides to you in completion funcs.
			// Words is an array that contains all words that have been typed by the user before hitting tab
			// Now cobra takes words[1] which is equal to the name of the comand and uses this to call completion on it
			// However in our case this does not work as words[1] points to 'konf' which is the wrapper and not the binary
			// In order to solve this we have to ensure that words[1] equates to konf-go, which is the binary.
			// Currently I have found, the fastest way to do this is by inserting a line to overwrite words[1]. This is
			// because the words[1] reference is used throughout the script and I would not want to replace all of it
			var b bytes.Buffer
			err := rootCmd.GenZshCompletion(&b)
			if err != nil {
				return err
			}
			anchor := "local -a completions" // this is basically a line early in the original script that we are going to cling onto
			genZsh := strings.Replace(b.String(), anchor, anchor+"\n    words[1]=\"msk-bin\"", 1)

			if _, err := os.Stdout.WriteString(zshHeader + genZsh); err != nil {
				return err
			}

		case "bash":
			var b bytes.Buffer
			err := rootCmd.GenBashCompletionV2(&b, true)
			if err != nil {
				return err
			}
			anchor := "local requestComp lastParam lastChar args"
			genBash := strings.Replace(b.String(), anchor, anchor+"\n    words[0]=\"msk-bin\"", 1)

			if _, err := os.Stdout.WriteString(genBash); err != nil {
				return err
			}

		default:
			return fmt.Errorf("msk currently does not support autocompletions for %s", args[0])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
