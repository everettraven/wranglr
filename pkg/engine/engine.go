package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

type Engine struct {
	sources []Source
}

type Source interface {
	Fetch(context.Context, *starlark.Thread) ([]any, error)
	Name() string
	Project() string
}

type output struct {
	Source  string `json:"source"`
	Project string `json:"project"`
	Items   []any  `json:"items"`
}

func (e *Engine) Run(ctx context.Context, thread *starlark.Thread) error {
	outputs := []output{}
	for _, source := range e.sources {
		items, err := source.Fetch(ctx, thread)
		if err != nil {
			return err
		}

		outputs = append(outputs, output{
			Source:  source.Name(),
			Project: source.Project(),
			Items:   items,
		})
	}

	outs := []string{}
	for _, output := range outputs {
		outBytes, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("marshalling output %v to json: %w", output, err)
		}

		outs = append(outs, string(outBytes))
	}

	fmt.Print(strings.Join(outs, "\n"))
	return nil
}

func (e *Engine) AddSource(source Source) {
	e.sources = append(e.sources, source)
}
