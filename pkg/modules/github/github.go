package github

import (
	"go.starlark.net/starlark"
)

func New() (string, starlark.Value) {
	return "github", &Module{}
}
