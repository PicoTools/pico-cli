package isx86

import (
	"fmt"

	"github.com/PicoTools/pico-cli/internal/scripts/aliases/shared"
	shr "github.com/PicoTools/pico-shared/shared"
	"github.com/PicoTools/plan/pkg/engine/object"
)

const name = "is_x86"

func GetApiName() string {
	return name
}

func FrontendIsX86(args ...object.Object) (object.Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expecting 1 argument, got %d", len(args))
	}
	id, ok := args[0].(*object.Int)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument int, got '%s'", args[0].TypeName())
	}
	return object.NewBool(shared.BackendIsArch(uint32(id.GetValue().(int64)), shr.ArchX86)), nil
}
