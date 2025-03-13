package agent

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// Cmd returns command "agents"
func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "agents",
		Short:   "Manage agents",
		GroupID: constants.BaseGroupId,
	}

	cmd.AddCommand(
		// list all registered agents
		listCmd(c),
	)

	return cmd
}

// listCmd returns command "list" for "agents"
func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List registered agents",
		Run: func(cmd *cobra.Command, args []string) {
			agents := agent.Agents.Get()
			for _, v := range agents {
				os := v.GetOs().StringShort()
				if v.GetIsPrivileged() {
					os = color.RedString(v.GetOs().StringShort())
				}
				last := color.GreenString(utils.HumanDurationC(v.GetLast()))
				if v.IsDelayed(0) {
					// if agent is delayed on predictable time
					last = color.YellowString(utils.HumanDurationC(v.GetLast()))
				}
				if v.IsDead(0) {
					// if agent is dead (or not?)
					last = color.RedString(utils.HumanDurationC(v.GetLast()))
				}
				notificator.Print("[%s] %15s  %6s %-20s %-16s %s",
					os,
					last,
					v.GetIdHex(),
					v.GetUsername(),
					v.GetHostname(),
					v.GetIntIp(),
				)
			}
		},
	}
}
