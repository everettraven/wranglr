package github

import "github.com/everettraven/synkr/pkg/plugins"

type Sourcer struct {
	sources []plugins.Source
}

func (s *Sourcer) AddSource(src plugins.Source) {
	s.sources = append(s.sources, src)
}

func (s *Sourcer) Sources() []plugins.Source {
	return s.sources
}
