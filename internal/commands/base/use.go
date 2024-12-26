package base

import (
	"fmt"
	"strconv"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/ant"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

func useCommand(c *console.Console) *cobra.Command {
	useCmd := &cobra.Command{
		Use:                   "use",
		Short:                 "switch on ant shell",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 16, 32)
			if err != nil {
				color.Red("invalid ant id")
				return
			}
			a := ant.Ants.GetById(uint32(id))
			if a == nil {
				color.Red("unknown ant id")
				return
			}
			if err := service.PollAntTasks(a); err != nil {
				color.Red("unable start polling tasks for ant: %s", err.Error())
				return
			}
			ant.ActiveAnt = a
			c.Menu(constants.AntMenuName).Prompt().Primary = func() string { return fmt.Sprintf("[%s] > ", color.MagentaString(args[0])) }
			c.SwitchMenu(constants.AntMenuName)
		},
	}
	carapace.Gen(useCmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range ant.Ants.Get() {
			suggestions = append(suggestions, v.GetIdHex())
		}
		return carapace.ActionValues(suggestions...)
	}))
	return useCmd
}
