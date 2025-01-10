package agent

import (
	"os"
	"strconv"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func taskDownloadCommand(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "download <task_id> <path>",
		Short:                 "Download output of task to file",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				notificator.PrintError("invalid task id")
				return
			}
			output, err := service.GetTaskOutput(id)
			if err != nil {
				switch status.Code(err) {
				case codes.NotFound:
					notificator.PrintError("unknown task id")
				default:
					notificator.PrintError("%s", err.Error())
				}
				return
			}
			if err := os.WriteFile(args[1], output, 0644); err != nil {
				notificator.PrintError("save output: %s", err.Error())
				return
			}
			notificator.PrintInfo("output saved to %s", args[1])
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

func taskListCommand(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list tasks (with capabilities) for agent",
		Run: func(*cobra.Command, []string) {
			for _, v := range task.Commands.GetTasks() {
				notificator.Print("[%s] (%d) %s",
					v.GetCreatedAt().Format("01/02 15:04:05"),
					v.GetId(),
					v.GetCapability().String(),
				)
			}
		},
	}
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
		taskListCommand(c),
	)
	return cmd
}
