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

func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List agents",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			agents := agent.Agents.Get()
			for _, v := range agents {
				os := v.GetOs().StringShort()
				if v.GetIsPrivileged() {
					os = color.RedString(v.GetOs().StringShort())
				}
				last := color.GreenString(utils.HumanDurationC(v.GetLast()))
				if v.IsDelay(0) {
					last = color.YellowString(utils.HumanDurationC(v.GetLast()))
				}
				if v.IsDead(0) {
					last = color.RedString(utils.HumanDurationC(v.GetLast()))
				}
				notificator.Print("[%s] (%15s) %6s %-20s %-16s %s",
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

func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "agents",
		Short:                 "Manage agents",
		DisableFlagsInUseLine: true,
		GroupID:               constants.BaseGroupId,
	}
	cmd.AddCommand(
		listCmd(c),
	)
	return cmd
}
