package builtins

import "go.starlark.net/starlark"

type BuiltinFunc func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)
