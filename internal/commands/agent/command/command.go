package command

import (
	"strconv"

	"github.com/PicoTools/pico-cli/internal/commands/agent/utils"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List commands for agent",
		Run: func(*cobra.Command, []string) {
			commands := task.Commands.Get()
			if len(commands) == 0 {
				notificator.PrintWarning("no commands exist yet")
				return
			}
			for _, v := range commands {
				notificator.Print("[%s] (%d) %s: %s",
					v.GetCreatedAt().Format("01/02 15:04:05"),
					v.GetId(),
					v.GetAuthor(),
					v.GetCmd(),
				)
			}
		},
	}
}

func getCmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "get <task_id>",
		Short:                 "Get output for command",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				notificator.PrintError("invalid task id")
				return
			}
			tg := task.Commands.GetById(id)
			if tg == nil {
				notificator.PrintError("unknown task id")
				return
			}
			for _, v := range tg.GetData().Get() {
				utils.PrintCommandData(c, v)
			}
			return
		},
	}
	carapace.Gen(cmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range task.Commands.Get() {
			suggestions = append(suggestions, strconv.Itoa(int(v.GetId())))
		}
		return carapace.ActionValues(suggestions...)
	}))
	return cmd
}

func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "commands",
		Short:                 "Show commands for agent",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
	}
	cmd.AddCommand(
		getCmd(c),
		listCmd(c),
	)
	return cmd
}
