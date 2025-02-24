package areserved28

import (
	"fmt"

	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/agent"
	commonv1 "github.com/PicoTools/pico/pkg/proto/common/v1"
	operatorv1 "github.com/PicoTools/pico/pkg/proto/operator/v1"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const name = "a_reserved28"

func GetApiName() string {
	return name
}

// arg0 - id (int)
// arg1 - args (str)
// arg2 - blob (str)
func FrontendAgentReserved28(args ...object.Object) (object.Object, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, fmt.Errorf("expecting 1, 2 or 3 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument 'int', got '%s'", args[0].TypeName())
	}
	var argsStr string
	if len(args) > 1 {
		// TODO - add object.Null support
		v, ok := args[1].(*object.Str)
		if !ok {
			return nil, fmt.Errorf("expecting 2nd argument 'str', got '%s'", args[1].TypeName())
		}
		argsStr = v.GetValue().(string)
	}
	var blobBytes []byte
	if len(args) > 2 {
		// TODO - add object.Null support
		v, ok := args[2].(*object.Str)
		if !ok {
			return nil, fmt.Errorf("expecting 3rd argument 'str', got '%s'", args[2].TypeName())
		}
		blobBytes = []byte(v.GetValue().(string))
	}
	if err := BackendAgentReserved28(uint32(id.GetValue().(int64)), argsStr, blobBytes); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAgentReserved28(id uint32, args string, blob []byte) error {
	cap := shared.CapReserved28

	agent := agent.Agents.GetById(id)
	if agent == nil {
		return fmt.Errorf("no agent with id %d", id)
	}

	if !cap.ValidateMask(agent.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("agent doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Reserved28{
			Reserved28: &commonv1.CapReserved28{
				Args: wrapperspb.String(args),
				Blob: wrapperspb.Bytes(blob),
			},
		},
	})
}
