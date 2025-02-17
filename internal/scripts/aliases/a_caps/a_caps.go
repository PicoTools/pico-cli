package acaps

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PicoTools/pico-cli/internal/storage/agent"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "a_caps"

func GetApiName() string {
	return name
}

func FrontendAgentCaps(args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expecting 1 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	caps, err := BackendAgentCaps(uint32(id.GetValue().(int64)))
	if err != nil {
		return nil, err
	}
	result := make([]object.Object, 0)
	for _, v := range caps {
		result = append(result, object.NewStr(v.String()))
	}
	return object.NewList(result), nil
}

func BackendAgentCaps(id uint32) ([]shared.Capability, error) {
	agent := agent.Agents.GetById(id)
	if agent == nil {
		return nil, fmt.Errorf("no agent with id %d", id)
	}
	caps := shared.SupportedCaps(agent.GetCaps())
	sort.Slice(caps, func(i, j int) bool {
		return strings.Compare(caps[i].String(), caps[j].String()) < 0
	})
	return caps, nil
}
