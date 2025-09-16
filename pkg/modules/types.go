package modules

import (
	"fmt"

	"go.starlark.net/starlark"
)

type BuiltinFunc func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

var modules = map[string]starlark.Value{}

func Register(name string, module starlark.Value) error {
	if _, ok := modules[name]; ok {
		return fmt.Errorf("module %q is already registered", name)
	}

	modules[name] = module

	return nil
}

func Modules() map[string]starlark.Value {
	return modules
}
