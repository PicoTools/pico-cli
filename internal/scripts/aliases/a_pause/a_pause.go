package apause

import (
	"fmt"

	merror "github.com/PicoTools/pico-cli/internal/scripts/aliases/m_error"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico-cli/internal/storage/ant"
	commonv1 "github.com/PicoTools/pico-shared/proto/gen/common/v1"
	operatorv1 "github.com/PicoTools/pico-shared/proto/gen/operator/v1"
	"github.com/PicoTools/pico-shared/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "a_pause"

func GetApiName() string {
	return name
}

func FrontendAntPause(args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expecting 2 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	pause, ok := args[1].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument int, got '%s'", args[1].TypeName())
	}
	if err := BackendAntPause(uint32(id.GetValue().(int64)), uint32(pause.GetValue().(int64))); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAntPause(id uint32, pause uint32) error {
	cap := shared.CapPause

	ant := ant.Ants.GetById(id)
	if ant == nil {
		return fmt.Errorf("no ant with id %d", id)
	}

	if !cap.ValidateMask(ant.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("ant doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Pause{
			Pause: &commonv1.CapPause{
				Pause: pause,
			},
		},
	})
}
