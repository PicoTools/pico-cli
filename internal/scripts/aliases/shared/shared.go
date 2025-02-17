package shared

import (
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico/pkg/shared"
)

func BackendIsOs(id uint32, os shared.AgentOs) bool {
	b := agent.Agents.GetById(id)
	if b == nil {
		return false
	}
	return b.GetOs() == os
}

func BackendIsArch(id uint32, arch shared.AgentArch) bool {
	b := agent.Agents.GetById(id)
	if b == nil {
		return false
	}
	return b.GetArch() == arch
}
