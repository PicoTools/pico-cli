package command

import (
	"strconv"
	"strings"

	"github.com/PicoTools/pico-cli/internal/commands/agent/utils"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

// Cmd returns command "commands"
func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commands",
		Short:   "Show commands entered for agent",
		GroupID: constants.CoreGroupId,
	}

	cmd.AddCommand(
		// get output of command
		getCmd(c),
		// list all commands enetered for agent
		listCmd(c),
	)

	return cmd
}

// listCmd return command "list" for "commands"
func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List entered commands for agent",
		Run: func(*cobra.Command, []string) {
			commands := task.Commands.Get()
			if len(commands) == 0 {
				notificator.PrintWarning("no commands entered yet")
				return
			}
			for _, v := range commands {
				if strings.Compare(v.GetAuthor(), service.GetUsername()) == 0 {
					// this command entered by current operator
					notificator.Print("[%s | %s] %s",
						v.GetCreatedAt().Format("02/01 15:04:05"),
						color.GreenString("%d", v.GetId()),
						v.GetCmd(),
					)
				} else {
					// this command entered by another operator
					notificator.Print("[%s | %s] %s: %s",
						v.GetCreatedAt().Format("02/01 15:04:05"),
						color.GreenString("%d", v.GetId()),
						color.RedString("%s", v.GetAuthor()),
						v.GetCmd(),
					)
				}
			}
		},
	}
}

// getCmd returns command "get" for "commands"
func getCmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get output for command specified by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// parse input
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				notificator.PrintError("invalid command's ID")
				return
			}
			// get command by ID
			command := task.Commands.GetById(id)
			if command == nil {
				notificator.PrintError("unknown command's ID")
				return
			}
			for _, v := range command.GetData().Get() {
				utils.PrintCommandData(c, v)
			}
		},
	}

	// command's IDs autocomplete
	carapace.Gen(cmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range task.Commands.Get() {
			suggestions = append(suggestions, strconv.Itoa(int(v.GetId())))
		}
		return carapace.ActionValues(suggestions...)
	}))

	return cmd
}
