package utils

import (
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
)

// PrintCommandData prints command's data
func PrintCommandData(c *console.Console, v task.TaskData) {
	switch data := v.(type) {
	case *task.Message:
		notificator.Print("%s", data.String())
	case *task.Task:
		PrintTaskData(c, data)
	}
}

// PrintTaskData print task's data
func PrintTaskData(_ *console.Console, data *task.Task) {
	// preambule (how much bytes and which status)
	notificator.Print("%s", data.StringStatus())
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
	// prepend '\n'
	if output[0] != '\n' {
		output = "\n" + output
	}
	// append '\n'
	if output[len(output)-1] != '\n' {
		output = output + "\n"
	}
	notificator.Print("%s", output)
}
