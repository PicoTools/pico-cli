package agent

import (
	"github.com/PicoTools/pico-cli/internal/commands/agent/alias"
	"github.com/PicoTools/pico-cli/internal/commands/agent/command"
	"github.com/PicoTools/pico-cli/internal/commands/agent/exit"
	"github.com/PicoTools/pico-cli/internal/commands/agent/info"
	"github.com/PicoTools/pico-cli/internal/commands/agent/last"
	"github.com/PicoTools/pico-cli/internal/commands/agent/task"
	"github.com/PicoTools/pico-cli/internal/commands/shared/use"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/reeflective/console"
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
			&cobra.Group{ID: constants.AliasGroupId, Title: constants.AliasGroupId},
			&cobra.Group{ID: constants.BaseGroupId, Title: constants.BaseGroupId},
			&cobra.Group{ID: constants.CoreGroupId, Title: constants.CoreGroupId},
		)

		// command (list/show commands for agent)
		cmd.AddCommand(command.Cmd(app))
		// use (switch on agent console)
		cmd.AddCommand(use.Cmd(app))
		// last (get last command output for agent)
		cmd.AddCommand(last.Cmd(app))
		// task (list/download tasks for agent)
		cmd.AddCommand(task.Cmd(app))
		// exit (switch back on base menu)
		cmd.AddCommand(exit.Cmd(app))
		// info (print full info about agent)
		cmd.AddCommand(info.Cmd(app))
		// alias (aliases for interaction with agent)
		for _, v := range alias.Cmd(app) {
			cmd.AddCommand(v)
		}

		cmd.InitDefaultHelpCmd()
		cmd.SetHelpCommandGroupID(constants.CoreGroupId)
		cmd.CompletionOptions.DisableDefaultCmd = true
		return cmd
	}
}
