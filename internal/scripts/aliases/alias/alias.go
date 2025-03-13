package alias

import (
	"fmt"
	"strings"

	"github.com/PicoTools/pico-cli/internal/scripts/aliases"
	"github.com/PicoTools/pico-cli/internal/service"
	"github.com/PicoTools/plan/pkg/engine/object"
	"github.com/PicoTools/plan/pkg/engine/visitor"
	"github.com/fatih/color"
	"github.com/go-faster/errors"
	"github.com/google/shlex"
)

const name = "alias"

func GetApiName() string {
	return name
}

// FrontendAlias performes new alias registration
// args[0] - alias name
// args[1] - closure which executes on alias invoke
// args[2] - alias'es description
// args[3] - alias'es usage
// args[4] - is command visible for other operators
func FrontendAlias(args ...object.Object) (object.Object, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("expecting 5 arguments, got %d", len(args))
	}
	name, ok := args[0].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 1st argument 'str', got '%s'", args[0].TypeName())
	}
	closure, ok := args[1].(*object.RuntimeFunc)
	if !ok {
		return nil, fmt.Errorf("expecting 2nd argument 'closure', got '%s'", args[1].TypeName())
	}
	description, ok := args[2].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 3rd argument 'str', got '%s'", args[2].TypeName())
	}
	usage, ok := args[3].(*object.Str)
	if !ok {
		return nil, fmt.Errorf("expecting 4th argument 'str', got '%s'", args[3].TypeName())
	}
	visible, ok := args[4].(*object.Bool)
	if !ok {
		return nil, fmt.Errorf("expecting 5th argument 'bool', got '%s'", args[4].TypeName())
	}
	newAlias := &aliases.Alias{}
	newAlias.SetDescription(description.GetValue().(string))
	newAlias.SetUsage(usage.GetValue().(string))
	newAlias.SetVisible(visible.GetValue().(bool))
	newAlias.SetClosure(closure)
	// save alias in storage
	aliases.Aliases[name.GetValue().(string)] = newAlias
	return object.NewNull(), nil
}

// BackendAlias invoke alias'es closure
func BackendAlias(id uint32, cmd string) error {
	s, err := shlex.Split(cmd)
	if err != nil {
		return errors.Wrap(err, "split command")
	}

	// get alias
	al, ok := aliases.Aliases[s[0]]
	if !ok {
		return fmt.Errorf("unknown alias '%s'", s[0])
	}

	// create command
	if err := service.NewCommand(id, cmd, al.GetVisible()); err != nil {
		return errors.Wrap(err, "create command")
	}
	defer func(id uint32) {
		// close command
		if err := service.CloseCommand(id); err != nil {
			color.Red(err.Error())
		}
	}(id)

	// bid
	arg0 := object.NewInt(int64(id))
	// cmd
	arg1 := object.NewStr(s[0])
	var temp []object.Object
	if len(s) != 1 {
		for i := 1; i < len(s); i++ {
			temp = append(temp, object.NewStr(s[i]))
		}
	}
	// args
	arg2 := object.NewList(temp)
	// raw
	arg3 := object.NewStr("")
	if len(strings.Split(cmd, " ")) != 1 {
		arg3 = object.NewStr(strings.Join(strings.Split(cmd, " ")[1:], " "))
	}

	// invoke closure for alias
	v := visitor.NewVisitor()
	v.InvokeRuntimeFunc(al.GetClosure(), arg0, arg1, arg2, arg3)
	if v.GetError() != nil {
		return v.GetError()
	}
	return nil
}
