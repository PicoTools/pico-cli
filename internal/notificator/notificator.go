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

// PrintfNotify prints message with NOTIFY level without new line
func PrintfNotify(format string, args ...any) {
	out.TransientPrintf("[%s] %s", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintNotify prints message with NOTIFY level
func PrintNotify(format string, args ...any) {
	out.TransientPrintf("[%s] %s\n", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintInfo prints message with INFO level
func PrintInfo(format string, args ...any) {
	out.TransientPrintf("[%s] %s\n", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintfInfo prints message with INFO level without new line
func PrintfInfo(format string, args ...any) {
	out.TransientPrintf("[%s] %s", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintWarning prints message with WARNING level
func PrintWarning(format string, args ...any) {
	out.TransientPrintf("[%s] %s\n", color.YellowString("!"), fmt.Sprintf(format, args...))
}

// PrintError prints message with ERROR level
func PrintError(format string, args ...any) {
	out.TransientPrintf("[%s] %s\n", color.RedString("-"), fmt.Sprintf(format, args...))
}
