package scripts

import (
	"fmt"
	"strings"

	"github.com/PicoTools/pico-cli/internal/scripts/aliases"
	"github.com/PicoTools/pico-cli/internal/scripts/aliases/alias"
)

// ProcessCommand processes and executes alias
func ProcessCommand(bid uint32, cmd string) error {
	cmd = strings.TrimSpace(cmd)

	c := strings.Split(cmd, " ")
	if aliases.IsAliasExist(c[0]) {
		return alias.BackendAlias(bid, cmd)
	}
	return fmt.Errorf("unknown alias '%s'", c[0])
}
