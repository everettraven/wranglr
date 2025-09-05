package builtins

import "go.starlark.net/starlark"

func TypeFromStarlarkList[T any](list *starlark.List) []T{
	types := []T{}

	if list != nil {
		for v := range list.Elements() {
			if typeInstance, ok := v.(T); ok {
				types = append(types, typeInstance)
			}
		}
	}

	return types
}
