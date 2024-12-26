package scripts

import (
	"github.com/PicoTools/pico-cli/internal/scripts/aliases"
	"github.com/PicoTools/plan/pkg/engine"
)

// Init registers all needing API
func Init() error {
	// register API
	registerApi()
	// register builtin scripts
	if err := registerBuiltin(); err != nil {
		return err
	}
	return nil
}

// Rebuild rebuilds script's storage
func Rebuild() error {
	// clear aliases
	aliases.Clear()
	// reset runtime
	engine.Reset()
	// initialize API
	Init()
	// reregister external scripts
	externalScripts := make([]*Script, 0)
	scripts.Range(func(k string, v *Script) bool {
		externalScripts = append(externalScripts, v)
		scripts.Delete(k)
		return true
	})
	for _, v := range externalScripts {
		RegisterExternalByPath(v.path)
	}
	return nil
}
