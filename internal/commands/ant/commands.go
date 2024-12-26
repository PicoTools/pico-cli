package ant

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func Commands(app *console.Console) console.Commands {
	return func() *cobra.Command {
		cmd := &cobra.Command{
			DisableFlagsInUseLine: true,
			SilenceErrors:         true,
			SilenceUsage:          true,
		}

		cmd.AddGroup(
			&cobra.Group{ID: constants.AliasGroupId, Title: constants.AliasGroupId},
			&cobra.Group{ID: constants.CoreGroupId, Title: constants.CoreGroupId},
		)

		// command
		cmd.AddCommand(commandCommand(app))
		// last
		cmd.AddCommand(lastCommand(app))
		// task
		cmd.AddCommand(taskCommand(app))
		// exit
		cmd.AddCommand(exitCommand(app))
		// aliases
		for _, v := range aliasCommands(app) {
			cmd.AddCommand(v)
		}

		cmd.InitDefaultHelpCmd()
		cmd.SetHelpCommandGroupID(constants.CoreGroupId)
		cmd.CompletionOptions.DisableDefaultCmd = true
		return cmd
	}
}
