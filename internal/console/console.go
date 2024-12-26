package console

import (
	"context"
	"fmt"
	"io"
	"os"

	antCmd "github.com/PicoTools/pico-cli/internal/commands/ant"
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

	// for notifications
	notificator.SetOut(app)

	// base menu
	base := app.NewMenu(constants.BaseMenuName)
	base.Short = "base operator cli"
	base.Prompt().Primary = func() string { return fmt.Sprintf("[%s] > ", color.MagentaString("pico")) }
	base.AddInterrupt(io.EOF, func(c *console.Console) {
		if utils.ExitConsolePrompt(c) {
			service.Close()
			os.Exit(0)
		}
	})
	base.SetCommands(baseCmd.Commands(app))

	// ant menu
	ant := app.NewMenu(constants.AntMenuName)
	ant.Short = "ant operator cli"
	ant.SetCommands(antCmd.Commands(app))

	// switch on base menu
	app.SwitchMenu(constants.BaseMenuName)
	return app.StartContext(ctx)
}
