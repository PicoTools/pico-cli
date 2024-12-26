package base

import (
	"fmt"

	"github.com/PicoTools/pico-cli/internal/storage/ant"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func antListCommand(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "list ants",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			ants := ant.Ants.Get()
			for _, v := range ants {
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
				fmt.Printf("[%s] (%15s) %6s %-20s %-16s %s\n",
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

func antCommand(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ants",
		Short:                 "manage ants",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		antListCommand(c),
	)
	return cmd
}
