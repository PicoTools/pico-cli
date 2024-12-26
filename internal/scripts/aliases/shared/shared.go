package shared

import (
	"github.com/PicoTools/pico-cli/internal/storage/ant"
	"github.com/PicoTools/pico-shared/shared"
)

func BackendIsOs(id uint32, os shared.AntOs) bool {
	b := ant.Ants.GetById(id)
	if b == nil {
		return false
	}
	return b.GetOs() == os
}

func BackendIsArch(id uint32, arch shared.AntArch) bool {
	b := ant.Ants.GetById(id)
	if b == nil {
		return false
	}
	return b.GetArch() == arch
}
