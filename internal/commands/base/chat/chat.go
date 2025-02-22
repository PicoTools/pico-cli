package chat

import (
	"strings"

	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// Cmd returns command "chat"
func Cmd(*console.Console) *cobra.Command {
	return &cobra.Command{
		Use:     "chat",
		Short:   "Send message in chat",
		GroupID: constants.BaseGroupId,
		Args:    cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			if err := service.SendChatMessage(strings.Join(args, " ")); err != nil {
				notificator.PrintError("send message: %s", err.Error())
				return
			}
		},
	}
}
