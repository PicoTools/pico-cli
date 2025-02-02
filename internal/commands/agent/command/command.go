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

func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "commands",
		Short:                 "Show commands for agent",
		Aliases:               []string{"t"},
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				for _, v := range task.Commands.Get() {
					notificator.Print("[%s] (%d) %s: %s",
						v.GetCreatedAt().Format("01/02 15:04:05"),
						v.GetId(),
						v.GetAuthor(),
						v.GetCmd(),
					)
				}
				return
			}
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
