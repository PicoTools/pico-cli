package agent

import (
	"fmt"
	"strings"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func infoCommand(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:                   "info",
		Short:                 "Get full information about agent",
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
		Run: func(cmd *cobra.Command, args []string) {
			agent := agent.ActiveAgent
			var result strings.Builder
			result.WriteString(fmt.Sprintf("%-16s %s\n", "ID:", agent.GetIdHex()))
			result.WriteString(fmt.Sprintf("%-16s %v\n", "Privileged:", agent.GetIsPrivileged()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Architecture:", agent.GetArch().String()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "OS:", agent.GetOs().String()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "OS meta:", agent.GetOsMeta()))
			result.WriteString(fmt.Sprintf("%-16s %d\n", "PID:", agent.GetPid()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "External IP:", agent.GetExtIp()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Internal IP:", agent.GetIntIp()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Domain:", agent.GetDomain()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Hostname:", agent.GetHostname()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Username:", agent.GetUsername()))
			result.WriteString(fmt.Sprintf("%-16s %s\n", "Process name:", agent.GetProcessName()))
			result.WriteString(fmt.Sprintf("%-16s %ds (%d%%)\n", "Sleep/Jitter:", agent.GetSleep(), agent.GetJitter()))
			result.WriteString(fmt.Sprintf("%-16s %s (%s)\n", "First checkout:", agent.GetFirst().Format("2006/01/02 15:04:05"), utils.HumanDurationC(agent.GetFirst())))
			result.WriteString(fmt.Sprintf("%-16s %s (%s)\n", "Last checkout:", agent.GetLast().Format("2006/01/02 15:04:05"), utils.HumanDurationC(agent.GetLast())))
			notificator.Printf("%s", result.String())
		},
	}
}
