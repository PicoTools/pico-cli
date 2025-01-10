package asleep

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

const name = "a_sleep"

func GetApiName() string {
	return name
}

func FrontendAgentSleep(args ...object.Object) (object.Object, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("expecting 2 or 3 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	sleep, ok := args[1].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument int, got '%s'", args[1].TypeName())
	}
	jitter := object.NewInt(0)
	if len(args) == 3 {
		jitter, ok = args[2].(*object.Int)
		if !ok {
			return nil, fmt.Errorf("expecting 3rd argument int, got '%s'", args[2].TypeName())
		}
	}
	if err := BackendAgentSleep(uint32(id.GetValue().(int64)), uint32(sleep.GetValue().(int64)), uint32(jitter.GetValue().(int64))); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentSleep(id uint32, sleep uint32, jitter uint32) error {
	cap := shared.CapSleep

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Sleep{
			Sleep: &commonv1.CapSleep{
				Sleep:  sleep,
				Jitter: jitter,
			},
		},
	})
}
