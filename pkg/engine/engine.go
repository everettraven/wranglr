package engine

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
)

type Engine struct {
	sources []Source
	printer Printer
}

type Source interface {
	Fetch(context.Context, *starlark.Thread) ([]any, error)
	Name() string
	Project() string
}

type Printer interface {
	Print(...SourceResult) error
}

type SourceResult struct {
	Source  string `json:"source"`
	Project string `json:"project"`
	Items   []any  `json:"items"`
}

func (e *Engine) Run(ctx context.Context, thread *starlark.Thread) error {
	outputs := []SourceResult{}
	for _, source := range e.sources {
		items, err := source.Fetch(ctx, thread)
		if err != nil {
			return err
		}

		outputs = append(outputs, SourceResult{
			Source:  source.Name(),
			Project: source.Project(),
			Items:   items,
		})
	}

	err := e.printer.Print(outputs...)
	if err != nil {
		return fmt.Errorf("printing source results: %w", err)
	}

	return nil
}

func (e *Engine) AddSource(source Source) {
	e.sources = append(e.sources, source)
}

func (e *Engine) SetPrinter(printer Printer) {
	e.printer = printer
}
