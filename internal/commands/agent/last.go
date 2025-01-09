package agent

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func lastCommand(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "last",
		Aliases:               []string{"l"},
		Short:                 "get output of last task",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			tg := task.Commands.GetLast()
			if tg == nil {
				notificator.PrintInfo("no tasks")
				return
			}
			for _, v := range tg.GetData().Get() {
				printTaskGroupData(c, v)
			}
		},
	}
}
