package acd

import (
	"fmt"

	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	commonv1 "github.com/PicoTools/pico/pkg/proto/common/v1"
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "a_cd"

func GetApiName() string {
	return name
}

func FrontendAgentCd(args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expecting 2 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument 'int', got '%s'", args[0].TypeName())
	}
	path, ok := args[1].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument 'str', got '%s'", args[1].TypeName())
	}
	if err := BackendAgentCd(uint32(id.GetValue().(int64)), path.GetValue().(string)); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentCd(id uint32, path string) error {
	cap := shared.CapCd

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Cd{
			Cd: &commonv1.CapCd{
				Path: path,
			},
		},
	})
}
