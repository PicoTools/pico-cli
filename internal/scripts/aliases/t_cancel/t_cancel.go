package tcancel

import (
	"fmt"

	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "t_cancel"

func GetApiName() string {
	return name
}

func FrontendTasksCancel(args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expecting 1 argument, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument 'int', got '%s'", args[0].TypeName())
	}
	if err := BackendTasksCancel(uint32(id.GetValue().(int64))); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendTasksCancel(id uint32) error {
	return service.CancelTasks(id)
}
