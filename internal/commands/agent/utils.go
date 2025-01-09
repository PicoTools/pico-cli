package agent

import (
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/storage/task"
	"github.com/reeflective/console"
)

func printTaskGroupData(c *console.Console, v task.TaskData) {
	switch data := v.(type) {
	case *task.Message:
		c.Printf("%s\n", data.String())
	case *task.Task:
		preambule := data.StringStatus()
		if preambule != "" {
			c.Printf("%s\n", preambule)
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
			c.Printf("%s\n", output)
		}
	}
}
