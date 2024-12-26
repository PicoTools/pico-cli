package ant

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/ant"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func exitCommand(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "exit",
		Short:                 "switch back on base console",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.UnpollAntTasks(ant.ActiveAnt); err != nil {
				color.Yellow("unable stop polling tasks for ant: %s", err.Error())
			}
			task.ResetStorage()
			ant.ActiveAnt = nil
			c.SwitchMenu(constants.BaseMenuName)
		},
	}
}
