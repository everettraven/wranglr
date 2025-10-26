package wranglr

import (
	"fmt"

	"github.com/everettraven/wranglr/pkg/modules"
	"github.com/everettraven/wranglr/pkg/modules/github"
	"github.com/everettraven/wranglr/pkg/modules/jira"
	"github.com/everettraven/wranglr/pkg/printers"
	"go.starlark.net/starlark"
)

type Module struct {
	Output string
}

func (m *Module) String() string        { return "wranglr" }
func (m *Module) Type() string          { return "Module" }
func (m *Module) Truth() starlark.Bool  { return starlark.False }
func (m *Module) Freeze()               {}
func (m *Module) Hash() (uint32, error) { return 0, fmt.Errorf("hashing not yet implemented") }

const RenderAttr = "render"

func (m *Module) Attr(name string) (starlark.Value, error) {
	switch name {
	case RenderAttr:
		return starlark.NewBuiltin(RenderAttr, RenderBuiltin(m.Output)), nil
	default:
		return nil, fmt.Errorf("unknown attribute %q", name)
	}
}

func (m *Module) AttrNames() []string {
	return []string{
		RenderAttr,
	}
}

func RenderBuiltin(output string) modules.BuiltinFunc {
	return func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		values := []starlark.Value{}
		for i, arg := range args {
			list, ok := arg.(*starlark.List)
			if !ok {
				return starlark.None, fmt.Errorf("wranglr.render(): positional arguments must be lists, but positional argument %d was type %s", i, arg.Type())
			}

			for elem := range list.Elements() {
				switch elem.(type) {
				case *github.Item, *jira.Item:
					// do nothing, valid case
				default:
					return starlark.None, fmt.Errorf("wranglr.render(): positional arguments must be lists of supported types, but positional argument %d contains unsupported elements", i)
				}

				values = append(values, elem)
			}
		}

		// TODO: a printer registry, or should each value be responsible for implementing an output interface?
		switch output {
		case "json":
			p := printers.JSON{}
			err := p.Print(values...)
			if err != nil {
				return starlark.None, err
			}
		case "interactive":
			p := printers.Interactive{}
			err := p.Print(values...)
			if err != nil {
				return starlark.None, err
			}
		default:
			return starlark.None, fmt.Errorf("unknown output format %q", output)
		}

		return starlark.None, nil
	}
}
