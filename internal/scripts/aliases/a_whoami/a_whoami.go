package awhoami

import (
	"fmt"

	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	commonv1 "github.com/PicoTools/pico-shared/proto/gen/common/v1"
	operatorv1 "github.com/PicoTools/pico-shared/proto/gen/operator/v1"
	"github.com/PicoTools/pico-shared/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "a_whoami"

func GetApiName() string {
	return name
}

func FrontendAgentWhoami(args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expecting 1 argument, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	if err := BackendAgentWhoami(uint32(id.GetValue().(int64))); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentWhoami(id uint32) error {
	cap := shared.CapWhoami

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Whoami{
			Whoami: &commonv1.CapWhoami{},
		},
	})
}
