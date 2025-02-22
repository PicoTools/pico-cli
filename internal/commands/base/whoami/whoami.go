package whoami

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// Cmd returns command "whoami"
func Cmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:     "whoami",
		Short:   "Get my operator's username",
		GroupID: constants.BaseGroupId,
		Run: func(*cobra.Command, []string) {
			notificator.Print("%s", service.GetUsername())
		},
	}
}
