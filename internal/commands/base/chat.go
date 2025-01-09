package base

import (
	"strings"

	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

func chatCommand(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "send message in chat",
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			if err := service.SendChatMessage(strings.Join(args, " ")); err != nil {
				notificator.PrintError("send message: %s", err.Error())
				return
			}
		},
	}
}
