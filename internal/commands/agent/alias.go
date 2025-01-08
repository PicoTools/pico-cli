package agent

import (
	"fmt"
	"strings"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/scripts"
	"github.com/PicoTools/pico-cli/internal/scripts/aliases"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func aliasCommands(*console.Console) []*cobra.Command {
	cmds := make([]*cobra.Command, 0)
	for k, v := range aliases.Aliases {
		cmd := &cobra.Command{
			Use:                   k,
			Short:                 v.GetDescription(),
			GroupID:               constants.AliasGroupId,
			DisableFlagsInUseLine: true,
			DisableFlagParsing:    true,
			Run: func(cmd *cobra.Command, args []string) {
				rawCmd := k + " " + strings.Join(args, " ")
				if err := scripts.ProcessCommand(agent.ActiveAgent.GetId(), rawCmd); err != nil {
					color.Red(err.Error())
				}
			},
		}
		cmd.SetHelpTemplate(fmt.Sprintf("%s\n\n%s\n", v.GetDescription(), v.GetUsage()))
		cmds = append(cmds, cmd)
	}
	return cmds
}
