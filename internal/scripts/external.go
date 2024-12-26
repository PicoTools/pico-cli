package scripts

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/PicoTools/pico-cli/internal/utils"
	"github.com/PicoTools/plan/pkg/engine/types"
	mlanUtils "github.com/PicoTools/plan/pkg/engine/utils"
	"github.com/PicoTools/plan/pkg/engine/visitor"
	"github.com/antlr4-go/antlr/v4"
	"github.com/go-faster/errors"
	"github.com/lrita/cmap"
)

var scripts cmap.Map[string, *Script]

type Script struct {
	path    string
	tree    antlr.ParseTree
	addedAt time.Time
}

func (s *Script) GetPath() string {
	return s.path
}

func (s *Script) SetPath(data string) {
	s.path = data
}

func (s *Script) GetTree() antlr.ParseTree {
	return s.tree
}

func (s *Script) SetTree(t antlr.ParseTree) {
	s.tree = t
}

func (s *Script) GetAddedAt() time.Time {
	return s.addedAt
}

func (s *Script) SetAddedAt(t time.Time) {
	s.addedAt = t
}

// GetScripts returns list of registered external scripts
func GetScripts() []*Script {
	temp := make([]*Script, 0)
	scripts.Range(func(k string, v *Script) bool {
		temp = append(temp, v)
		return true
	})
	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].addedAt.Before(temp[j].addedAt)
	})
	return temp
}

// IsExternalScriptExists returns true if external scripts already exists in storage
func IsExternalScriptExists(path string) bool {
	_, ok := scripts.Load(path)
	return ok
}

// RemoveExternalByPath removes external script by pass and rebuild script's storage
func RemoveExternalByPath(path string) error {
	var err error

	if path, err = utils.GetAbsPath(path); err != nil {
		return fmt.Errorf("unable get absolute path of script: %s", err)
	}

	script, ok := scripts.Load(path)
	if !ok {
		return fmt.Errorf("script %s not registered", path)
	}

	scripts.Delete(script.GetPath())
	if err := Rebuild(); err != nil {
		return err
	}
	return nil
}

// ReloadExternalByPath reloads content of script by path
func ReloadExternalByPath(path string) error {
	var err error

	if path, err = utils.GetAbsPath(path); err != nil {
		return fmt.Errorf("unable get absolute path of script: %s", err)
	}

	script, ok := scripts.Load(path)
	if !ok {
		return fmt.Errorf("script %s not registered", path)
	}

	temp, err := processExternalScript(script.GetPath())
	if err != nil {
		return err
	}

	scripts.Delete(path)
	if err := Rebuild(); err != nil {
		return err
	}

	scripts.Store(path, temp)
	return nil
}

// RegisterExternalByPath registers external script by path
func RegisterExternalByPath(path string) error {
	var err error

	if path, err = utils.GetAbsPath(path); err != nil {
		return fmt.Errorf("unable get absolute path of script: %s", err)
	}

	if IsExternalScriptExists(path) {
		return fmt.Errorf("script %s already registered. Reload it manually", path)
	}

	temp, err := processExternalScript(path)
	if err != nil {
		return err
	}

	scripts.Store(path, temp)
	return nil
}

// processExternalScript processes external script by its path
func processExternalScript(path string) (*Script, error) {
	temp := &Script{
		path:    path,
		tree:    nil,
		addedAt: time.Now(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read script file")
	}

	temp.tree, err = mlanUtils.CreateAST(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "create ast")
	}

	v := visitor.NewVisitor()
	if res := v.Visit(temp.tree); res != types.Success {
		return nil, v.GetError()
	}

	return temp, nil
}
