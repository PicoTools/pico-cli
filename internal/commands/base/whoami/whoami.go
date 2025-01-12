package whoami

import (
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func Cmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:     "whoami",
		Short:   "Get username of operator",
		GroupID: constants.BaseGroupId,
		Run: func(*cobra.Command, []string) {
			notificator.Print("%s", service.GetUsername())
		},
	}
}
