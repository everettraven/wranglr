package engine

import (
	"context"
	"fmt"
	"slices"

	"github.com/everettraven/synkr/pkg/plugins"
	"go.starlark.net/starlark"
)

type Engine struct {
	plugins []plugins.Plugin
}

func New(plugins ...plugins.Plugin) *Engine {
	return &Engine{
		plugins: plugins,
	}
}

// TODO: Pulling data from sources should be asynchronous to improve performance
func (e *Engine) Run(ctx context.Context, thread *starlark.Thread) ([]plugins.SourceResult, error) {
	results := []plugins.SourceResult{}

	for _, plugin := range e.plugins {
		for _, source := range plugin.Sources() {
			result, err := source.Fetch(ctx, thread)
			if err != nil {
				return nil, fmt.Errorf("fetching data for source %q: %w", source.Name(), err)
			}

			if result == nil {
				continue
			}

			if ind := slices.IndexFunc(results, func(e plugins.SourceResult) bool {
				return e.Source == result.Source && e.Project == result.Project
			}); ind != -1 {
				results[ind].Items = appendUniqueResults(results[ind].Items, result.Items...)
				continue
			}


			results = append(results, *result)
		}
	}

	return results, nil
}

func appendUniqueResults(base []plugins.SourceEntry, toAppends ...plugins.SourceEntry) []plugins.SourceEntry {
	out := append([]plugins.SourceEntry{}, base...)

	for _, toAppend := range toAppends {
		if slices.ContainsFunc(out, func(e plugins.SourceEntry) bool {
			return e.Identifier() == toAppend.Identifier()
		}) {
			continue
		}

		out = append(out, toAppend)
	}

	return out
}
