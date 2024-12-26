package base

import (
	"os"

	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func exitCommand(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "exit",
		Short: "exit operator cli",
		Run: func(*cobra.Command, []string) {
			if utils.ExitConsolePrompt(c) {
				service.Close()
				os.Exit(0)
			}
		},
	}
}
