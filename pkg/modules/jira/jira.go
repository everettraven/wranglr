package jira

import (
	"go.starlark.net/starlark"
)

func New() (string, starlark.Value) {
	return "jira", &Module{}
}
