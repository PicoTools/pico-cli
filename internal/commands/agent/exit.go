package agent

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
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
			if err := service.UnpollAgentTasks(agent.ActiveAgent); err != nil {
				color.Yellow("unable stop polling tasks for agent: %s", err.Error())
			}
			task.ResetStorage()
			agent.ActiveAgent = nil
			c.SwitchMenu(constants.BaseMenuName)
		},
	}
}
