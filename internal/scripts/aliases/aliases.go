package aliases

import "github.com/PicoTools/plan/pkg/engine/object"

// storage of all aliases
var Aliases = make(map[string]*Alias)

// Clear reset storage with aliases
func Clear() {
	Aliases = make(map[string]*Alias)
}

// IsAliasExist returns true if alias exists in storage
func IsAliasExist(n string) bool {
	_, ok := Aliases[n]
	return ok
}

type Alias struct {
	description string
	usage       string
	visible     bool
	closure     *object.RuntimeFunc
}

func (a *Alias) GetDescription() string {
	return a.description
}

func (a *Alias) SetDescription(description string) {
	a.description = description
}

func (a *Alias) GetUsage() string {
	return a.usage
}

func (a *Alias) SetUsage(usage string) {
	a.usage = usage
}

func (a *Alias) GetVisible() bool {
	return a.visible
}

func (a *Alias) SetVisible(flag bool) {
	a.visible = flag
}

func (a *Alias) GetClosure() *object.RuntimeFunc {
	return a.closure
}

func (a *Alias) SetClosure(closure *object.RuntimeFunc) {
	a.closure = closure
}
