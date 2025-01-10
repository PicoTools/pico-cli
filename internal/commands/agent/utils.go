package agent

import (
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
)

// printCommandData prints command's data
func printCommandData(_ *console.Console, v task.TaskData) {
	switch data := v.(type) {
	case *task.Message:
		notificator.Print("%s", data.String())
	case *task.Task:
		preambule := data.StringStatus()
		if preambule != "" {
			notificator.Print("%s", preambule)
		}
		if data.GetOutputLen() == 0 {
			return
		}
		if data.GetIsOutputBig() {
			notificator.PrintWarning("output too big, use: tasks download %d <path to save>", data.GetId())
			return
		}
		if data.GetIsBinary() {
			notificator.PrintWarning("output is possible binary, use: tasks download %d <path to save>", data.GetId())
			return
		}
		output := data.GetOutputString()
		if output != "" {
			notificator.Print("%s", output)
		}
	}
}
