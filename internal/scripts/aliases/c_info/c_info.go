package cinfo

import (
	"fmt"

	"github.com/PicoTools/pico-cli/internal/notificator"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "c_info"

func GetApiName() string {
	return name
}

func FrontendConsoleInfo(args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expecting 1 arguments, got %d", len(args))
	}
	msg, ok := args[0].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expectign 1st argument str, got '%s'", args[1].TypeName())
	}
	notificator.PrintInfo("%s", msg.GetValue().(string))
	return object.NewNull(), nil
}