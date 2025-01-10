package base

import (
	"github.com/PicoTools/pico-cli/internal/commands/shared"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/reeflective/console"
	"github.com/reeflective/console/commands/readline"
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
			&cobra.Group{ID: constants.BaseGroupId, Title: constants.BaseGroupId},
		)

		// chat (send messages to chat)
		cmd.AddCommand(chatCommand(app))
		// exit (exit from cli)
		cmd.AddCommand(exitCommand(app))
		// agent (manage agents)
		cmd.AddCommand(agentCommand(app))
		// use (switch to agent console)
		cmd.AddCommand(shared.UseCommand(app))
		// script (manage scripts)
		cmd.AddCommand(scriptCommand(app))
		// readline (manipulate readline variables)
		cmd.AddCommand(readline.Commands(app.Shell()))

		cmd.InitDefaultHelpCmd()
		cmd.CompletionOptions.DisableDefaultCmd = true
		return cmd
	}
}
