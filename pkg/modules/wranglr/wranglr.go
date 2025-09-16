package wranglr

import (
	"go.starlark.net/starlark"
)

func New(output string) (string, starlark.Value) {
	return "wranglr", &Module{Output: output}
}
