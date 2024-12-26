package base

import (
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func Commands(app *console.Console) console.Commands {
	return func() *cobra.Command {
		rootCmd := &cobra.Command{
			DisableFlagsInUseLine: true,
			SilenceErrors:         true,
			SilenceUsage:          true,
		}

		// chat
		rootCmd.AddCommand(chatCommand(app))
		// exit
		rootCmd.AddCommand(exitCommand(app))
		// ant
		rootCmd.AddCommand(antCommand(app))
		// use
		rootCmd.AddCommand(useCommand(app))
		// script
		rootCmd.AddCommand(scriptCommand(app))

		rootCmd.InitDefaultHelpCmd()
		rootCmd.CompletionOptions.DisableDefaultCmd = true
		return rootCmd
	}
}
