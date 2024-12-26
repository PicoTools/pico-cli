package appid

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

const name = "a_ppid"

func GetApiName() string {
	return name
}

func FrontendAntPpid(args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expecting 2 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	ppid := object.NewInt(0)
	if len(args) == 2 {
		ppid, ok = args[1].(*object.Int)
		if !ok {
			return nil, fmt.Errorf("expecting 2nd argument int, got '%s'", args[1].TypeName())
		}
	}
	if err := BackendAntPpid(uint32(id.GetValue().(int64)), ppid.GetValue().(int64)); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendAntPpid(id uint32, ppid int64) error {
	cap := shared.CapPpid

	ant := ant.Ants.GetById(id)
	if ant == nil {
		return fmt.Errorf("no ant with id %d", id)
	}

	if !cap.ValidateMask(ant.GetCaps()) {
		return merror.BackendMessageError(id, fmt.Sprintf("ant doesn't support %s", cap.String()))
	}

	return service.NewTask(id, &operatorv1.CreateTaskRequest{
		Cap: uint32(cap),
		Args: &operatorv1.CreateTaskRequest_Ppid{
			Ppid: &commonv1.CapPpid{
				Ppid: uint64(ppid),
			},
		},
	})
}
