package mwarning

import (
	"fmt"

	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/pico/pkg/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "m_warning"

func GetApiName() string {
	return name
}

func FrontendMessageWarning(args ...object.Object) (object.Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expecting 2 arguments, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	msg, ok := args[1].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expectign 2nd argument str, got '%s'", args[1].TypeName())
	}
	if err := BackendMessageWarning(uint32(id.GetValue().(int64)), msg.GetValue().(string)); err != nil {
		return nil, err
	}
	return object.NewNull(), nil
}

func BackendMessageWarning(id uint32, message string) error {
	return service.NewCommandMessage(id, shared.WarningMessage, message)
}
