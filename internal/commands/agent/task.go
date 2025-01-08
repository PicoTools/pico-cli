package agent

import (
	"os"
	"strconv"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func taskDownloadCommand(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "download <task_id> <path>",
		Short:                 "download output of task to file",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				color.Red("invalid task id")
				return
			}
			output, err := service.GetTaskOutput(id)
			if err != nil {
				switch status.Code(err) {
				case codes.NotFound:
					color.Red("unknown task id")
				default:
					color.Red(err.Error())
				}
				return
			}
			if err := os.WriteFile(args[1], output, 0644); err != nil {
				color.Red("save output: %s", err.Error())
				return
			}
			color.Green("output saved to %s", args[1])
		},
	}
	// autocomplete
	// arg1: task id
	// arg2: fs
	carapace.Gen(cmd).PositionalCompletion(carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range task.Commands.GetTasks() {
			suggestions = append(suggestions, strconv.Itoa(int(v.GetId())))
		}
		return carapace.ActionValues(suggestions...)
	}), carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.ActionFiles()
	}))
	return cmd
}

func taskCommand(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "tasks",
		Short:                 "show tasks for agent",
		Aliases:               []string{"t"},
		DisableFlagsInUseLine: true,
		GroupID:               constants.CoreGroupId,
	}
	cmd.AddCommand(
		taskDownloadCommand(c),
	)
	return cmd
}
