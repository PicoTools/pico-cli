package aexecassembly

import (
	"fmt"
	"os"

	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	commonv1 "github.com/PicoTools/pico/pkg/proto/common/v1"
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
	"github.com/go-faster/errors"
)

const name = "a_exec_assembly"

func GetApiName() string {
	return name
}

func FrontendAgentExecuteAssembly(args ...object.Object) (object.Object, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("expecting 2 or 3 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument 'int', got '%s'", args[0].TypeName())
	}
	path, ok := args[1].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument 'str', got '%s'", args[1].TypeName())
	}
	arg := object.NewStr("")
	if len(args) == 3 {
		arg, ok = args[2].(*object.Str)
		if !ok {
			return nil, fmt.Errorf("expecting 3rd argument 'str', got '%s'", args[2].TypeName())
		}
	}
	if err := BackendAgentExecuteAssembly(uint32(id.GetValue().(int64)), path.GetValue().(string), arg.GetValue().(string)); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentExecuteAssembly(id uint32, path, args string) error {
	cap := shared.CapExecAssembly

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	// get data from file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("no such file")
		} else if os.IsPermission(err) {
			err = errors.New("permission denied")
		}
		return merror.BackendMessageError(id, fmt.Sprintf("unable open local file by path %s: %s", path, err.Error()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_ExecAssembly{
			ExecAssembly: &commonv1.CapExecAssembly{
				Args: args,
				Blob: data,
			},
		},
	})
}
