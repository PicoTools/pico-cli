package notificator

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/reeflective/console"
)

var out *console.Console

// SetOut sets output for printing
func SetOut(c *console.Console) {
	out = c
}

// PrintChat prints chat's message in console
func PrintChat(format string, args ...any) {
	out.TransientPrintf("[%s] %s", color.CyanString("+"), fmt.Sprintf(format, args...))
}
