package console

import (
	"context"
	"fmt"
	"io"
	"os"

	agentCmd "github.com/PicoTools/pico-cli/internal/commands/agent"
	baseCmd "github.com/PicoTools/pico-cli/internal/commands/base"
	"github.com/PicoTools/pico-cli/internal/constants"
	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/fatih/color"
	"github.com/reeflective/console"
)

func Run(ctx context.Context) error {
	app := console.New("pico-cli")

	// apply tweaks on shell
	applyTweaks(app.Shell())

	// print functions for notifications
	notificator.SetPreOut(app.TransientPrintf)
	notificator.SetPostOut(app.Printf)

	// base menu
	base := app.NewMenu(constants.BaseMenuName)
	base.Short = "base operator cli"
	base.Prompt().Primary = func() string { return fmt.Sprintf("[%s] > ", color.HiCyanString("pico")) }
	base.AddInterrupt(io.EOF, func(c *console.Console) {
		if utils.ExitConsolePrompt(c) {
			service.Close()
			os.Exit(0)
		}
	})
	base.SetCommands(baseCmd.Cmds(app))

	// agent menu
	agent := app.NewMenu(constants.AgentMenuName)
	agent.Short = "agent operator cli"
	agent.SetCommands(agentCmd.Cmds(app))

	// switch on base menu
	app.SwitchMenu(constants.BaseMenuName)

	return app.StartContext(ctx)
}
