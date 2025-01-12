package base

import (
	"github.com/PicoTools/pico-cli/internal/commands/base/agent"
	"github.com/PicoTools/pico-cli/internal/commands/base/chat"
	"github.com/PicoTools/pico-cli/internal/commands/base/exit"
	"github.com/PicoTools/pico-cli/internal/commands/base/script"
	"github.com/PicoTools/pico-cli/internal/commands/base/whoami"
	"github.com/PicoTools/pico-cli/internal/commands/shared/use"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/reeflective/console"
	"github.com/reeflective/console/commands/readline"
	"github.com/spf13/cobra"
)

func Cmds(app *console.Console) console.Commands {
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
		cmd.AddCommand(chat.Cmd(app))
		// exit (exit from cli)
		cmd.AddCommand(exit.Cmd(app))
		// agent (manage agents)
		cmd.AddCommand(agent.Cmd(app))
		// use (switch to agent console)
		cmd.AddCommand(use.Cmd(app))
		// script (manage scripts)
		cmd.AddCommand(script.Cmd(app))
		// readline (manipulate readline variables)
		cmd.AddCommand(readline.Commands(app.Shell()))
		// whoami (get username of operator)
		cmd.AddCommand(whoami.Cmd(app))

		cmd.InitDefaultHelpCmd()
		cmd.CompletionOptions.DisableDefaultCmd = true
		return cmd
	}
}
