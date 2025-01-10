package notificator

import (
	"fmt"

	"github.com/fatih/color"
)

var printf func(format string, args ...any) (int, error)

// SetOut sets output for printing
func SetOut(fn func(string, ...any) (int, error)) {
	printf = fn
}

// Printf prints regular message by formatter without new line
func Printf(format string, args ...any) {
	printf(format, args...)
}

// Print prints regular message by formatter
func Print(format string, args ...any) {
	printf("%s\n", fmt.Sprintf(format, args...))
}

// PrintfNotify prints message with NOTIFY level without new line
func PrintfNotify(format string, args ...any) {
	printf("[%s] %s", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintNotify prints message with NOTIFY level
func PrintNotify(format string, args ...any) {
	printf("[%s] %s\n", color.CyanString("*"), fmt.Sprintf(format, args...))
}

// PrintfInfo prints message with INFO level without new line
func PrintfInfo(format string, args ...any) {
	printf("[%s] %s", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintInfo prints message with INFO level
func PrintInfo(format string, args ...any) {
	printf("[%s] %s\n", color.GreenString("+"), fmt.Sprintf(format, args...))
}

// PrintfWarning prints message with WARNING level without new line
func PrintfWarning(format string, args ...any) {
	printf("[%s] %s", color.YellowString("!"), fmt.Sprintf(format, args...))
}

// PrintWarning prints message with WARNING level
func PrintWarning(format string, args ...any) {
	printf("[%s] %s\n", color.YellowString("!"), fmt.Sprintf(format, args...))
}

// PrintfError prints message with ERROR level without new line
func PrintfError(format string, args ...any) {
	printf("[%s] %s\n", color.RedString("-"), fmt.Sprintf(format, args...))
}

// PrintError prints message with ERROR level
func PrintError(format string, args ...any) {
	printf("[%s] %s\n", color.RedString("-"), fmt.Sprintf(format, args...))
}
