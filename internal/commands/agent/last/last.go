package last

import (
	"github.com/PicoTools/pico-cli/internal/commands/agent/utils"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func Cmd(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "last",
		Aliases:               []string{"l"},
		Short:                 "Get output of last task",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			tg := task.Commands.GetLast()
			if tg == nil {
				notificator.PrintInfo("no commands exist yet")
				return
			}
			for _, v := range tg.GetData().Get() {
				utils.PrintCommandData(c, v)
			}
		},
	}
}
