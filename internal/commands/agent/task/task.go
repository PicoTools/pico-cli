package task

import (
	"os"
	"strconv"

	"github.com/PicoTools/pico-cli/internal/commands/agent/utils"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/fatih/color"
	"github.com/reeflective/console"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Cmd returns command "tasks"
func Cmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tasks",
		Short:   "Show tasks processed by agent",
		GroupID: constants.CoreGroupId,
	}

	cmd.AddCommand(
		// get output of task (capability invoke)
		getCmd(c),
		// download raw results of task in file
		downloadCmd(c),
		// list all tasks for agent
		listCmd(c),
	)

	return cmd
}

// listCmd returns command "list" for "tasks"
func listCmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List tasks (invokes of capabilities) for agent",
		Run: func(*cobra.Command, []string) {
			tasks := task.Commands.GetTasks()
			if len(tasks) == 0 {
				notificator.PrintWarning("no tasks exist yet")
				return
			}
			for _, v := range tasks {
				notificator.Print("[%s | %s] %s",
					v.GetCreatedAt().Format("02/01 15:04:05"),
					color.GreenString("%d", v.GetId()),
					v.GetCapability().String(),
				)
			}
		},
	}
}

// downloadCmd returns command "download" for "tasks"
func downloadCmd(*console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "download <id> <path>",
		Short:                 "Download output of task to file",
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// parse input
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				notificator.PrintError("invalid task's ID")
				return
			}
			// get output of task by ID
			output, err := service.GetTaskOutput(id)
			if err != nil {
				switch status.Code(err) {
				case codes.NotFound:
					notificator.PrintError("unknown task's ID")
				default:
					notificator.PrintError("%s", err.Error())
				}
				return
			}
			// write output of task to file
			if err := os.WriteFile(args[1], output, 0640); err != nil {
				notificator.PrintError("save output: %s", err.Error())
				return
			}
			notificator.PrintInfo("output saved to %s", args[1])
		},
	}

	// autocomplete
	// arg1: task's ID
	// arg2: FS
	carapace.Gen(cmd).PositionalCompletion(tasksIdsCompleter(), carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		return carapace.ActionFiles()
	}))

	return cmd
}

// getCmd returns command "get" for "tasks"
func getCmd(c *console.Console) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get output of task specified by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// parse input
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				notificator.PrintError("invalid task's ID")
				return
			}
			// get task by ID
			task := task.Commands.GetTaskById(id)
			if task == nil {
				notificator.PrintError("unknown task's ID")
				return
			}
			utils.PrintTaskData(c, task)
		},
	}

	// task's IDs autocomplete
	carapace.Gen(cmd).PositionalCompletion(tasksIdsCompleter())

	return cmd
}

// taskIdsCompleter returns completer for all task's IDs
func tasksIdsCompleter() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		var suggestions []string
		for _, v := range task.Commands.GetTasks() {
			suggestions = append(suggestions, strconv.Itoa(int(v.GetId())))
		}
		return carapace.ActionValues(suggestions...)
	})
}
