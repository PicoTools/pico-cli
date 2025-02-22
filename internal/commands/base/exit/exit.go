package exit

import (
	"os"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// Cmd returns command "exit"
func Cmd(c *console.Console) *cobra.Command {
	return &cobra.Command{
		Use:     "exit",
		Short:   "Exit operator's CLI",
		GroupID: constants.BaseGroupId,
		Run: func(*cobra.Command, []string) {
			if utils.ExitConsolePrompt(c) {
				_ = service.Close()
				os.Exit(0)
			}
		},
	}
}
