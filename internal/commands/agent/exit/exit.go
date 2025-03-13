package exit

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func Cmd(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "exit",
		Short:                 "Switch back on base console",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.UnpollAgentTasks(agent.GetActiveAgent()); err != nil {
				notificator.PrintWarning("unable stop polling tasks for agent: %s", err.Error())
			}
			task.ResetStorage()
			agent.SetActiveAgent(nil)
			c.SwitchMenu(constants.BaseMenuName)
		},
	}
}
