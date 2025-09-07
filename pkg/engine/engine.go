package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

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

type resultPool struct {
	mu      sync.Mutex
	results map[string][]plugins.SourceEntry
}

func (e *Engine) Run(ctx context.Context, thread *starlark.Thread) ([]plugins.SourceEntry, error) {
	results := &resultPool{
		mu:      sync.Mutex{},
		results: make(map[string][]plugins.SourceEntry),
	}

	resultChan := make(chan plugins.SourceEntry)

	processWaitGroup := sync.WaitGroup{}
	processWaitGroup.Add(1)

	// Process incoming results
	go func() {
		for {
			shouldExit := false
			select {
			case <-ctx.Done():
				shouldExit = true
			case entry, ok := <-resultChan:
				if !ok {
					// channel closed which means we are done processing
					shouldExit = true
					break
				}

				results.mu.Lock()
				existing := results.results[entry.Source()]
                existing = appendUniqueResults(existing, entry)
                results.results[entry.Source()] = existing
				results.mu.Unlock()
			}

			if shouldExit {
				break
			}
		}
		processWaitGroup.Done()
	}()

	sourceWaitGroup := sync.WaitGroup{}
	sourceErrs := []error{}

	for _, plugin := range e.plugins {
		for _, source := range plugin.Sources() {
			sourceWaitGroup.Add(1)
			go func() {
				err := source.Fetch(ctx, thread, resultChan)
				if err != nil {
					sourceErrs = append(sourceErrs, fmt.Errorf("fetching data for source %q: %w", source.Name(), err))
				}

				sourceWaitGroup.Done()
			}()
		}
	}

    sourceWaitGroup.Wait()
    close(resultChan)

    processWaitGroup.Wait()

	out := []plugins.SourceEntry{}

	for _, entry := range results.results {
		out = append(out, entry...)
	}

	return out, errors.Join(sourceErrs...)
}

func appendUniqueResults(base []plugins.SourceEntry, toAppends ...plugins.SourceEntry) []plugins.SourceEntry {
	out := slices.Clone(base)

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
