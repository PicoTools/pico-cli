package last

import (
	"github.com/PicoTools/pico-cli/internal/commands/agent/utils"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// Cmd returns command "last"
func Cmd(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:     "last",
		Short:   "Get output of last entered command by operator (yourself)",
		GroupID: constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			// get last command created by operator
			command := task.Commands.GetLastCommandByOperator(service.GetUsername())
			if command == nil {
				notificator.PrintWarning("you didn't enter any commands")
				return
			}
			for _, v := range command.GetData().Get() {
				utils.PrintCommandData(c, v)
			}
		},
	}
}
