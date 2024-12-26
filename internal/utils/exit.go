package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/reeflective/console"
)

// ExitConsolePrompt asks operator to submit console exit
func ExitConsolePrompt(c *console.Console) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Confirm exit (y/n): ")
	text, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(text)

	if (answer == "Y") || (answer == "y") {
		return true
	}
	return false
}
