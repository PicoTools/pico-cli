package notificator

import (
	"fmt"

	"github.com/fatih/color"
)

var prePrintf func(format string, args ...any) (int, error)
var postPrintf func(format string, args ...any) (int, error)

// SetPreOut sets function to print content below console cursor
func SetPreOut(fn func(string, ...any) (int, error)) {
	prePrintf = fn
}

// SetPostOut sets function to print content under console cursor
func SetPostOut(fn func(string, ...any) (int, error)) {
	postPrintf = fn
}

// PostPrintf prints regular message by formatter wihtout new line
func PostPrintf(format string, args ...any) {
	postPrintf(format, args...)
}

// PostPrint prints regular message by formatter
func PostPrint(format string, args ...any) {
	postPrintf("%s\n", fmt.Sprintf(format, args...))
}

// Printf prints regular message by formatter without new line
func Printf(format string, args ...any) {
	prePrintf(format, args...)
}

// Print prints regular message by formatter
func Print(format string, args ...any) {
	prePrintf("%s\n", fmt.Sprintf(format, args...))
}

// PrintfNotify prints message with NOTIFY level without new line
func PrintfNotify(format string, args ...any) {
	prePrintf("[%s] %s", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintNotify prints message with NOTIFY level
func PrintNotify(format string, args ...any) {
	prePrintf("[%s] %s\n", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintfInfo prints message with INFO level without new line
func PrintfInfo(format string, args ...any) {
	prePrintf("[%s] %s", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintInfo prints message with INFO level
func PrintInfo(format string, args ...any) {
	prePrintf("[%s] %s\n", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintfWarning prints message with WARNING level without new line
func PrintfWarning(format string, args ...any) {
	prePrintf("[%s] %s", color.YellowString("!"), fmt.Sprintf(format, args...))
}

// PrintWarning prints message with WARNING level
func PrintWarning(format string, args ...any) {
	prePrintf("[%s] %s\n", color.YellowString("!"), fmt.Sprintf(format, args...))
}

// PrintfError prints message with ERROR level without new line
func PrintfError(format string, args ...any) {
	prePrintf("[%s] %s\n", color.RedString("-"), fmt.Sprintf(format, args...))
}

// PrintError prints message with ERROR level
func PrintError(format string, args ...any) {
	prePrintf("[%s] %s\n", color.RedString("-"), fmt.Sprintf(format, args...))
}
