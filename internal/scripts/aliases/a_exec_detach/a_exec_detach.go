package aexecdetach

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

const name = "a_exec_detach"

func GetApiName() string {
	return name
}

func FrontendAgentExecDetach(args ...object.Object) (object.Object, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("expecting 2 or 3 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	cmd, ok := args[1].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument str, got '%s'", args[1].TypeName())
	}
	arg := object.NewStr("")
	if len(args) == 3 {
		arg, ok = args[2].(*object.Str)
		if !ok {
			return nil, fmt.Errorf("expecting 3rd argument str, got '%s'", args[2].TypeName())
		}
	}
	if err := BackendAgentExecDetach(uint32(id.GetValue().(int64)), cmd.GetValue().(string), arg.GetValue().(string)); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentExecDetach(id uint32, cmd string, args string) error {
	cap := shared.CapExecDetach

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_ExecDetach{
			ExecDetach: &commonv1.CapExecDetach{
				Cmd:  cmd,
				Args: args,
			},
		},
	})
}
