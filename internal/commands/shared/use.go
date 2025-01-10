package shared

import (
	"fmt"
	"strconv"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

func UseCommand(c *console.Console) *cobra.Command {
	useCmd := &cobra.Command{
		Use:                   "use",
		Short:                 "Switch on agent shell",
		DisableFlagsInUseLine: true,
		GroupID:               constants.BaseGroupId,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 16, 32)
			if err != nil {
				notificator.PrintError("invalid agent id")
				return
			}
			a := agent.Agents.GetById(uint32(id))
			if a == nil {
				notificator.PrintError("unknown agent id")
				return
			}
			if agent.ActiveAgent != nil {
				if err := service.UnpollAgentTasks(agent.ActiveAgent); err != nil {
					notificator.PrintWarning("unable stop polling tasks for agent: %s", err.Error())
				}
				task.ResetStorage()
			}
			if err := service.PollAgentTasks(a); err != nil {
				notificator.PrintError("unable start polling tasks for agent: %s", err.Error())
				return
			}
			agent.ActiveAgent = a
			c.Menu(constants.AgentMenuName).Prompt().Primary = func() string { return fmt.Sprintf("[%s] > ", color.MagentaString(args[0])) }
			c.SwitchMenu(constants.AgentMenuName)
		},
	}
	carapace.Gen(useCmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range agent.Agents.Get() {
			suggestions = append(suggestions, v.GetIdHex())
		}
		return carapace.ActionValues(suggestions...)
	}))
	return useCmd
}
